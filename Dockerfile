FROM alpine:3.13.5


COPY  ./config/asterix.local.crt.cer /usr/local/share/ca-certificates/my-cert.crt

RUN apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates
#RUN cat /usr/local/share/ca-certificates/my-cert.crt >> /etc/ssl/certs/ca-certificates.crt && \
 #   apk --no-cache add \
  #      curl
COPY ./config /etc/krakend/
COPY ./binaries/integration-hub /usr/bin/krakend


VOLUME [ "/etc/krakend" ]

ARG appVersion
ARG commitId
ENV APP_VERSION=$appVersion
ENV COMMIT_ID=$commitId

RUN echo IMAGE_VERSION: $APP_VERSION-$COMMIT_ID

ARG appVersion
ARG commitId
ENV APP_VERSION=$appVersion
ENV COMMIT_ID=$commitId

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json"]

EXPOSE 8090
