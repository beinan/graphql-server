FROM golang:1.10

ADD . /go/src/github.com/beinan/graphql-server
WORKDIR /go/src/github.com/beinan/graphql-server

RUN go get github.com/pilu/fresh

RUN go build
RUN go install

CMD fresh
