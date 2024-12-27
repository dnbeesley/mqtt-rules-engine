FROM golang:1.22-alpine3.21 AS build

RUN mkdir /app
COPY go.mod /app/
COPY go.sum /app/
COPY *.go /app/
WORKDIR /app
RUN go build .

FROM golang:1.22-alpine3.21
RUN mkdir /app
COPY --from=build /app/mqtt-rules-engine /app
ENV DATA-DIR=/data
CMD [ "/app/mqtt-rules-engine" ]