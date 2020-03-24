FROM golang:1.14.0-alpine3.11

ENV LOG_LEVEL 'info'
ENV ETCD_URL ''
ENV ETCD_SERVICE_PATH ''
ENV ETCD_TIMEOUT_SECONDS '30'

WORKDIR /etcd-hashring

WORKDIR /etcd-hashring
ADD go.mod /etcd-hashring/
RUN go mod download

WORKDIR /etcd-hashring/simple-server
ADD simple-server/go.mod /etcd-hashring/simple-server/
RUN go mod download

ADD / /etcd-hashring
RUN go build -o /usr/bin/etcd-hashring-simple-server

ADD simple-server/startup.sh /

CMD [ "/startup.sh" ]

