FROM golang:1.10 as build

ENV REALIZE_VERSION=1.3

RUN go get github.com/tockins/realize
WORKDIR /go/src/github.com/tockins/realize
# RUN git checkout v${REALIZE_VERSION}
RUN go build -o /realize

FROM alpine:latest

COPY --from=build /realize /

ENTRYPOINT ["./realize"]
CMD ["start"]
