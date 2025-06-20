FROM golang:1.23-alpine

# Create non-root user
RUN adduser -D -s /bin/sh appuser

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN apk add --no-cache curl && \
    curl -o air.tar.gz -L https://github.com/cosmtrek/air/releases/download/v1.44.0/air_1.44.0_linux_amd64.tar.gz && \
    tar -xvzf air.tar.gz -C /usr/local/bin && \
    rm air.tar.gz

# Install as root, run as appuser
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Change ownership and switch to non-root user
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

CMD ["sh", "-c", "swag init -g main.go && air"]

