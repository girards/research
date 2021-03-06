package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
)

// CommonCrawlDotOrg is a source to process subdomains from http://commoncrawl.org
type CommonCrawlDotOrg struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *CommonCrawlDotOrg) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {

	var resultLabel = "commoncrawl"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		uniqFilter := map[string]bool{}

		req, err := http.NewRequest(http.MethodGet, "https://index.commoncrawl.org/CC-MAIN-2018-17-index?url=*."+domain+"&output=json", nil)
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
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New(resp.Status)))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
						return
					}
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

	}(domain, results)
	return results
}
