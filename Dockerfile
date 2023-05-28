# UPX
FROM golang:1.19.5-alpine AS builder
RUN apk add --update --no-cache \
  build-base \
  upx
WORKDIR /go/src/github.com/hades_api/api/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
  go build -a -installsuffix cgo -ldflags="-s -w" -o web_api ./app/main.go && \
  upx web_api

FROM alpine:latest AS runtime
WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/hades_api/api .
CMD ["./web_api"]