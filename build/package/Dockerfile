FROM golang:1.25.3-alpine AS builder

WORKDIR /usr/local/src

RUN apk add --no-cache bash git make gcc gettext musl-dev ca-certificates

COPY ./go.mod ./go.sum .

RUN go mod download

COPY ./ ./

ARG TAG=unknown

RUN go build -o ./build/package/subscriptions -ldflags="-X main.version=$TAG"  ./cmd/sub/main.go

FROM alpine:3.22 AS runner

RUN apk add --no-cache ca-certificates

COPY ./configs/dev.docker.yaml /etc/subscriptions/config.yaml

ENV POSTGRES_URL=
ENV CONFIG_PATH=/etc/subscriptions/config.yaml

ENV GIN_MODE=info

COPY --from=builder /usr/local/src/build/package/subscriptions /usr/bin/subscriptions

CMD ["/usr/bin/subscriptions"]
