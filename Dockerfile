# SOURCE CODE AND DEPENDENCIES
FROM golang:1.23.0-alpine AS builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o bin/app cmd/app/main.go

# FINAL STAGE
FROM alpine AS final
RUN mkdir logs
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/bin/app /app
EXPOSE 8080
CMD ["/app"]
