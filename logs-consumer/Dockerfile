FROM golang:1.18-alpine AS build
RUN apk add build-base

WORKDIR /app
COPY . .
RUN go build -o logs main.go

FROM golang:1.18-alpine AS run
COPY --from=build /app/logs /app/logs

CMD /app/logs