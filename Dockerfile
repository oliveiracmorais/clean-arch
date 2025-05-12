#Dockerfile 
FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/ordersystem
ENTRYPOINT ["./server"] 

