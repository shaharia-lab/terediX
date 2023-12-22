FROM --platform=$BUILDPLATFORM golang:1.20-alpine as builder

RUN apk add --no-cache make git
WORKDIR /workspace/teredix

COPY . /workspace/teredix/
RUN go mod download

COPY . /workspace/teredix
RUN make dist

# -----------------------------------------------------------------------------

FROM alpine:3.16

LABEL org.opencontainers.image.source https://github.com/shaharia-lab/teredix

RUN apk add --no-cache ca-certificates git bash curl jq openssh-client gnupg

COPY --from=builder /workspace/teredix/dist/teredix /usr/local/bin/teredix

CMD ["/usr/local/bin/teredix"]