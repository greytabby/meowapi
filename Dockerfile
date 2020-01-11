# Builder image
FROM golang:1.13 as builder
WORKDIR /go/src/github.com/greytabby/meowapi
RUN go get -v github.com/labstack/echo
RUN go get -v github.com/go-gorp/gorp
RUN go get -v github.com/go-sql-driver/mysql
RUN go get -v github.com/dgrijalva/jwt-go
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -v -o meowapi

# Runtime image
FROM alpine:3.7
COPY --from=builder /go/src/github.com/greytabby/meowapi/meowapi /meowapi
EXPOSE 8080
ENTRYPOINT ["/meowapi"]
