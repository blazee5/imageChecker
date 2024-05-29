FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app /app

EXPOSE 3000

WORKDIR /app/cmd

CMD ["./main"]