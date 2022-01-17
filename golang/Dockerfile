FROM golang:1.14-stretch

COPY ./entrypoint.sh /tmp/entrypoint.sh
ADD . /go/src/app
WORKDIR /go/src/app

ENV MONGODB_URI=${MONGODB_URI}
ENV PORT=${PORT}
ENV TEST_PORT=${TEST_PORT}

CMD apt-get update -y &&\
    apt-get install -y inotify-tools &&\
    chmod +x /tmp/entrypoint.sh

#RUN go get go.mongodb.org/mongo-driver/mongo
#RUN update-ca-certificates

ENTRYPOINT ["/tmp/entrypoint.sh"]
