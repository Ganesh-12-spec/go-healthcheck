package checker

import (
	"net/http"
	"sync"
	"time"
)

type Result struct{
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

func CheckURL(url string) Result {
	start := time.Now()
	resp,err := http.Get(url)
	duration := time.Since(start)
	

	if err != nil{
		return  Result{URL: url, Err:err, Duration:duration}
	}
	defer resp.Body.Close()

	return Result{URL: url, StatusCode: resp.StatusCode, Duration:duration}
}
func CheckALL(urls []string, maxConcurrent int) []Result {
	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(urls))
	sem := make(chan struct{}, maxConcurrent)

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			// hint: sem <- struct{}{} — slot le
			sem <- struct{}{}
			// hint: defer func() { <-sem }() — slot free karo function khatam hote hi
			defer func() {<-sem }()
			result := CheckURL(u)
			resultsChan <- result
		}(url)
	}

	wg.Wait()
	close(resultsChan)

	var results []Result
	for r := range resultsChan {
		results = append(results, r)
	}
	return results
}