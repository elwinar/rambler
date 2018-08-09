FROM golang
COPY . /go/src/github.com/elwinar/rambler
WORKDIR /go/src/github.com/elwinar/rambler
RUN go get ./...
RUN go build -ldflags="-s -linkmode external -extldflags -static -w"

FROM scratch
MAINTAINER Romain Baugue <romain.baugue@elwinar.com>
COPY --from=0 /go/src/github.com/elwinar/rambler/rambler /
CMD ["/rambler", "apply", "-a"]
