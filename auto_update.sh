#!/bin/bash

set -e

cd /var/docker/bing || exit 1

echo "[AUTO-UPDATE] 🔄 Pulling latest code from Git..."

# ✅ ตรวจสอบและแก้ permission ให้กับ .git
if [ -d ".git" ]; then
    echo "[AUTO-UPDATE] 🔍 Checking .git permissions..."
    OWNER=$(stat -c "%U" .git)
    if [ "$OWNER" != "$USER" ]; then
        echo "[AUTO-UPDATE] 🛠 Fixing .git ownership to $USER..."
        sudo chown -R "$USER":"$USER" .git
    fi
fi

# ✅ อัปเดต Git repo
git reset --hard HEAD
if ! git pull origin main; then
    echo "[AUTO-UPDATE] ❌ Git pull failed"
    exit 1
fi

# ✅ ตรวจว่า Dockerfile มีการเปลี่ยนแปลง
if [ Dockerfile -nt .docker_image_timestamp ]; then
    echo "[AUTO-UPDATE] 🛠️ Dockerfile changed. Will rebuild image."
    rm -f .docker_image_timestamp
else
    echo "[AUTO-UPDATE] ✅ Dockerfile unchanged. Checking for existing image..."
fi

# ✅ ตรวจสอบว่า image มีอยู่หรือไม่
if ! docker image inspect bing-bot-image > /dev/null 2>&1 || [ ! -f .docker_image_timestamp ]; then
    echo "[AUTO-UPDATE] 🔧 Building new image..."
    docker build -t bing-bot-image .
    touch .docker_image_timestamp
else
    echo "[AUTO-UPDATE] ✅ Docker image already exists."
fi

# ✅ ลบ container เดิม (ถ้ามี)
docker rm -f bing-bot-container 2>/dev/null || true

# ✅ รัน container ใหม่
docker run -d --name bing-bot-container --restart always bing-bot-image

# ✅ เคลียร์ log
find /var/log/ -type f -exec truncate -s 0 {} \;

echo "[AUTO-UPDATE] ✅ Done."
