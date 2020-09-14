FROM golang:1.13-alpine

WORKDIR /build

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache build-base \
    && rm -rf /var/cache/apk/* 

ADD ./go.mod ./go.sum ./

RUN go mod download

ADD ./server.go ./

RUN go build -o server .

WORKDIR /app

RUN cp /build/server .

EXPOSE 3000

CMD ["/app/server"]