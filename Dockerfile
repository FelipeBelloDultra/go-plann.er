FROM golang:1.22.4-alpine

WORKDIR /planner

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./bin/planner ./cmd/planner/planner.go

EXPOSE 8080

ENTRYPOINT [ "/planner/bin/planner" ]
