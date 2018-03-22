FROM golang:1.10 as build

COPY ./ /go/src/github.com/relnod/evo
WORKDIR /go/src/github.com/relnod/evo

RUN go get -v github.com/gopherjs/gopherjs

RUN cd cmd/evo && go get -v -d -tags=js ./...
RUN cd cmd/evo && gopherjs build -v -m -o evo.js .

RUN cd cmd/evo-client && go get -v -d -tags=js ./...
RUN cd cmd/evo-client && gopherjs build -v -m -o evo-client.js .

FROM nginx

COPY --from=build /go/src/github.com/relnod/evo/cmd/evo/evo.js /usr/share/nginx/html
COPY --from=build /go/src/github.com/relnod/evo/cmd/evo-client/evo-client.js /usr/share/nginx/html

COPY cmd/evo-web/static/index.html /usr/share/nginx/html
COPY cmd/evo-web/static/app.html /usr/share/nginx/html
COPY cmd/evo-web/static/client.html /usr/share/nginx/html
