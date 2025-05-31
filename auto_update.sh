#!/bin/bash

set -e

cd /var/docker/bing || exit 1

# ✅ อัปเดต Git repo
echo "[AUTO-UPDATE] Pulling latest code..."
git reset --hard HEAD
git pull origin main || {
    echo "[AUTO-UPDATE] Git pull failed"
    exit 1
}

# ✅ ตรวจสอบว่า image มีอยู่หรือไม่
if ! docker image inspect bing-bot-image > /dev/null 2>&1; then
    echo "[AUTO-UPDATE] Image not found. Building new image..."
    docker build -t bing-bot-image .
else
    echo "[AUTO-UPDATE] Image already exists. Skip build."
fi

# ✅ ลบ container เดิม (ถ้ามี)
docker rm -f bing-bot-container 2>/dev/null || true

# ✅ รัน container ใหม่
docker run -d --name bing-bot-container --restart always bing-bot-image

# ✅ เคลียร์ log
find /var/log/ -type f -exec truncate -s 0 {} \;

echo "[AUTO-UPDATE] ✅ Done."
