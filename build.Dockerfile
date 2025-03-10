FROM golang:1.21
ARG TARGETARCH
ARG VERSION
COPY . /go/src/github.com/elwinar/rambler
WORKDIR /go/src/github.com/elwinar/rambler
RUN go get ./...
RUN go build -ldflags="-X main.VERSION=${VERSION} -s -linkmode external -extldflags -static -w"

FROM scratch
COPY --from=0 /go/src/github.com/elwinar/rambler/rambler /
