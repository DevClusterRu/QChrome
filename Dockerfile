FROM golang:1.18.2-alpine
WORKDIR /app
COPY . .
EXPOSE 9598
RUN go mod tidy

