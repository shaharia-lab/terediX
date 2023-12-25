FROM alpine:3.16

LABEL org.opencontainers.image.source https://github.com/shaharia-lab/teredix

RUN apk add --no-cache ca-certificates git bash curl jq openssh-client gnupg

COPY teredix /usr/local/bin/teredix

CMD ["/usr/local/bin/teredix", "--version"]