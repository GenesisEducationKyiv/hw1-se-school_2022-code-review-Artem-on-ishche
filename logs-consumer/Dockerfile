FROM golang:1.18-alpine AS build
RUN apk add build-base

WORKDIR /logs-consumer
COPY . /logs-consumer

RUN go mod download
RUN go build -o /gses2.app/logs

FROM golang:1.18-alpine as run
COPY --from=build /logs-consumer /logs-consumer

EXPOSE 8080

CMD "/gses2.app/logs"