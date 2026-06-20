# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /app/out ./cmd/api

RUN mkdir -p /data

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /app/out /app/bff
COPY --from=build --chown=nonroot:nonroot /data /data

ENV PORT=8080 SQLITE_PATH=/data/kinji.db

EXPOSE 8080

USER nonroot

ENTRYPOINT ["/app/bff"]