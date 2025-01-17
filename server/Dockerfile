FROM --platform=linux/amd64 golang:1.23.4-bookworm AS build
LABEL authors="aleksandr"

WORKDIR /api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . ./
RUN go build -o api ./cmd/app/main.go

FROM --platform=linux/amd64 debian:bookworm
#RUN useradd appuser
WORKDIR /api
COPY --from=build /api/api .
COPY --from=build /api/config.yaml.deply ./config.yaml
#USER appuser
ENTRYPOINT ["/api/api"]