FROM golang:1.22-alpine3.21 AS build

RUN mkdir /app
WORKDIR /app
COPY go.mod /app/
COPY go.sum /app/
RUN go get github.com/eclipse/paho.mqtt.golang
COPY *.go /app/
RUN go build .

FROM golang:1.22-alpine3.21
RUN apk add python3
RUN mkdir /app
COPY --from=build /app/mqtt-rules-engine /app
WORKDIR /app
CMD [ "/app/mqtt-rules-engine" ]