# ✅ ใช้ base image ของ Playwright ที่มี Chromium/Firefox/WebKit ติดตั้งแล้ว
FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ✅ ติดตั้ง Go
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential wget curl unzip ca-certificates fonts-noto \
    && wget -O /tmp/go.tar.gz https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# ✅ กำหนด PATH และ ENV ที่จำเป็น
ENV PATH="/usr/local/go/bin:/root/go/bin:$PATH"
ENV HOME=/root
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PWDEBUG=0
ENV DISPLAY=

# ✅ ติดตั้ง Playwright driver สำหรับ Go
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest && \
    /root/go/bin/playwright install --with-deps

# ✅ ตั้งค่า Working Directory
WORKDIR /app

# ✅ ใช้ Docker layer cache สำหรับ dependency
COPY go.mod go.sum ./
RUN go mod download

# ✅ คัดลอกโค้ดโปรเจกต์เข้า
COPY . .

# ✅ คอมไพล์ Go Binary
RUN go build -o bot

# ✅ คำสั่งรันหลัก
CMD ["/app/bot"]
