ARG GO_VERSION
ARG DOCKER_VERSION=23.0.6-cli-alpine3.18
#FROM docker:20.10-dind AS dind
FROM docker:$DOCKER_VERSION AS dind

FROM golang:$GO_VERSION as go

#RUN apk update && apk add bash docker-compose mosquitto-clients curl
#RUN apk add git
#RUN apk add gettext
#RUN apk add jq

#ADD dind_bashrc /root/.bashrc

ENV TERM=xterm-256color

COPY --from=dind /usr/local/bin/docker /usr/local/bin/docker
#COPY --from=go /usr/local/go /usr/local/go
#RUN ln -s /usr/local/go/bin/go /usr/bin/go

#RUN apk update ; apk add build-base

#ENTRYPOINT ["/usr/local/go/bin/go"]

