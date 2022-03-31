FROM alpine:3.15

LABEL maintainer="Dmitry Mozzherin"

ENV LAST_FULL_REBUILD 2022-03-31

WORKDIR /bin

COPY ./bhlindex/bhlindex /bin

RUN mkdir -p /opt/bhl

ENTRYPOINT [ "bhlindex" ]

CMD ["rest"]
