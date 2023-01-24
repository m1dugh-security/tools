package urls

import (
    "github.com/m1dugh/gocrawler/pkg/gocrawler"
    "github.com/m1dugh-security/tools/go/recon-engine/pkg/types"
)

func CrawlProgram(prog *types.ReconedProgram, crawler *gocrawler.Crawler) {
    ch := make(chan types.ReconedUrl)

    crawler.AddCallback(func(res gocrawler.CrawlResponse) {
        ch <- types.ReconedUrl{
            Endpoint: res.URL,
            Status: res.Response.StatusCode,
            ResponseLength: res.Response.ContentLength,
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

