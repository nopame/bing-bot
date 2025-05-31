# ✅ Base Image จาก Playwright ที่มี Ubuntu + Firefox
FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ✅ Set Environment สำหรับ Go Toolchain ไม่ให้ดาวน์โหลดเอง
ARG GOTOOLCHAIN=local
ENV GOTOOLCHAIN=$GOTOOLCHAIN

# ✅ ติดตั้ง Dependency ที่จำเป็น
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget curl unzip git ca-certificates build-essential \
    && rm -rf /var/lib/apt/lists/*

# ✅ ติดตั้ง Playwright CLI สำหรับ Go
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

# ✅ ติดตั้ง Browser ที่ต้องใช้ (เช่น Firefox)
RUN /root/go/bin/playwright install --with-deps firefox

# ✅ ตั้ง Working Directory
WORKDIR /app

# ✅ คัดลอกไฟล์ทั้งหมดเข้าไป
COPY . .

# ✅ Build Go Project โดยใช้ GOTOOLCHAIN=local
RUN export GOTOOLCHAIN=local && go mod tidy && go build -o bot

# ✅ คำสั่งรันเมื่อ container เริ่มต้น
CMD ["./bot"]
