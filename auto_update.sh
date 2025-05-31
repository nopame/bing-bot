#!/bin/bash

set -e

cd /var/docker/bing || exit 1

echo "[AUTO-UPDATE] 🔄 Pulling latest code from Git..."

# ✅ แก้ปัญหา 'dubious ownership'
git config --global --add safe.directory "$(pwd)"

# ✅ ตรวจสอบและแก้ permission ของ .git (ถ้าใช้ sudo หรือ root)
if [ -d ".git" ]; then
    echo "[AUTO-UPDATE] 🛠 Fixing permissions for .git..."
    chown -R "$(id -u):$(id -g)" .git
fi

# ✅ รีเซ็ตและดึงโค้ดจาก Git
git reset --hard HEAD
if ! git pull origin main; then
    echo "[AUTO-UPDATE] ❌ Git pull failed"
    exit 1
fi

# ✅ ตรวจสอบว่ามีการเปลี่ยน Dockerfile หรือไม่
if [ -f .dockerfile_hash ] && cmp -s .dockerfile_hash Dockerfile; then
    echo "[AUTO-UPDATE] ✅ Dockerfile unchanged. Skipping rebuild."
else
    echo "[AUTO-UPDATE] 🛠️ Dockerfile changed. Will rebuild image."
    cp Dockerfile .dockerfile_hash
fi

# ✅ ตรวจสอบว่า image มีอยู่หรือไม่
if ! docker image inspect bing-bot-image > /dev/null 2>&1; then
    echo "[AUTO-UPDATE] 📦 Image not found. Building new image..."
else
    echo "[AUTO-UPDATE] ✅ Docker image already exists."
fi

# ✅ สร้าง image ใหม่ทุกครั้งหากมีการเปลี่ยน Dockerfile
echo "[AUTO-UPDATE] 🔧 Building new image..."
docker build -t bing-bot-image .

# ✅ ลบ container เดิม (ถ้ามี)
docker rm -f bing-bot-container 2>/dev/null || true

# ✅ รัน container ใหม่
docker run -d --name bing-bot-container --restart always bing-bot-image

# ✅ เคลียร์ log system อย่างปลอดภัย
find /var/log/ -type f -exec truncate -s 0 {} \;

echo "[AUTO-UPDATE] ✅ Done."
