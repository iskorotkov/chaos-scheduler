FROM alpine:latest as base

# prepare
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app

# restore packages
COPY ["go.mod", "go.sum", "./"]
RUN go get -d -v ./...

# build
COPY . .
RUN go install -v ./...

# run
FROM base as run
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/scheduler /app
ENTRYPOINT ["./app"]
