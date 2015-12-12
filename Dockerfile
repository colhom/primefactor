FROM golang:1.5
MAINTAINER colin.hom@coreos.com

ADD ./primefactor.go ./
RUN go build primefactor.go
ENTRYPOINT ["./primefactor"]