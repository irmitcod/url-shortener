# deploying stage
FROM golang:1.16

WORKDIR /build

RUN apt update && apt -y install libvips42

ADD . .

EXPOSE 8080

CMD ["./url-shortener"]