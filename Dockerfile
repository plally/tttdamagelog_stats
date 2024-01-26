FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY internal ./internal
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

COPY db/migrations/* /app/db/migrations/


EXPOSE 8080

CMD ["/server"]
