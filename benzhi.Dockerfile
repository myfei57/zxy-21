FROM golang:1.23.12-bookworm

ENV GOPROXY=off GOSUMDB=off GOFLAGS=-mod=vendor

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

RUN go build -mod=vendor ./...

EXPOSE 8080
CMD ["go", "run", "-mod=vendor", "./cmd/limsd"]
