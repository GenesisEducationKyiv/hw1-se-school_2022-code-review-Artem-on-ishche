# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY .env  ./

RUN mkdir data
COPY data/*.go ./data/

RUN mkdir emails
COPY emails/*.go ./emails/

RUN mkdir exchange_rate
COPY exchange_rate/*.go ./exchange_rate/

RUN mkdir smtp
COPY smtp/*.go ./smtp/

RUN go build -o /gses2.app/api

EXPOSE 8080

CMD [ "/gses2.app/api" ]