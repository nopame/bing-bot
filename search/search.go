package search

import (
	"fmt"
	"net/url"

	"github.com/playwright-community/playwright-go"
	"scraper-bing/sendresult"
	"scraper-bing/getkeyword"
	"scraper-bing/config"
	"scraper-bing/utils"
)

var browser playwright.Browser
var context playwright.BrowserContext

func InitBrowser() error {
	fmt.Println("🚀 Initializing Playwright...")
	sendresult.StartSendQueue()

	if !config.OpenBrowser {
		fmt.Println("🛑 Browser mode disabled (config.OpenBrowser = false)")
		return nil
	}

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("❌ Failed to start Playwright: %v", err)
	}

	fmt.Println("🔥 Launching Firefox with UI...")
	browser, err = pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to launch Firefox: %v", err)
	}

	context, err = browser.NewContext()
	if err != nil {
		return fmt.Errorf("❌ Failed to create context: %v", err)
	}

	context.Route("**/*", func(route playwright.Route) {
		request := route.Request()
		resourceType := request.ResourceType()
		url := request.URL()

		if config.ResourceBlockLog {
			fmt.Printf("🔄 Request: %s | URL: %s\n", resourceType, url)
		}

		if resourceType == "image" || resourceType == "stylesheet" ||
			resourceType == "font" || resourceType == "media" ||
			resourceType == "video" || url[:5] == "data:" {
			if config.ResourceBlockLog {
				fmt.Printf("🚫 Blocking: %s\n", url)
			}
			route.Abort()
			return
		}

		route.Continue()
	})

	context.SetDefaultNavigationTimeout(0)
	fmt.Println("-> Firefox launched successfully with resource blocking enabled!")

	return nil
}

func CloseBrowser() {
	if browser != nil {
		fmt.Println("🛑 Closing Browser...")
		browser.Close()
	}
}

func ClickNextPage(page playwright.Page) bool {
	nextButton := page.Locator("a[aria-label='Next page']")
	count, err := nextButton.Count()
	if err != nil {
		fmt.Println("-> Error checking Next button:", err)
		return false
	}
	if count == 0 {
		fmt.Println("-> No more pages to process.")
		return false
	}

	if err := nextButton.Click(); err != nil {
		fmt.Println("-> Failed to click Next button.")
		return false
	}

	state := playwright.LoadState("networkidle")
	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:  &state,
		Timeout: playwright.Float(8000),
	})

	return true
}

func SearchBing(job *getkeyword.KeywordData) error {
	fmt.Printf("[SEARCH] -> Searching Bing : %s (QID: %d)\n", job.Query, job.ID)
	utils.PrintDivider()

	searchURL := fmt.Sprintf("https://www.bing.com/search?q=%s&setlang=%s&cc=%s",
		url.QueryEscape(job.Query), job.Lang, job.Region)

	if !config.OpenBrowser {
		fmt.Println("[SEARCH] -> Running in Headless Mode...")
		pw, err := playwright.Run()
		if err != nil {
			return fmt.Errorf("[SEARCH] -> Failed to start Playwright: %v\n", err)
		}
		defer pw.Stop()

		browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("[SEARCH] -> Failed to launch headless browser: %v\n", err)
		}
		defer browser.Close()

		context, err := browser.NewContext()
		if err != nil {
			return fmt.Errorf("[SEARCH] -> Failed to create browser context: %v\n", err)
		}
		return ProcessSearchResults(job, searchURL, context)
	}

	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("[SEARCH] -> Failed to create new page: %v\n", err)
	}
	defer page.Close()

	if _, err := page.Goto(searchURL); err != nil {
		return fmt.Errorf("[SEARCH] -> Failed to open Bing: %v\n", err)
	}

	state := playwright.LoadState("networkidle")
	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:  &state,
		Timeout: playwright.Float(8000),
	})

	return ProcessSearchResults(job, searchURL, context)
}

func ProcessSearchResults(job *getkeyword.KeywordData, searchURL string, context playwright.BrowserContext) error {
	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("[SEARCH] ❌ Failed to create new page: %v\n", err)
	}
	defer page.Close()

	if _, err := page.Goto(searchURL); err != nil {
		return fmt.Errorf("[SEARCH] ❌ Failed to open Bing: %v\n", err)
	}

	state := playwright.LoadState("networkidle")
	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:  &state,
		Timeout: playwright.Float(8000),
	})

	pageNum := 1
	for {
		_, err := page.WaitForSelector("h2 a", playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(5000),
		})
		if err != nil {
			fmt.Printf("[SEARCH] ⚠️ No <h2 a> found on page %d (timeout)\n", pageNum)
		}

		elements, err := page.Locator("h2 a").All()
		if err != nil {
			return fmt.Errorf("[SEARCH] ❌ Failed to extract links: %v\n", err)
		}

		var results []sendresult.ResultItem
		for _, element := range elements {
			href, err := element.GetAttribute("href")
			if err == nil && href != "" {
				results = append(results, sendresult.ResultItem{
					URL:  href,
					Page: pageNum,
				})
			}
		}

		if len(results) == 0 {
			fmt.Printf("[SEARCH] ⚠️ No results found on page %d for query: %s\n", pageNum, job.Query)
		}

		sendresult.AddToSendQueue(*job, results)

		if !ClickNextPage(page) {
			job.Status = "completed"
			break
		}

		page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State:  &state,
			Timeout: playwright.Float(8000),
		})

		pageNum++
	}

	fmt.Printf("[SEARCH] ✅ Finished processing: %s (QID: %d) | Status: %s\n", job.Query, job.ID, job.Status)
	return nil
}
