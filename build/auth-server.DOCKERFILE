FROM golang:1.18.1

WORKDIR /app

ADD ./go.mod ./
ADD ./go.sum ./

RUN go mod download

ADD ./api/auth ./api/auth
ADD ./cmd/auth ./cmd/auth
ADD ./internal/auth ./internal/auth
ADD ./configs/auth ./configs/auth
ADD ./infrastructure/auth ./infrastructure/auth

RUN go build -o /server /app/cmd/auth

EXPOSE 8000

CMD ["/server"]