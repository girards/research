package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestArchiveIs(t *testing.T) {
	domain := "apple.com"
	source := ArchiveIs{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// stop after 20
	counter := 0

	for result := range source.ProcessDomain(ctx, domain) {
		counter++
		//if counter == 20 {
		//	cancel()
		//}
		fmt.Println(result.Success)
		//t.Log(result.Success)
	}

	fmt.Printf("found '%v' results\n", counter)
}

func TestArchiveIsRecursive(t *testing.T) {
	domain := "apple.com"
	source := &ArchiveIs{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := &core.EnumerationOptions{
		Recursive: true,
		Sources:   []core.Source{source},
	}

	for result := range core.EnumerateSubdomains(ctx, domain, options) {
		results = append(results, result)
		fmt.Println(result)

	}

	if !(len(results) >= 5) {
		t.Errorf("expected more than 5 result(s), got '%v'", len(results))
		t.Error(ctx.Err())
	}

	fmt.Println(len(results), ctx.Err())
}

// TODO: fix tests to add the new context version of the API

//func TestArchiveIsMultiThreaded(t *testing.T) {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := ArchiveIs{}
//	results := []*core.Result{}
//
//	wg := sync.WaitGroup{}
//	mx := sync.Mutex{}
//
//	for _, domain := range domains {
//		wg.Add(1)
//		go func(domain string) {
//			defer wg.Done()
//			for result := range source.ProcessDomain(domain) {
//				t.Log(result)
//				if result.IsSuccess() && result.IsFailure() {
//					t.Error("got a result that was a success and failure")
//				}
//				mx.Lock()
//				results = append(results, result)
//				mx.Unlock()
//			}
//		}(domain)
//	}
//
//	wg.Wait() // collect results
//
//	if len(results) <= 4 {
//		t.Errorf("expected at least 4 results, got '%v'", len(results))
//	}
//}
//
//func ExampleArchiveIs() {
//	domain := "bing.com"
//	source := ArchiveIs{}
//	results := []*core.Result{}
//
//	for result := range source.ProcessDomain(domain) {
//		results = append(results, result)
//	}
//
//	fmt.Println(len(results) >= 20)
//	// Output: true
//}
//
//func ExampleArchiveIs_multi_threaded() {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := ArchiveIs{}
//	results := []*core.Result{}
//
//	wg := sync.WaitGroup{}
//	mx := sync.Mutex{}
//
//	for _, domain := range domains {
//		wg.Add(1)
//		go func(domain string) {
//			defer wg.Done()
//			for result := range source.ProcessDomain(domain) {
//				mx.Lock()
//				results = append(results, result)
//				mx.Unlock()
//			}
//		}(domain)
//	}
//
//	wg.Wait() // collect results
//
//	fmt.Println(len(results) >= 4)
//	// Output: true
//}
//
//func BenchmarkArchiveIsSingleThreaded(b *testing.B) {
//	domain := "bing.com"
//	source := ArchiveIs{}
//
//	for n := 0; n < b.N; n++ {
//		results := []*core.Result{}
//		for result := range source.ProcessDomain(domain) {
//			results = append(results, result)
//		}
//	}
//}
//
//func BenchmarkArchiveIsMultiThreaded(b *testing.B) {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := ArchiveIs{}
//	wg := sync.WaitGroup{}
//	mx := sync.Mutex{}
//
//	for n := 0; n < b.N; n++ {
//		results := []*core.Result{}
//
//		for _, domain := range domains {
//			wg.Add(1)
//			go func(domain string) {
//				defer wg.Done()
//				for result := range source.ProcessDomain(domain) {
//					mx.Lock()
//					results = append(results, result)
//					mx.Unlock()
//				}
//			}(domain)
//		}
//
//		wg.Wait() // collect results
//	}
//}
