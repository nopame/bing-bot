# ✅ ใช้ Base Image ที่รองรับ Playwright แบบ Headless เท่านั้น (ไม่มี GUI)
FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ✅ ติดตั้ง Go
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential wget curl unzip ca-certificates \
    && wget -O /tmp/go.tar.gz https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# ✅ กำหนด PATH ให้ใช้ Go ได้
ENV PATH="/usr/local/go/bin:$PATH"

# ✅ ติดตั้ง Playwright CLI (ใช้ Go Version)
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

# ✅ ติดตั้ง Dependencies ของ Playwright
RUN /root/go/bin/playwright install --with-deps firefox

# ✅ กำหนดให้ Playwright ใช้ Headless Mode เสมอ
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1

# ปิด Debug UI
ENV PWDEBUG=0

# ไม่ใช้ X11
ENV DISPLAY=

# ✅ ตั้งค่า Working Directory
WORKDIR /app

# ✅ คัดลอกโค้ดโปรเจคเข้า Docker
COPY . /app

# ✅ คอมไพล์โค้ด
RUN go mod tidy && go build -o bot

# ✅ คำสั่งเริ่มต้นรัน Bot
CMD ["/app/bot"]
