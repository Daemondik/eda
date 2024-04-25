FROM golang:latest

RUN go install github.com/cosmtrek/air@latest

WORKDIR /cmd

COPY . .

RUN go mod download

RUN go mod tidy

CMD ["air"]