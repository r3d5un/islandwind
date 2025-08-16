# BASE IMAGE
FROM golang:1.25-alpine AS base

WORKDIR /app

COPY . /app
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/api

# RUNNER
FROM alpine:latest
WORKDIR /app
COPY --from=base /app/api .

EXPOSE 8080
CMD [ "./api" ]

