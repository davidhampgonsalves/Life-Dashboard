FROM golang:1.23-alpine
# Set the Current Working Directory inside the container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main /app/main.go
EXPOSE 8080
ENV GIN_MODE=release
CMD ["/app/main"]