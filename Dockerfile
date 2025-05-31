FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ติดตั้ง Go
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential wget curl unzip ca-certificates fonts-noto \
    && wget -O /tmp/go.tar.gz https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz && apt-get clean && rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:$PATH"

# Playwright ใช้ browser จาก base image นี้
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PWDEBUG=0
ENV DISPLAY=

WORKDIR /app

# ใช้ layer cache ให้ดีขึ้น
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot

CMD ["/app/bot"]
