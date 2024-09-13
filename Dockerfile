FROM alpine:latest
ARG VERSION=latest
RUN apk add \
   curl \
   ca-certificates \
   openssl \
   tar \
   gzip \
   pigz \
   bzip2 \
   coreutils \
   lzip \
   xz-dev \
   lzop \
   xz \
   zstd \
   tzdata \
   && \
   rm -rf /var/cache/apk/*

ADD install /install
RUN chmod +x /install \
   && /install ${VERSION} \
   && rm /install

CMD ["/usr/local/bin/vtsbackup", "run"]
