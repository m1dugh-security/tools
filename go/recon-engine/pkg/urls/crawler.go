package urls

import (
    "github.com/m1dugh/gocrawler/pkg/gocrawler"
    "github.com/m1dugh/recon-engine/pkg/types"
    "net/http"
)

func CrawlProgram(prog *types.ReconedProgram, crawler *gocrawler.Crawler) {
    ch := make(chan types.ReconedUrl)

    crawler.AddCallback(func(res *http.Response, body string) {
        ch <- types.ReconedUrl{
            Endpoint: res.Request.URL.String(),
            Status: res.StatusCode,
            ResponseLength: res.ContentLength,
        }
    })
    go func(ch chan types.ReconedUrl) {
        for val := range ch {
            prog.Urls.AddElement(val)
        }
    }(ch)
    crawler.Crawl(prog.GetUrls())
    close(ch)
}

