FROM golang:1.24-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@v1.62.0

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
