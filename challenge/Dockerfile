FROM golang:1.20-alpine

WORKDIR /app
COPY . .

RUN go build -o server .

ENV FLAG=CTS{default_flag}
COPY flag.txt /flag.txt

EXPOSE 1337
CMD ["./server"]
