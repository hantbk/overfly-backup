FROM alpine:latest
ARG VERSION=latest
RUN apk add \
   curl \
   ca-certificates \
   openssl \
   # replace busybox utils
   tar \
   gzip \
   pigz \
   bzip2 \
   # there is no pbzip2 yet
   lzip \
   xz-dev \
   lzop \
   xz \
   # pixz is in edge atm
   zstd \
   && \
   rm -rf /var/cache/apk/*

ADD install /install
RUN /install ${VERSION} && rm /install

CMD ["/usr/local/bin/vtsbackup", "run"]