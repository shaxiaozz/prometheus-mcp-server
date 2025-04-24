FROM golang:1.24.2-alpine as builder
WORKDIR /data/prometheus-mcp-server-code
ENV GOPROXY=https://goproxy.cn
RUN apk add --no-cache upx ca-certificates tzdata
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o prometheus-mcp-server

FROM golang:1.24.2-alpine as runner
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /data/prometheus-mcp-server-code/prometheus-mcp-server /prometheus-mcp-server
EXPOSE 8000
CMD ["/prometheus-mcp-server"]