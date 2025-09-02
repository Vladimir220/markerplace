FROM golang:1.24.5

WORKDIR /app

RUN rm /etc/localtime

RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime

CMD go mod tidy; go run .
