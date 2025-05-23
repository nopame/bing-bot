package getkeyword

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"crypto/tls"
	"net/http"
	"sync"

	"database/sql"
	"scraper-bing/config"
	"scraper-bing/utils"
)

// ✅ โครงสร้างของ Keyword งาน
type KeywordData struct {
	ID       int           `json:"id"` // ✅ ใช้ `q_id` เป็น `ID`
	Query    string        `json:"query"`
	Lang     string        `json:"lang"`
	Region   string        `json:"region"`
	Priority sql.NullInt64 `json:"priority"`
	Status   string        `json:"status"`
	KID      sql.NullInt64 `json:"k_id"`
	SetID    sql.NullInt64 `json:"set_id"`
	JID      int           `json:"j_id"`
}

// ✅ ใช้ Mutex ป้องกันการดึงงานซ้ำซ้อน
var queueLock sync.Mutex

// ✅ ฟังก์ชันนับจำนวน Active Jobs ใน `sync.Map`
func CountActiveJobs(activeJobs *sync.Map) int {
	count := 0
	activeJobs.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// ✅ **ฟังก์ชันดึงงานจาก Server**
func getServerJob() (*KeywordData, error) {
	url := config.API_URL + config.API_PREFIX +"/job/keyword"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // true for self-signed
			},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.AUTH_TOKEN)

	// ✅ Debug Log
	fmt.Println("🔄 Fetching jobs from:", url)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch job from server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("❌ Server returned status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to read response body: %v", err)
	}

	// ✅ แปลง JSON เป็น `map[string]interface{}`
	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return nil, fmt.Errorf("❌ Failed to parse JSON: %v", err)
	}

	// ✅ ตรวจสอบค่าก่อนใช้ ป้องกัน panic
	job := KeywordData{
		ID:       getInt(jsonData, "q_id"), // ✅ ใช้ `q_id` เป็น `ID`
		Query:    getString(jsonData, "query"),
		Lang:     getString(jsonData, "lang"),
		Region:   getString(jsonData, "region"),
		Status:   getString(jsonData, "status"),
		Priority: toNullInt64(getInt(jsonData, "priority")),
		KID:      toNullInt64(getInt(jsonData, "k_id")),
		SetID:    toNullInt64(getInt(jsonData, "set_id")),
		JID:      getInt(jsonData, "j_id"),
	}

	fmt.Println("✅ Received Job:", job)
	return &job, nil
}

// ✅ **ฟังก์ชันช่วยดึงค่าจาก `map[string]interface{}`**
func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0 // ✅ ถ้าไม่มีค่าให้คืน 0
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return "" // ✅ ถ้าไม่มีค่าให้คืน ""
}

// ✅ **ฟังก์ชันแปลง `int` เป็น `sql.NullInt64`**
func toNullInt64(value int) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{Int64: 0, Valid: false}
	}
	return sql.NullInt64{Int64: int64(value), Valid: true}
}

// ✅ **ฟังก์ชันดึง Keywords เข้า Queue**
func FetchKeywords(queue chan KeywordData, activeJobs *sync.Map) {
	queueLock.Lock()
	defer queueLock.Unlock()

	fmt.Printf("📥 Queue: %d/%d | Active: %d\n", len(queue), config.QueueSize, CountActiveJobs(activeJobs))

	job, err := getServerJob()
	if err != nil {
		fmt.Println("❌ Error fetching job:", err)
		return
	}

	// ✅ ตรวจสอบว่า `job` มีค่าและยังไม่มีใน ActiveJobs
	if job != nil {
		if _, exists := activeJobs.Load(job.ID); exists {
			fmt.Println("⚠️ Job already in progress, skipping:", job.ID)
			return
		}

		queue <- *job
		activeJobs.Store(job.ID, true)
		fmt.Println("✅ Added Job to Queue:", job.Query, "(ID:", job.ID, ")")
	}

	fmt.Println("✅ Fetching Completed.")
	utils.PrintDivider()
}
