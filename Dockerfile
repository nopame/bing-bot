FROM mcr.microsoft.com/playwright:v1.50.1-jammy

# ✅ ติดตั้ง Go และเครื่องมือที่จำเป็น
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential wget curl unzip git ca-certificates fonts-noto \
    && wget -O /tmp/go.tar.gz https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz && apt-get clean && rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:/root/go/bin:$PATH"
ENV HOME=/root
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PWDEBUG=0
ENV DISPLAY=

# ✅ ดึง playwright-go tag ที่ยังมี driver อยู่ และ build driver
RUN git clone https://github.com/playwright-community/playwright-go.git /tmp/playwright-go && \
    cd /tmp/playwright-go && \
    git checkout tags/v0.170.1 && \
    cd cmd/playwright-go-driver && \
    go build -o /root/.playwright-go/driver .

# ✅ ตั้ง Working Directory
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot

CMD ["/app/bot"]
