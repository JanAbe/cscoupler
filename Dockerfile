FROM golang:1.12.9

WORKDIR /go/src/github.com/janabe/cscoupler
ADD . /go/src/github.com/janabe/cscoupler

RUN go get ./...

EXPOSE 3000
CMD [ "go", "run", "main.go" ]