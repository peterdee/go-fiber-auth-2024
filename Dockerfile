FROM golang:1.22 as builder
WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /server

FROM alpine
WORKDIR /release
COPY --from=builder /server .
COPY --from=builder /build/assets/favicon.ico ./assets/favicon.ico
ENV ENV_SOURCE env
EXPOSE 2024
ENTRYPOINT ["/release/server"]
