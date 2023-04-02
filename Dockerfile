FROM scratch
COPY teredix /usr/bin/teredix
ENTRYPOINT ["/usr/bin/teredix"]