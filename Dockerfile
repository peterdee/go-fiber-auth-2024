FROM golang:alpine3.19 as builder
WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server main.go

FROM alpine as release
WORKDIR /release
COPY --from=builder /build/server /release/server
COPY --from=builder /build/assets/favicon.ico /release/assets/favicon.ico
ENV ENV_SOURCE "env"
EXPOSE 2024
RUN chmod +x server
ENTRYPOINT ["./server"]
