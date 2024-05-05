FROM golang:1.22 as builder
WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server main.go
ENV ENV_SOURCE env
CMD ["/build/server"]
