FROM golang:1.26

WORKDIR /app

COPY . .

RUN go build -o scheduler

CMD ["./scheduler"]

EXPOSE 7540