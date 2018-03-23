FROM golang:1.10 as build

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN cd cmd/evo-server && go get -v ./...
RUN cd cmd/evo-server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /main

FROM scratch

COPY --from=build /main .

EXPOSE 8080

CMD ["./main"]
