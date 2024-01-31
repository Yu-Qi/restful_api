FROM golang:1.20.12-alpine3.19 AS builder
WORKDIR /app

# install gcc for cgo and curl for health check
RUN apk --update --no-cache add g++

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
COPY go.mod go.sum ./
RUN go mod download

# Copy the source and build
COPY . .
ARG CMD_DIR=.
RUN go build -a -o main ${CMD_DIR}/main.go

######## Start a new stage from scratch #######
FROM alpine:3.16.7
RUN apk --no-cache add ca-certificates curl tzdata
WORKDIR /root/
COPY --from=builder /app/main .
ARG APP_PORT=8130
ENV APP_PORT ${APP_PORT}
EXPOSE 8130
CMD ["./main"]
