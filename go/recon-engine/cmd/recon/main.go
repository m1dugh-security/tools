package main

import (
    "fmt"
    "log"
    "github.com/m1dugh-security/tools/go/recon-engine/pkg/recon"
)

func main() {
    options := recon.DefaultOptions()
    options.ThreadsPerProgram = 5
    options.CrawlerConfig.MaxThreads = 5
    options.MaxConcurrentPrograms = 5
    eng, err := recon.New(options)
    if err != nil {
        log.Fatal(err)
    }
    defer eng.Close()
    err = eng.FindPrograms()
    if err != nil {
        log.Fatal(err)
    }

    stages := recon.DefaultStages()
    stages.ScanSubdomains = recon.Never
    stages.Crawl = recon.Never

    _, err = eng.Recon(stages)
    if err != nil {
        log.Fatal(err)
    }

    eng.Close()

    fmt.Println("END")
}
