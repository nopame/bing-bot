FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ✅ ติดตั้ง Go
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential wget curl unzip ca-certificates fonts-noto \
    && wget -O /tmp/go.tar.gz https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# ✅ ENV
ENV PATH="/usr/local/go/bin:/root/go/bin:$PATH"
ENV HOME=/root
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PWDEBUG=0
ENV DISPLAY=

# ✅ ติดตั้ง playwright CLI เพื่อใช้ดาวน์โหลด driver ใน runtime
RUN npm install -g playwright && \
    playwright install --with-deps

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot

CMD ["/app/bot"]
