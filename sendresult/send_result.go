package sendresult

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"scraper-bing/config"
	"scraper-bing/getkeyword"
	"scraper-bing/utils"
)

type ResultItem struct {
	URL  string
	Page int
}

type SubmissionData struct {
	URL    string `json:"url"`
	KID    int    `json:"k_id"`
	JID    int    `json:"j_id"`
	QID    int    `json:"q_id"`
	Page   int    `json:"page"`
	Region string `json:"region"`
}

var sendQueue = make(chan []SubmissionData, 400)
var wg sync.WaitGroup

func AddToSendQueue(job getkeyword.KeywordData, results []ResultItem) {
	var submissions []SubmissionData
	for _, result := range results {
		submissions = append(submissions, SubmissionData{
			URL:    result.URL,
			KID:    int(job.KID.Int64),
			JID:    job.JID,
			QID:    job.ID,
			Page:   result.Page,
			Region: job.Region,
		})
	}

	if len(submissions) == 0 {
		fmt.Printf("[QUEUE] ❌ No valid results to enqueue | QID: %d | Query: %s\n", job.ID, job.Query)
		return
	}

	fmt.Printf("[QUEUE] ✅ Enqueue %d items | QID: %d | KID: %d | Region: %s\n",
		len(submissions), job.ID, job.KID.Int64, job.Region)

	wg.Add(1)
	sendQueue <- submissions
}

func sendWorker() {
	for submissions := range sendQueue {
		if hasPageOverLimit(submissions, 200) {
			fmt.Printf("⛔️ Queue stopped (page > 200)\n")
		} else {
			post(submissions)
			utils.PrintDivider()
		}
		wg.Done()
	}
}

func hasPageOverLimit(items []SubmissionData, limit int) bool {
	for _, item := range items {
		if item.Page > limit {
			return true
		}
	}
	return false
}

func StartSendQueue() {
	fmt.Println("🚀 Send Queue Started...")
	go sendWorker()
	go sendWorker()
}

func doRequestWithRetry(client *http.Client, req *http.Request, retries int) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= retries; attempt++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return resp, err
}

func post(results []SubmissionData) {
	if len(results) == 0 {
		fmt.Println("[POST] -> No results to submit.")
		return
	}

	fmt.Printf("[POST] -> Preparing to send %d results\n", len(results))

	jsonData, err := json.Marshal(results)
	if err != nil {
		fmt.Println("[POST] -> JSON Encode Error:", err)
		return
	}
	fmt.Printf("[POST] -> JSON Payload Size: %d bytes\n", len(jsonData))

	url := config.API_URL + config.API_PREFIX + "/job/save-results"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("[POST] -> Request Creation Error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if config.AUTH_TOKEN != "" {
		req.Header.Set("Authorization", "Bearer "+config.AUTH_TOKEN)
	}

	fmt.Println("[POST] -> Sending request to:", url)
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	resp, err := doRequestWithRetry(client, req, 3)
	if err != nil {
		fmt.Println("[POST] -> API Request Error:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[POST] -> Error reading response body:", err)
		return
	}

	status := "error"
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status = "success"
	}

	fmt.Printf("[POST] -> HTTP Status: %d\n", resp.StatusCode)
	for _, result := range results {
		fmt.Printf("[POST] -> %s | Page: %d | KID: %d | Region: %s | URL: %s\n",
			status, result.Page, result.KID, result.Region, utils.TruncateString(result.URL, 40))
	}
	fmt.Printf("[POST] -> Response Body: %s\n", string(bodyBytes))
}
