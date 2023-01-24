package recon

import (
    "io"
    "log"
    "github.com/m1dugh/gocrawler/pkg/gocrawler"
    datamanager "github.com/m1dugh/recon-engine/internal/database"
    "github.com/m1dugh/recon-engine/pkg/types"
    "github.com/m1dugh/recon-engine/internal/broadcast"
    "github.com/m1dugh/nmapgo/pkg/nmapgo"
)

type StageScope int

const (
    Always StageScope = iota
    IfNew   
    Never
)

type Stages struct {
    ScopeExtraction  StageScope
    FindSubdomains  StageScope
    ScanSubdomains  StageScope
    HttpProbe       StageScope
    FetchRobots     StageScope
    Crawl           StageScope
    Save            StageScope
}

func DefaultStages() *Stages {
    return &Stages{
        ScopeExtraction: Always,
        FindSubdomains: Always,
        ScanSubdomains: IfNew,
        HttpProbe: Always,
        FetchRobots: Always,
        Crawl: IfNew,
        Save: Always,
    }
}

type aggregatedLogger struct {
    errLog *log.Logger
    warnLog *log.Logger
    infoLog *log.Logger
}

func (l *aggregatedLogger) Info(s string, values ...interface{}) {
    l.infoLog.Printf(s + "\n", values...)
}

func (l *aggregatedLogger) Warn(s string, values ...interface{}) {
    l.warnLog.Printf(s + "\n", values...)
}

func (l *aggregatedLogger) Error(s string, values ...interface{}) {
    l.errLog.Printf(s + "\n", values...)
}
type Options struct {
    MaxConcurrentPrograms   uint
    ThreadsPerProgram       uint
    ProbeHttp               bool
    ScannerOptions          *nmapgo.Options
    CrawlerConfig           *gocrawler.Config
    DBConfig                *datamanager.Config
}

func DefaultOptions() *Options {
    return &Options{
        MaxConcurrentPrograms: 5,
        ThreadsPerProgram: 10,
        ProbeHttp: true,
        ScannerOptions: nmapgo.NewOptions(),
        CrawlerConfig: gocrawler.DefaultConfig(),
        DBConfig: datamanager.DefaultConfig(),
    }
}


type ReconEngine struct {
    Options         *Options
    Programs        []*types.ReconedProgram
    masterThrottler *types.ThreadThrottler
    scanner         *nmapgo.Scanner
    DataManager     *datamanager.DataManager
    Logger          *aggregatedLogger
    bot             *broadcast.DiscordBot
    writer          io.WriteCloser
}

