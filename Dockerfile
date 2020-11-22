FROM golang:alpine

RUN adduser -h /leave_a_message -D lam

USER lam

ADD . /leave_a_message
WORKDIR /leave_a_message

RUN go build

RUN rm -f go.* *.go Dockerfile

ENTRYPOINT /leave_a_message/leave-a-message
