#!/bin/bash

# 生成自签名 SSL 证书（仅用于测试）

CERT_DIR="$(dirname "$0")/../nginx/ssl"
mkdir -p "$CERT_DIR"

echo "生成自签名证书..."

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$CERT_DIR/key.pem" \
    -out "$CERT_DIR/cert.pem" \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=PLMS/OU=IT/CN=localhost"

chmod 600 "$CERT_DIR/key.pem"
chmod 644 "$CERT_DIR/cert.pem"

echo "证书已生成:"
echo "  证书: $CERT_DIR/cert.pem"
echo "  私钥: $CERT_DIR/key.pem"
echo ""
echo "注意: 自签名证书会被浏览器标记为不安全，仅用于测试！"
echo "生产环境请使用 Let's Encrypt 或正规 CA 签发的证书。"
