
FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN make

EXPOSE 6379

CMD ["./bin/pastebin"]