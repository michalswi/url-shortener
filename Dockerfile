ARG GOLANG_VERSION
ARG ALPINE_VERSION

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk --no-cache add make git; \
    adduser -D -h /tmp/dummy dummy

USER dummy
WORKDIR /tmp/dummy

COPY --chown=dummy Makefile Makefile
COPY --chown=dummy go.mod go.mod
COPY --chown=dummy go.sum go.sum

RUN go mod download

COPY --chown=dummy home home
COPY --chown=dummy links links
COPY --chown=dummy server server
COPY --chown=dummy main.go main.go

RUN make build

FROM alpine:${ALPINE_VERSION}

ARG VERSION
ARG NAME

# Shortener Configuration
ENV PORT="8080"

# Copy from builder
COPY --from=builder /tmp/dummy/${APPNAME}-${VERSION} /usr/bin/${APPNAME}

# Exec
CMD ["url-shortener"]