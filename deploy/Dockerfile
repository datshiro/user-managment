# syntax=docker/dockerfile:1
# Build Stage
FROM golang:1.22 AS build-stage

WORKDIR /app
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server-app


# Run Stage
FROM scratch
WORKDIR /

COPY --from=build-stage /server-app /server-app

ENTRYPOINT ["./server-app"]

