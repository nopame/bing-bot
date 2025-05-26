#!/bin/bash
cd /var/docker/bing || exit 1

# ดึงอัปเดตจาก Git
git pull origin main

# ลบ container เดิม (ถ้ามี) แล้ว build และ run ใหม่
docker rm -f bing-bot-container 2>/dev/null || true
docker build -t bing-bot-image .
docker run -d --name bing-bot-container --restart always bing-bot-image
