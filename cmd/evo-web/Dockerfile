FROM golang:1.10 as build

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN go get -v github.com/gopherjs/gopherjs

RUN cd cmd/evo && go get -v -d -tags=js ./...
RUN cd cmd/evo && gopherjs build -v -m -o evo.js .

RUN cd cmd/evo-client && go get -v -d -tags=js ./...
RUN cd cmd/evo-client && gopherjs build -v -m -o evo-client.js .

FROM nginx:alpine

COPY --from=build /go/src/github.com/relnod/evo/cmd/evo-web/scripts/configure.sh /
RUN chmod +x /configure.sh
COPY --from=build /go/src/github.com/relnod/evo/cmd/evo-web/static /usr/share/nginx/html

ENTRYPOINT ["/configure.sh", "bla", "/usr/share/nginx/html/client.html", ";", "nginx", "-g", "daemon off;"]
CMD ["nginx", "-g", "daemon off;"]
