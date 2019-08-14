FROM golang:1.10-alpine AS builder

RUN apk add git
RUN go get github.com/influxdata/influxdb1-client/v2
RUN go get github.com/shirou/gopsutil/cpu
RUN go get github.com/shirou/gopsutil/mem
RUN apk add curl

ADD ./src /go/src/app
WORKDIR /go/src/app

ENV PORT=8915
ENV DEV1READ1=
ENV DEV1READ2=
ENV DEVICE1=
ENV ADDR=

HEALTHCHECK --interval=30s --timeout=30s CMD curl -sS http://localhost:8915/ || exit 1

CMD ["go", "run", "messages.go"]
