FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /bin/worker ./cmd/worker
RUN CGO_ENABLED=0 go build -o /bin/migrate ./cmd/migrate

FROM alpine:3.20 AS api
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /bin/api /usr/local/bin/api
ENTRYPOINT ["api"]

FROM alpine:3.20 AS worker
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /bin/worker /usr/local/bin/worker
ENTRYPOINT ["worker"]

FROM alpine:3.20 AS migrate
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /bin/migrate /usr/local/bin/migrate
COPY db/migrations /migrations
ENTRYPOINT ["migrate"]
