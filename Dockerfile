FROM golang:1.9-alpine

EXPOSE 8080

COPY main.go .

RUN go build main.go
RUN chmod +x main.go

CMD ./main
