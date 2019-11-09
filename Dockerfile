FROM golang:1.12.9

WORKDIR /go/src/github.com/janabe/cscoupler

# copies the contents of the current dir (.) to the container destination
ADD . /go/src/github.com/janabe/cscoupler

# get all dependencies
RUN go get ./...

EXPOSE 3000

# runs: go run main.go
CMD [ "go", "run", "main.go" ]
