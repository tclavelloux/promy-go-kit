## Multistage build
FROM golang:1.26 as build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /app

## Multistage deploy
FROM gcr.io/distroless/base-debian10

# Kodata simulation for Yokai framework
ENV KO_DATA_PATH=/var/run/ko

WORKDIR /

# Copy application configs
COPY --from=build /src/configs /configs

# Copy database migrations to KO_DATA_PATH (required for Yokai SQL migrations)
COPY --from=build /src/db/migrations /var/run/ko/db/migrations

# Copy compiled application
COPY --from=build /app /app

# Expose ports
# 8080 - Main HTTP API
# 8081 - Metrics and health dashboard (optional, for internal monitoring)
EXPOSE 8080 8081

ENTRYPOINT ["/app"]
CMD ["run"]
