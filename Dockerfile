FROM golang:1.18-alpine as builder
RUN apk update
RUN apk add chromium
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main

FROM alpine:3
RUN apk update
RUN apk add chromium
COPY --from=builder /app/main /bin/main
EXPOSE 9598
ENTRYPOINT ["/bin/main"]