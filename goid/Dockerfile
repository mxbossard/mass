FROM docker:20.10-dind AS dind

FROM golang:1.18 as go

#RUN apk update && apk add bash docker-compose mosquitto-clients curl
#RUN apk add git
#RUN apk add gettext
#RUN apk add jq

ENV TERM=xterm-256color

COPY entrypoint.sh /root/entrypoint.sh
ENTRYPOINT ["/root/entrypoint.sh"]

COPY --from=dind /usr/local/bin /usr/local/docker/bin
COPY --from=dind /certs /certs

ENV PATH="${PATH}:/usr/local/docker/bin"

#ENV DOCKER_TLS_CERTDIR=/certs

