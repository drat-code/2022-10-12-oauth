FROM golang:1.18-alpine

RUN apk update --no-cache && apk upgrade --no-cache
RUN apk --no-cache add gcc g++ make git
RUN go install github.com/codegangsta/gin@latest

WORKDIR /app

CMD ["sh", "-c", "PORT=5000 gin --immediate -port=$GIN_PORT -appPort=5000 run main.go"]