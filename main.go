package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"scraper-bing/config"
	"scraper-bing/getkeyword"
	"scraper-bing/search"
	"scraper-bing/utils"
)

func ensurePlaywrightDriverInstalled() {
	err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"firefox"},
	})
	if err != nil {
		log.Fatalf("❌ Failed to install Playwright driver: %v", err)
	}
	log.Println("✅ Playwright driver installed successfully.")
}

func main() {
	utils.PrintDivider()
	fmt.Printf("\U0001F680 Starting Bing Scraper | Max Workers: %d | Queue Size: %d\n", config.MaxConcurrentJobs, config.QueueSize)

	// ✅ ตรวจสอบว่ามี Driver ของ Playwright แล้วหรือยัง
	ensurePlaywrightDriverInstalled()

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
				case job := <-queue:
					fmt.Printf("✅ Worker %d received job: %s (ID: %d)\n", workerID, job.Query, job.ID)
					<-workerPool

					fmt.Printf("🛠️ Worker %d processing: %s (ID: %d)\n", workerID, job.Query, job.ID)
					if err := search.SearchBing(&job); err != nil {
						log.Printf("❌ Worker %d failed: %v\n", workerID, err)
					}

					activeJobs.Delete(job.ID)
					workerPool <- struct{}{}

				default:
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("🎉 All tasks completed!")
	utils.PrintDivider()
}
