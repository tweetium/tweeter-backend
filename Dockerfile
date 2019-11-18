# --- Build the binary in a minimal Go container
FROM golang:1.13.4-alpine3.10 as build-phase

WORKDIR /app

# Download from go.mod / go.sum before getting all files
# for better layer caching.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Now we can get everything.
COPY . .

RUN go build main.go

# --- Create container with just the binary
FROM alpine:3.10

COPY --from=build-phase /app/main /app

ENTRYPOINT [ "/app" ]
