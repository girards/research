package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/subfinder/research/core"
)

// Baidu is a source to process subdomains from https://baidu.com
type Baidu struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Baidu) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {

	var resultLabel = "baidu"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		for currentPage := 1; currentPage <= 750; currentPage++ {
			url := "https://www.baidu.com/s?rn=10&pn=" + strconv.Itoa(currentPage) + "&wd=site%3A" + domain + "+-www.+&oq=site%3A" + domain + "+-www.+"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
				return
			}

			req.WithContext(ctx)

			resp, err := core.HTTPClient.Do(req)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
				return
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New(resp.Status)))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				if ctx.Err() != nil {
					return
				}

				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
						resp.Body.Close()
						return
					}
				}
			}

			resp.Body.Close()

		}

	}(domain, results)
	return results
}
