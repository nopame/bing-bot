# ✅ ใช้ base image ของ Microsoft Playwright ที่มี Chromium/Firefox/WebKit ติดตั้งไว้แล้ว
FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ✅ ติดตั้ง Go
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential wget curl unzip ca-certificates fonts-noto \
    && wget -O /tmp/go.tar.gz https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# ✅ กำหนด ENV และ PATH
ENV PATH="/usr/local/go/bin:/root/go/bin:$PATH"
ENV HOME=/root
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PWDEBUG=0
ENV DISPLAY=

# ✅ ติดตั้ง playwright-go-driver แล้ววางไว้ใน ~/.playwright-go/driver
RUN go install github.com/playwright-community/playwright-go/cmd/playwright-go-driver@latest && \
    mkdir -p /root/.playwright-go && \
    cp /root/go/bin/playwright-go-driver /root/.playwright-go/driver

# ✅ เตรียมโปรเจกต์
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ คอมไพล์ Go binary
RUN go build -o bot

# ✅ รัน bot
CMD ["/app/bot"]
