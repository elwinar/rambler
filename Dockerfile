FROM scratch
MAINTAINER Romain Baugue <romain.baugue@elwinar.com>

COPY rambler /
CMD ["/rambler", "apply", "-a"]
