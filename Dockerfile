FROM golang:1.25.6-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:3.20

WORKDIR /app
COPY --from=build /app/api /app/api
COPY config.yaml /app/config.yaml

EXPOSE 8080
CMD ["/app/api"]
