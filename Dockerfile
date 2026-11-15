FROM golang:1.26-alpine

WORKDIR /app

COPY . .

RUN go build -o scheduler

CMD ["./scheduler"]