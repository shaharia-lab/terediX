FROM scratch
COPY example /usr/bin/teredix
ENTRYPOINT ["/usr/bin/teredix"]