FROM --platform=linux/amd64 golang:1.23.6-bookworm AS build
LABEL authors="aleksandr"

WORKDIR /api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./ ./
RUN go build -o consumer ./cmd/consumer/main.go

FROM --platform=linux/amd64 debian:bookworm
#RUN useradd appuser
WORKDIR /api
COPY --from=build /api/consumer .
#USER appuser
ENTRYPOINT ["/api/consumer"]