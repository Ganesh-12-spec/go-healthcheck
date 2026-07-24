package checker

import (
	"net/http"
	"sync"
	"time"
)

type Result struct{
	URl        string
  StatusCode int
	Duration   time.Duration
	Err        error
}

func CheckURl(url string) Result {
	start := time.Now()
	resp,err := http.Get(url)
	duration := time.Since(start)
	

	if err != nil{
		return  Result{URl :url,Err:err,Duration:duration}
	}
	defer resp.Body.Close()

	return Result{URl: url, StatusCode: resp.StatusCode, Duration:duration}
}
func CheckALL(urls []string) []Result {
	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(urls))

	for _,url := range urls {
		wg.Add(1)
		go func(u string){
			defer wg.Done()
			result := CheckURl(u)
			resultsChan <- result
		}(url)
	}
	wg.Wait()
  close(resultsChan)

	var results []Result
	for r:= range resultsChan{
		results = append(results, r)

	}
	return results

}