FROM golang:1.23

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    PATH="/root/go/bin:${PATH}"

WORKDIR /app

RUN apt-get update && apt-get install -y curl && \
    curl -o air.tar.gz -L https://github.com/cosmtrek/air/releases/download/v1.44.0/air_1.44.0_linux_amd64.tar.gz && \
    tar -xvzf air.tar.gz -C /usr/local/bin && \
    rm air.tar.gz && \
    go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["sh", "-c", "swag init -g main.go && air"]

