FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o go-chat-app

RUN chmod +x go-chat-app

EXPOSE 4000

EXPOSE 8080


CMD ["./go-chat-app"]