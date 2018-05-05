FROM alpine:latest

# Based on markpfletcher/graphviz-docker
# https://github.com/markpfletcher/graphviz-docker/blob/master/Dockerfile
RUN mkdir /graphviz && \
  apk add --update ca-certificates graphviz && \
  rm -rf /var/cache/apk/*

COPY ./dreitafel /usr/bin/
COPY ./dreitafel-web /usr/bin/
