package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-plann.er/internal/api/spec"
	"go-plann.er/internal/pgstore"
	"go.uber.org/zap"
)

type store interface {
	CreateTrip(ctx context.Context, pool *pgxpool.Pool, params spec.CreateTripRequest) (uuid.UUID, error)
	GetParticipant(ctx context.Context, participantID uuid.UUID) (pgstore.Participant, error)
	GetTrip(ctx context.Context, id uuid.UUID) (pgstore.Trip, error)
	ConfirmParticipant(ctx context.Context, participantID uuid.UUID) error
	UpdateTrip(ctx context.Context, arg pgstore.UpdateTripParams) error
	GetTripActivities(ctx context.Context, tripID uuid.UUID) ([]pgstore.Activity, error)
	CreateActivity(ctx context.Context, arg pgstore.CreateActivityParams) (uuid.UUID, error)
}

type Mailer interface {
	SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error
}

type API struct {
	store     store
	logger    *zap.Logger
	validator *validator.Validate
	pool      *pgxpool.Pool
	mailer    Mailer
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger, mailer Mailer) API {
	validator := validator.New(validator.WithRequiredStructEnabled())
	return API{pgstore.New(pool), logger, validator, pool, mailer}
}

// Confirms a participant on a trip.
// (PATCH /participants/{participantId}/confirm)
func (api API) PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request, participantID string) *spec.Response {
	id, err := uuid.Parse(participantID)
	if err != nil {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "invalid UUID"})
	}

	participant, err := api.store.GetParticipant(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "participant not found"})
		}

		api.logger.Error("failed to get participant", zap.Error(err), zap.String("participant_id", participantID))
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	if participant.IsConfirmed {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
			Message: "participant already confirmed",
		})
	}

	if err := api.store.ConfirmParticipant(r.Context(), id); err != nil {
		api.logger.Error("failed to confirm participant", zap.Error(err), zap.String("participant_id", participantID))
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	return spec.PatchParticipantsParticipantIDConfirmJSON204Response(nil)
}

// Create a new trip
// (POST /trips)
func (api API) PostTrips(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.CreateTripRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "invalid json"})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "invalid input: " + err.Error()})
	}

	tripID, err := api.store.CreateTrip(r.Context(), api.pool, body)
	if err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "failed to create trip, try again"})
	}

	go func() {
		if err := api.mailer.SendConfirmTripEmailToTripOwner(tripID); err != nil {
			api.logger.Error(
				"failed to send confirmation email on PostTrips",
				zap.Error(err),
				zap.String("trip_id", tripID.String()),
			)
		}
	}()

	return spec.PostTripsJSON201Response(spec.CreateTripResponse{TripID: tripID.String()})
}

// Get a trip details.
// (GET /trips/{tripId})
func (api API) GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDJSON400Response(spec.Error{Message: "invalid UUID"})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.GetTripsTripIDJSON400Response(spec.Error{Message: "trip not found"})
		}

		api.logger.Error("failed to get trips", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	return spec.GetTripsTripIDJSON200Response(spec.GetTripDetailsResponse{
		Trip: spec.GetTripDetailsResponseTripObj{
			ID:          trip.ID.String(),
			Destination: trip.Destination,
			StartsAt:    trip.StartsAt.Time,
			EndsAt:      trip.EndsAt.Time,
			IsConfirmed: trip.IsConfirmed,
		},
	})
}

// Update a trip.
// (PUT /trips/{tripId})
func (api API) PutTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var body spec.UpdateTripRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "invalid json"})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "invalid inputs: " + err.Error()})
	}

	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "invalid UUID"})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "trip not found"})
		}

		api.logger.Error("failed to get trips", zap.Error(err), zap.String("trip_id", tripID))
		return spec.PutTripsTripIDJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	if err := api.store.UpdateTrip(r.Context(), pgstore.UpdateTripParams{
		Destination: body.Destination,
		StartsAt:    pgtype.Timestamp{Valid: true, Time: body.StartsAt},
		EndsAt:      pgtype.Timestamp{Valid: true, Time: body.EndsAt},
		IsConfirmed: trip.IsConfirmed,
		ID:          trip.ID,
	}); err != nil {
		api.logger.Error("failed to update trip", zap.Error(err), zap.String("trip_id", tripID))
		return spec.PutTripsTripIDJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	return spec.PutTripsTripIDJSON204Response(nil)
}

// Get a trip activities.
// (GET /trips/{tripId}/activities)
func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "invalid UUID"})
	}

	if _, err := api.store.GetTrip(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "trip not found"})
		}

		api.logger.Error("failed to get trips", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	activities, err := api.store.GetTripActivities(r.Context(), id)
	if err != nil {
		api.logger.Error("failed to get trips activities", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	activityMap := make(map[string][]spec.GetTripActivitiesResponseInnerArray)
	for _, activity := range activities {
		date := activity.OccursAt.Time.Format(time.DateOnly)
		activityMap[date] = append(activityMap[date], spec.GetTripActivitiesResponseInnerArray{
			ID:       activity.ID.String(),
			OccursAt: activity.OccursAt.Time,
			Title:    activity.Title,
		})
	}

	var response spec.GetTripActivitiesResponse
	for date, activities := range activityMap {
		parsedDate, _ := time.Parse(time.DateOnly, date)
		response.Activities = append(response.Activities, spec.GetTripActivitiesResponseOuterArray{
			Date:       parsedDate,
			Activities: activities,
		})
	}

	return spec.GetTripsTripIDActivitiesJSON200Response(response)
}

// Create a trip activity.
// (POST /trips/{tripId}/activities)
func (api API) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var body spec.CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "invalid json"})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "invalid inputs: " + err.Error()})
	}

	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "invalid UUID"})
	}

	if _, err := api.store.GetTrip(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "trip not found"})
		}

		api.logger.Error("failed to create trip activities", zap.Error(err), zap.String("trip_id", tripID))
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	activityID, err := api.store.CreateActivity(r.Context(), pgstore.CreateActivityParams{
		TripID:   id,
		Title:    body.Title,
		OccursAt: pgtype.Timestamp{Valid: true, Time: body.OccursAt},
	})
	if err != nil {
		api.logger.Error("failed to create activity", zap.Error(err), zap.String("trip_id", tripID))
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{
			Message: "something went wrong, try again",
		})
	}

	return spec.PostTripsTripIDActivitiesJSON201Response(spec.CreateActivityResponse{
		ActivityID: activityID.String(),
	})
}

// Confirm a trip and send e-mail invitations.
// (GET /trips/{tripId}/confirm)
func (api API) GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}

// Invite someone to the trip.
// (POST /trips/{tripId}/invites)
func (api API) PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}

// Get a trip links.
// (GET /trips/{tripId}/links)
func (api API) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}

// Create a trip link.
// (POST /trips/{tripId}/links)
func (api API) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}

// Get a trip participants.
// (GET /trips/{tripId}/participants)
func (api API) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}
