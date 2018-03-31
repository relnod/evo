FROM golang:1.10 as build

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN make build-web

FROM nginx:alpine

COPY --from=build /go/src/github.com/relnod/evo/cmd/evo-web/scripts/startup.sh /
RUN chmod +x /startup.sh
COPY --from=build /go/src/github.com/relnod/evo/cmd/evo-web/static /usr/share/nginx/html

ENTRYPOINT ["/startup.sh", "/usr/share/nginx/html/client.html"]
CMD ["localhost:8080"]
