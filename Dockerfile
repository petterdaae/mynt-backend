FROM golang:buster

WORKDIR $GOPATH/src/backend

ENV  GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN  go mod download

COPY cmd/       ./cmd
COPY internal/  ./internal
COPY sql/       ./sql
COPY main.go    main.go

RUN go build -o backend .

CMD ["./backend"]
