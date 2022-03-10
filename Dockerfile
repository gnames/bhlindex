FROM alpine:3.15

LABEL maintainer="Dmitry Mozzherin"

ENV LAST_FULL_REBUILD 2021-04-07

WORKDIR /bin

COPY ./bhlindex/bhlindex /bin

ENTRYPOINT [ "bhlindex" ]

CMD ["server"]
