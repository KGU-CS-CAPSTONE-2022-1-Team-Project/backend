FROM golang:1.18.1

WORKDIR /app

ADD ./go.mod ./
ADD ./go.sum ./

RUN go mod download

ADD ./api/gateway ./api/gateway
ADD ./cmd/gateway ./cmd/gateway
ADD ./configs/gateway ./configs/gateway
RUN go build -o /server /app/cmd/gateway

EXPOSE 10321

CMD ["/server"]