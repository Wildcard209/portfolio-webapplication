FROM golang:1.20

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


WORKDIR /app

RUN apt-get update && apt-get install -y curl && \
    curl -o air.tar.gz -L https://github.com/cosmtrek/air/releases/download/v1.44.0/air_1.44.0_linux_amd64.tar.gz && \
    tar -xvzf air.tar.gz -C /usr/local/bin && \
    rm air.tar.gz

COPY go.mod ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]

