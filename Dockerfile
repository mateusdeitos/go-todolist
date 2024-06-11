# primeiro builda o projeto
FROM golang:1.22-alpine AS build

WORKDIR /usr/app/src

COPY src/ .
RUN go mod download

RUN go build -v -o server

# depois roda
FROM alpine:3

WORKDIR /usr/app/src

RUN apk add go

COPY --from=build /usr/app/src .

CMD ["/usr/app/src/server"]

EXPOSE 9000
