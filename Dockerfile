# syntax=docker/dockerfile:1.4
FROM golang:1.25-alpine AS builder
LABEL maintainer="jaam8"

WORKDIR /build

COPY go.* ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
COPY migrations ./migrations
COPY .env ./.env

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o app cmd/app/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/.env .
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/app .

CMD ["./app"]
