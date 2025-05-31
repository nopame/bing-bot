# ✅ ใช้ Base Image ที่มี Playwright ติดตั้งพร้อม
FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ✅ ติดตั้ง Go
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget curl unzip git ca-certificates build-essential \
    && wget -O /tmp/go.tar.gz https://go.dev/dl/go1.22.3.linux-amd64.tar.gz \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz && apt-get clean && rm -rf /var/lib/apt/lists/*

# ✅ เพิ่ม PATH
ENV PATH="/usr/local/go/bin:$PATH"

# ✅ ติดตั้ง playwright-go CLI ที่จำเป็น
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

# ✅ ติดตั้ง Firefox สำหรับการใช้งานกับ Playwright
RUN /root/go/bin/playwright install --with-deps firefox

# ✅ กำหนด Environment
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PWDEBUG=0
ENV DISPLAY=

# ✅ ตั้งค่า Working Directory
WORKDIR /app

# ✅ คัดลอกโปรเจกต์เข้าไป
COPY . .

# ✅ ติดตั้ง Module และ Build
RUN go mod tidy && go build -o bot

# ✅ คำสั่งรัน
CMD ["/app/bot"]
