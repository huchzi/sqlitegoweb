FROM golang:1.21

WORKDIR /app

COPY sqlitedb.go .
COPY go.mod .
COPY go.sum .
COPY *.html .
COPY assets/* ./assets/

RUN go build .

EXPOSE 8080

CMD [ "./sqlitedb" ]