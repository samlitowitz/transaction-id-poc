############################
# STEP 1 build executable binary
############################
FROM golang:1.12-alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
# Create appuser.
RUN adduser -D -g '' appuser
WORKDIR $GOPATH/src/github.com/samlitowitz/transaction-id-poc
COPY . .
# Fetch dependencies.
# Using go mod.
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

############################
# STEP 1.1 build app
############################
FROM builder AS app-builder
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app cmd/app/*.go

# App
FROM scratch AS app
# Import the user and group files from the builder.
COPY --from=app-builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=app-builder /go/bin/app /go/bin/app
# Use an unprivileged user.
USER appuser
# Run the binary.
ENTRYPOINT ["/go/bin/app"]

############################
# STEP 1.2 build echo
############################
FROM builder AS echo-builder
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/echo cmd/echo/*.go

# App
FROM scratch AS echo
# Import the user and group files from the builder.
COPY --from=echo-builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=echo-builder /go/bin/echo /go/bin/echo
# Use an unprivileged user.
USER appuser
# Run the binary.
ENTRYPOINT ["/go/bin/echo"]

############################
# STEP 1.3 build nats-sub
############################
# NATS Sub
FROM builder AS nats-sub-builder
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/nats-sub cmd/nats-sub/*.go

FROM scratch AS nats-sub
# Import the user and group files from the builder.
COPY --from=nats-sub-builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=nats-sub-builder /go/bin/nats-sub /go/bin/nats-sub
# Use an unprivileged user.
USER appuser
# Run the binary.
ENTRYPOINT ["/go/bin/nats-sub"]

############################
# STEP 1.4 build random-requester
############################
# Random Requester
FROM builder AS random-requester-builder
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/random-requester cmd/random-requester/*.go

FROM scratch AS random-requester
# Import the user and group files from the builder.
COPY --from=random-requester-builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=random-requester-builder /go/bin/random-requester /go/bin/random-requester
# Use an unprivileged user.
USER appuser
# Run the binary.
ENTRYPOINT ["/go/bin/random-requester"]