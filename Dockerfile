FROM golang:1.13
RUN mkdir /code

ENV GOBIN /go/bin
ENV PATH $PATH:$GOBIN

WORKDIR /code

CMD make serve