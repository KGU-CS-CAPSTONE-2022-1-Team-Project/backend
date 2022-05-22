FROM golang:1.18.1

WORKDIR /app

ADD ./go.mod ./
ADD ./go.sum ./

RUN go mod download

ADD ./api/partner ./api/partner
ADD ./cmd/partner ./cmd/partner
ADD ./configs/partner ./configs/partner
ADD ./infrastructure/partner ./infrastructure/partner
ADD ./tool ./tool
ADD ./proto ./proto

RUN go build -o /server /app/cmd/partner

EXPOSE 15663

CMD ["/server"]