FROM relnod/evo-base:latest

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN go get github.com/cortesi/modd/cmd/modd

EXPOSE 8080

CMD ["modd"]
