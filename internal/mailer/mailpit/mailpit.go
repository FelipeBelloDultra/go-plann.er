package mailpit

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
	"go-plann.er/internal/pgstore"
)

type store interface {
	GetTrip(context.Context, uuid.UUID) (pgstore.Trip, error)
}

type Mailpit struct {
	store store
}

type ParticipantToSendEmail struct {
	Name  string
	Email string
}

func NewMailpit(pool *pgxpool.Pool) Mailpit {
	return Mailpit{
		store: pgstore.New(pool),
	}
}

func (mp Mailpit) SendConfirmTripEmailToTripParticipant(participant ParticipantToSendEmail, tripID uuid.UUID) error {
	ctx := context.Background()
	trip, err := mp.store.GetTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("mailpit: failed to get trip for SendConfirmTripEmailToTripParticipant: %w", err)
	}

	msg := mail.NewMsg()
	if err := msg.From("mailpit@plann.er"); err != nil {
		return fmt.Errorf("mailpit: failed to set From in email SendConfirmTripEmailToTripParticipant: %w", err)
	}

	if err := msg.To(participant.Email); err != nil {
		return fmt.Errorf("mailpit: failed to set To in email SendConfirmTripEmailToTripParticipant: %w", err)
	}

	msg.Subject("Confirm your trip")
	msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
        Hello, %s!
        Your trip to %s starting on %s needs to be confirmed.
        Click the button below to confirm.
        `,
		participant.Name, trip.Destination, trip.StartsAt.Time.Format(time.DateOnly),
	))

	client, err := mail.NewClient("mailpit", mail.WithTLSPortPolicy(mail.NoTLS), mail.WithPort(1025))
	if err != nil {
		return fmt.Errorf("mailpit: failed create email client SendConfirmTripEmailToTripParticipant: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("mailpit: failed send email client SendConfirmTripEmailToTripParticipant: %w", err)
	}

	return nil
}

func (mp Mailpit) SendConfirmTripEmailToTripParticipants(participants []ParticipantToSendEmail, tripID uuid.UUID) error {
	ctx := context.Background()
	trip, err := mp.store.GetTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("mailpit: failed to get trip for SendConfirmTripEmailToTripParticipants: %w", err)
	}

	msg := mail.NewMsg()
	if err := msg.From("mailpit@plann.er"); err != nil {
		return fmt.Errorf("mailpit: failed to set From in email SendConfirmTripEmailToTripParticipants: %w", err)
	}

	for _, participant := range participants {
		if err := msg.To(participant.Email); err != nil {
			return fmt.Errorf("mailpit: failed to set To in email SendConfirmTripEmailToTripParticipants: %w", err)
		}

		msg.Subject("Confirm your trip")
		msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
            Hello, %s!
            Your trip to %s starting on %s needs to be confirmed.
            Click the button below to confirm.
            `,
			participant.Name, trip.Destination, trip.StartsAt.Time.Format(time.DateOnly),
		))

		client, err := mail.NewClient("mailpit", mail.WithTLSPortPolicy(mail.NoTLS), mail.WithPort(1025))
		if err != nil {
			return fmt.Errorf("mailpit: failed create email client SendConfirmTripEmailToTripParticipants: %w", err)
		}

		if err := client.DialAndSend(msg); err != nil {
			return fmt.Errorf("mailpit: failed send email client SendConfirmTripEmailToTripParticipants: %w", err)
		}
	}

	return nil
}

func (mp Mailpit) SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error {
	ctx := context.Background()
	trip, err := mp.store.GetTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("mailpit: failed to get trip for SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg := mail.NewMsg()
	if err := msg.From("mailpit@plann.er"); err != nil {
		return fmt.Errorf("mailpit: failed to set From in email SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := msg.To(trip.OwnerEmail); err != nil {
		return fmt.Errorf("mailpit: failed to set To in email SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg.Subject("Confirm your trip")
	msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
		Hello, %s!
		Your trip to %s starting on %s needs to be confirmed.
        Click the button below to confirm.
		`,
		trip.OwnerName, trip.Destination, trip.StartsAt.Time.Format(time.DateOnly),
	))

	client, err := mail.NewClient("mailpit", mail.WithTLSPortPolicy(mail.NoTLS), mail.WithPort(1025))
	if err != nil {
		return fmt.Errorf("mailpit: failed create email client SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("mailpit: failed send email client SendConfirmTripEmailToTripOwner: %w", err)
	}

	return nil
}
