FROM golang:1.18.1-alpine3.15

RUN go version
ENV GOPATH=/

COPY . .

RUN go mod download
RUN go build -o sha256sum ./cmd/main.go

CMD ["./sha256sum"]