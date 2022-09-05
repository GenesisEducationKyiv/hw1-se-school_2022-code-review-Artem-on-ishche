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

RUN mkdir rates
COPY rates/*.go ./rates/

RUN mkdir smtp
COPY implementations/smtp/*.go ./smtp/

RUN go build -o /gses2.app/api

EXPOSE 8080

CMD [ "/gses2.app/api" ]