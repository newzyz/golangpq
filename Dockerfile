FROM golang:1.18.4 AS Production
RUN mkdir /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY server.go ./

RUN go build -o /docker-gs-ping

CMD [ "/docker-gs-ping" ]

