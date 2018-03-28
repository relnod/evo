FROM golang:1.10 as build

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN go get -v github.com/gopherjs/gopherjs

RUN make build-web

FROM nginx

COPY --from=build /go/src/github.com/relnod/evo/cmd/evo-web/static /usr/share/nginx/html
