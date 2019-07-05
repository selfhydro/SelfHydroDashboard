FROM golang:alpine

WORKDIR /go/src/StateInfluxDBSubscriber
COPY . .

RUN apk update && apk add git
RUN go get github.com/influxdata/influxdb1-client/v2

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["StateInfluxDBSubscriber"]
