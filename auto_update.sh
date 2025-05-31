#!/bin/bash
cd /var/docker/bing || exit 1

# ดึงการอัปเดตจาก Git
git pull origin main

# รีบิวด์ Docker container
docker rm -f bing-bot-container 2>/dev/null || true
docker build -t bing-bot-image .
docker run -d --name bing-bot-container --restart always bing-bot-image

# เคลียร์ log system อย่างปลอดภัย (ไม่ลบไฟล์)
find /var/log/ -type f -exec truncate -s 0 {} \;
