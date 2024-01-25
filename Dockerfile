FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/* cmd/
COPY internal/* internal/
COPY db/* db/

RUN CGO_ENABLED=0 GOOS=linux go build -o /server


EXPOSE 8080

CMD ["/server"]
