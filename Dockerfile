# primeiro builda o projeto
FROM golang:1.22-alpine AS build

WORKDIR /usr/app

COPY . .
RUN go mod download

COPY . .

RUN go build -v -o server

# depois roda
FROM alpine:3

WORKDIR /usr/app

COPY --from=build /usr/app .

CMD ["/usr/app/server"]

EXPOSE 9000
