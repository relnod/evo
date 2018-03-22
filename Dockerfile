FROM golang:1.10

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN cd cmd/evo-server && go get -v ./...
RUN go build -v cmd/evo-server/main.go

EXPOSE 8080

CMD ["./main"]
