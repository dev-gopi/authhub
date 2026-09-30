FROM golang:1.25-bookworm AS builder

WORKDIR /src

ARG SERVICE=auth-api
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod verify

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/authhub \
    ./cmd/${SERVICE}

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/healthcheck \
    ./cmd/healthcheck


FROM debian:bookworm-slim

ENV TZ=UTC

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       ca-certificates \
       tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --system authhub \
    && useradd \
       --system \
       --gid authhub \
       --create-home \
       --home-dir /app \
       authhub

WORKDIR /app

COPY --from=builder /out/authhub /app/authhub
COPY --from=builder /out/healthcheck /app/healthcheck

RUN chown -R authhub:authhub /app \
    && chmod 0555 /app/authhub \
    && chmod 0555 /app/healthcheck

USER authhub:authhub

EXPOSE 8080

ENTRYPOINT ["/app/authhub"]