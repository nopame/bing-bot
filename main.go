package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"scraper-bing/config"
	"scraper-bing/getkeyword"
	"scraper-bing/search"
	"scraper-bing/utils"
)

func main() {
	// ✅ กำหนด ENV ให้ playwright-go ใช้ browser ที่ติดตั้งไว้ใน container
	os.Setenv("PLAYWRIGHT_BROWSERS_PATH", "/ms-playwright")
	os.Setenv("PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD", "1")

	utils.PrintDivider()
	fmt.Printf("\U0001F680 Starting Bing Scraper | Max Workers: %d | Queue Size: %d\n", config.MaxConcurrentJobs, config.QueueSize)

	// ✅ เปิด Browser ถ้าถูกเปิดใช้งาน
	err := search.InitBrowser()
	if err != nil {
		log.Fatalf("❌ Failed to initialize browser: %v", err)
	}
	defer search.CloseBrowser()

	// ✅ สร้าง Queue และ ActiveJobs
	queue := make(chan getkeyword.KeywordData, config.QueueSize)
	activeJobs := &sync.Map{}
	var wg sync.WaitGroup

	// ✅ Fetch Jobs ให้ Queue เต็มล่วงหน้า
	go func() {
		for {
			fmt.Println("[Trigger] 🔄 FetchKeywords()...")
			getkeyword.FetchKeywords(queue, activeJobs)
			time.Sleep(2 * time.Second)
		}
	}()

	// ✅ ใช้ Worker Pool บริหาร Worker
	workerPool := make(chan struct{}, config.MaxConcurrentJobs)

	// ✅ สร้าง Worker ตามจำนวนที่กำหนด
	for i := 0; i < config.MaxConcurrentJobs; i++ {
		workerPool <- struct{}{}
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case job := <-queue: // ✅ ดึงงานจาก Queue
					fmt.Printf("✅ Worker %d received job: %s (ID: %d)\n", workerID, job.Query, job.ID)
					<-workerPool // ✅ รับสิทธิ์ทำงานจาก Worker Pool

					fmt.Printf("🛠️ Worker %d processing: %s (ID: %d)\n", workerID, job.Query, job.ID)
					if err := search.SearchBing(&job); err != nil {
						log.Printf("❌ Worker %d failed: %v\n", workerID, err)
					}

					// ✅ เมื่องานเสร็จ ลบออกจาก ActiveJobs
					activeJobs.Delete(job.ID)

					// ✅ คืน Worker Slot กลับไป
					workerPool <- struct{}{}

				default:
					time.Sleep(500 * time.Millisecond) // ✅ รอรับงานใหม่ ไม่ใช้ CPU 100%
				}
			}
		}(i)
	}

	// ✅ รอให้ Worker ทำงานเสร็จ
	wg.Wait()
	fmt.Println("🎉 All tasks completed!")
	utils.PrintDivider()
}
