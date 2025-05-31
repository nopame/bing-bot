#!/bin/bash

set -e

cd /var/docker/bing || exit 1

echo "[AUTO-UPDATE] 🔄 Pulling latest code from Git..."
git reset --hard HEAD
git pull origin main || {
    echo "[AUTO-UPDATE] ❌ Git pull failed"
    exit 1
}

# เช็คว่า Dockerfile เปลี่ยนแปลงหรือไม่
NEW_HASH=$(md5sum Dockerfile | awk '{print $1}')
OLD_HASH_FILE=".last_dockerfile_hash"

NEED_BUILD=false
if [ ! -f "$OLD_HASH_FILE" ] || [ "$NEW_HASH" != "$(cat $OLD_HASH_FILE)" ]; then
    echo "[AUTO-UPDATE] 🛠️ Dockerfile changed. Will rebuild image."
    NEED_BUILD=true
else
    echo "[AUTO-UPDATE] ✅ Dockerfile unchanged."
fi

# เช็คว่า image มีอยู่หรือไม่
if ! docker image inspect bing-bot-image > /dev/null 2>&1; then
    echo "[AUTO-UPDATE] 🛠️ Docker image not found. Will build image."
    NEED_BUILD=true
else
    echo "[AUTO-UPDATE] ✅ Docker image already exists."
fi

# สร้าง image ถ้าจำเป็น
if [ "$NEED_BUILD" = true ]; then
    echo "[AUTO-UPDATE] 🔧 Building new image..."
    docker build -t bing-bot-image .
    echo "$NEW_HASH" > "$OLD_HASH_FILE"
else
    echo "[AUTO-UPDATE] ⏭️ Skipping build."
fi

# ลบ container เดิมถ้ามี
echo "[AUTO-UPDATE] 🧹 Removing old container (if exists)..."
docker rm -f bing-bot-container 2>/dev/null || true

# รัน container ใหม่
echo "[AUTO-UPDATE] 🚀 Starting new container..."
docker run -d --name bing-bot-container --restart always bing-bot-image

# ล้าง log อย่างปลอดภัย
echo "[AUTO-UPDATE] 🧹 Clearing system logs..."
find /var/log/ -type f -exec truncate -s 0 {} \;

echo "[AUTO-UPDATE] ✅ Update complete."
