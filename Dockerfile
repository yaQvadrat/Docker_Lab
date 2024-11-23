# SOURCE CODE AND DEPENDENCIES
FROM golang:1.23.0-alpine AS builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o bin/app cmd/app/main.go

# FINAL STAGE
FROM alpine AS final
EXPOSE 8080
RUN mkdir /app && mkdir /app/logs
WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app
USER appuser

COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/bin/app /app/app
CMD ["./app"]
