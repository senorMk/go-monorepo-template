# Build the Go API. Build context is the repo root (see .github/workflows).
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download

COPY apps/api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app
ENV APP_ENV=production
ENV PORT=3000

COPY --from=build /out/api /app/api

EXPOSE 3000
ENTRYPOINT ["/app/api"]
