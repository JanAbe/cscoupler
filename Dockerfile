FROM golang:1.12.9

# ARG SSH_PRIVATE_KEY
# RUN mkdir /root/.ssh/
# RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
# RUN chmod 600 /root/.ssh/id_rsa
# RUN touch /root/.ssh/config && echo "StrictHostKeyChecking no " > /root/.ssh/config

WORKDIR /go/src/cscoupler
COPY . .

# RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
# RUN go get -d -v ./...
# RUN go install ./...

EXPOSE 3000

# CMD [ "./main" ]

# RUN rm /root/.ssh/id_rsa