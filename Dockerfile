FROM golang:1.12.9 as builder

WORKDIR /cscoupler

# copies the contents of the current dir (.) to the container destination
ADD . /cscoupler

# get all dependencies
RUN go get ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN mkdir resumes
COPY --from=builder /cscoupler/.secret.json .
COPY --from=builder /cscoupler/main .

EXPOSE 3000

CMD [ "./main" ]
