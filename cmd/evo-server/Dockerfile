FROM relnod/evo-base:latest as build

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN make build-server

FROM scratch

COPY --from=build /go/src/github.com/relnod/evo/out/evo-server .

EXPOSE 8080

CMD ["./evo-server"]
