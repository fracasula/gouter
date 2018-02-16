FROM golang:1.9-alpine

EXPOSE 8080

COPY http http
COPY main.go .

RUN go build main.go
RUN chmod +x main.go
RUN rm -f main.go
RUN rm -Rf http

CMD ./main
