FROM golang:1.8

WORKDIR /go/src/github.com/waypoint/waypoint
COPY . .

RUN go-wrapper install

CMD ["go-wrapper", "run"]
