FROM golang:1.18.1

WORKDIR /app

ADD ./go.mod ./
ADD ./go.sum ./

RUN go mod download

ADD ./api/owner ./api/owner
ADD ./cmd/owner ./cmd/owner
ADD ./internal/owner ./internal/owner
ADD ./configs/owner ./configs/owner
ADD ./infrastructure/owner ./infrastructure/owner
ADD ./tool ./tool
ADD ./proto ./proto

RUN go build -o /server /app/cmd/owner

EXPOSE 8000

CMD ["/server"]