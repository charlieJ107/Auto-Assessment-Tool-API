FROM golang:1.19.3
LABEL authors="Charlie"
WORKDIR /app

COPY . .

RUN go build -o main .

ENV GIN_MODE=release

EXPOSE 80

CMD ["./main"]