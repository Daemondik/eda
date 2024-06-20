FROM golang:latest

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod tidy

RUN go test ./... || exit 1

CMD ["air"]