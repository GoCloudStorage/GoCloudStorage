FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct\
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

#将配置文件复制到镜像中
COPY etc /service/etc
#编译后放在service下
RUN go build -ldflags="-s -w" -o /service/storageRpc service/storage/rpc/main.go

#构建小镜像
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
#将二进制文件复制到镜像中的app目录下
COPY --from=builder /service/storageRpc /app/storageRpc
COPY --from=builder /service/etc /app/etc
CMD ["./storageRpc", "-c", "etc/config.yaml"]
