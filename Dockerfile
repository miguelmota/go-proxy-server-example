FROM golang

EXPOSE 8000

ENV APP_PATH=/go/src/app/

ADD . /go/src/app
RUN go install app

WORKDIR /go/src/app
ENTRYPOINT /go/bin/app
