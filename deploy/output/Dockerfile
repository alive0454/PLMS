# PLMS 生产镜像（Alpine 轻量版）
FROM alpine:3.18

# 安装必要工具
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制二进制文件
COPY plms-server .
RUN chmod +x plms-server

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./plms-server"]
