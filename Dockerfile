FROM golang:1.12.9 as builder

WORKDIR /go/src/github.com/janabe/cscoupler

# copies the contents of the current dir (.) to the container destination
ADD . /go/src/github.com/janabe/cscoupler

# get all dependencies
RUN go get ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

COPY --from=builder /go/src/github.com/janabe/cscoupler/.secret.json .
COPY --from=builder /go/src/github.com/janabe/cscoupler/main .

EXPOSE 3000

# runs: go run main.go
CMD [ "./main" ]
