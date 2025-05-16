FROM golang:1.24

WORKDIR /

COPY . .

RUN go build -o main .

EXPOSE 8082

CMD ["./main"]