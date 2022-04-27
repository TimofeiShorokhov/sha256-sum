FROM golang:1.18-alpine AS buildenv
WORKDIR /sha256-sum
ADD . /sha256-sum
RUN go mod download
RUN go build -o sha256-sum cmd/main.go

RUN chmod +x sha256-sum

FROM alpine:latest
WORKDIR /app
COPY --from=buildenv /sha256-sum .
COPY --from=buildenv /sha256-sum/config.yml ./

ENTRYPOINT ["/app/sha256-sum"]