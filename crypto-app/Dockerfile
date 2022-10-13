FROM golang:1.18-alpine AS build
RUN apk add build-base

WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM golang:1.18-alpine AS run

COPY --from=build /app/main /app/main
COPY --from=build /app/.env .env

EXPOSE 8080
CMD /app/main