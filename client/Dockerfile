# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------

FROM golang:1.21.1 AS builder

WORKDIR /go/src/

COPY proto ./proto
COPY certs ./certs

COPY go.mod go.sum ./
COPY client/main.go ./main.go

RUN go get -d -v golang.org/x/net/html

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/main ./grpc_client
ENTRYPOINT ["./grpc_client"]
