scraper-bing/
│── main.go        # ระบบ Queue + Worker + Fetch งานใหม่อัตโนมัติ
│── search.go      # ฟังก์ชันเปิด Browser และค้นหา Bing
│── get_keyword.go # ดึงคีย์เวิร์ดใหม่เข้าระบบ และป้องกันทำซ้ำ
│── send_result.go # ส่งผลลัพธ์ไปยัง Server (Mock API)
│── go.mod         # ไฟล์ Go Modules สำหรับจัดการ dependencies

--------------------
#build app for window
go build -o bing-bot.exe
--------------------
#Startup Folder
กด Win + R แล้วพิมพ์:
shell:startup
--------------------
# คำสั่ง git หากใช้ครั้งแรก
git init
git remote add origin https://github.com/nopame/bing-bot.git
git branch -M main
git add .
git commit -m "First commit"
git push -u origin main

# git pull ดึงไฟล์
git pull origin main
--------------------
# linux
go build -o search_bing && chmod +x /var/www/app/go/bing/search_bing && chcon -t bin_t /var/www/app/go/bing/search_bing && restorecon -v /var/www/app/go/bing/search_bing && systemctl daemon-reload && systemctl restart bing
--------------------
# Docker
docker rm -f bing-bot-container 2>/dev/null || true && docker build -t bing-bot-image . && docker run -d --name bing-bot-container --restart always bing-bot-image