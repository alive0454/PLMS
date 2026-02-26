# PLMS 生产镜像（使用预编译二进制）
FROM alpine:3.19

# 安装必要工具
RUN apk add --no-cache ca-certificates tzdata wget

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制二进制文件
COPY plms-server .
RUN chmod +x plms-server

# 创建日志目录
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 运行应用
CMD ["./plms-server"]
