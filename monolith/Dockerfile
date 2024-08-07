FROM golang:1.21.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main .

EXPOSE 8080

# Add a short sleep to ensure the database is ready
CMD ["sh", "-c", "sleep 5 && ./main"]