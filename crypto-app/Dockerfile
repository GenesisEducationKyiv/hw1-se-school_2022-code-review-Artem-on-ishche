FROM golang:1.18-alpine AS build
RUN apk add build-base

WORKDIR /crypto-app
COPY . /crypto-app

RUN go mod download
RUN go build -o /gses2.app/api

FROM golang:1.18-alpine as run
COPY --from=build /crypto-app /crypto-app

EXPOSE 8080

CMD "/gses2.app/api"