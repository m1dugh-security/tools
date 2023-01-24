package recon

import (
    "log"
    "os"
    "io"
    "fmt"
    "github.com/m1dugh/recon-engine/pkg/programs"
    "github.com/m1dugh/recon-engine/internal/broadcast"
    datamanager "github.com/m1dugh/recon-engine/internal/database"
    "github.com/m1dugh/recon-engine/pkg/subdomains"
    "github.com/m1dugh/recon-engine/pkg/portsrecon"
    "github.com/m1dugh/recon-engine/pkg/httprobe"
    "github.com/m1dugh/recon-engine/pkg/urls"
    "github.com/m1dugh/recon-engine/pkg/types"
    ptypes "github.com/m1dugh/recon-engine/pkg/programs/types"
    "github.com/m1dugh/nmapgo/pkg/nmapgo"
    "github.com/m1dugh/gocrawler/pkg/gocrawler"
    "errors"
    "regexp"
)



func New(options *Options) (*ReconEngine, error) {
    if options == nil {
        options = DefaultOptions()
    }

    pipeReader, pipeWriter := io.Pipe()
    writer := io.MultiWriter(os.Stderr, pipeWriter)
    logger := &aggregatedLogger{
        errLog: log.New(writer, "[ERR] ", log.Ldate | log.Ltime),
        warnLog: log.New(writer, "[WARN] ", log.Ldate | log.Ltime),
        infoLog: log.New(writer, "[INFO] ", log.Ldate | log.Ltime),
    }

    token := os.Getenv("DISCORD_TOKEN")
    bot, err := broadcast.NewDiscord(token)
    bot.Start()
    bot.SetReader(pipeReader)
    if err != nil {
        msg := fmt.Sprintf("Error when creating discord bot %s\n", err)
        return nil, errors.New(msg)
    }

    masterThrottler := types.NewThreadThrottler(options.MaxConcurrentPrograms)

    scanner, err := nmapgo.NewScanner(options.ScannerOptions)
    if err != nil {
        msg := fmt.Sprintf("Error while loading nmapgo scanner: %s\n", err)
        logger.Error(msg)
        return nil, errors.New(msg)
    }


    res := &ReconEngine{
        Options: options,
        Programs: nil,
        masterThrottler: masterThrottler,
        scanner: scanner,
        DataManager: datamanager.New(options.DBConfig),
        Logger: logger,
        bot: bot,
        writer: pipeWriter,
    }
    err = res.DataManager.Init()
    if err != nil {
        msg := fmt.Sprintf("Error while starting data manager: %s", err)
        logger.Error(msg)
        return nil, errors.New(msg)
    }

    return res, nil
}

var _webRegex *regexp.Regexp = regexp.MustCompile(`web|api|http`)

func (eng *ReconEngine) FindPrograms() error {
    eng.Logger.Info("Fetching programs")
    programs, err := programs.GetPrograms(nil)
    if err != nil {
        eng.Logger.Error("%s", err)
        return errors.New("ReconEngine.FindPrograms: could not fetch programs")
    }

    for _, prog := range programs {
        prog.Throttler.MaxThreads = eng.Options.ThreadsPerProgram
        eng.Logger.Info("Extracting scope info for %s", prog.Program.Code())
        prog.ExtractScopeInfo()
    }

    eng.Programs = programs
    return nil
}

func fetchSubdomainsWorker(prog *types.ReconedProgram,
throttler *types.ThreadThrottler) {
    defer throttler.Done()
    err := subdomains.FetchSubdomains(prog)
    if err != nil {
        return
    }
}

func (eng *ReconEngine) FetchSubdomains() {
    for _, prog := range eng.Programs {
        eng.masterThrottler.RequestThread()
        eng.Logger.Info("Fetching subdomains for %s", prog.Program.Code())
        go fetchSubdomainsWorker(prog, eng.masterThrottler)
    }

    eng.masterThrottler.Wait()
    eng.Logger.Info("Finished fetching subdomains")
}

func (eng *ReconEngine) HttpProbe() {
    for _, prog := range eng.Programs {
        eng.masterThrottler.RequestThread()
        go func(throttler *types.ThreadThrottler,
            prog *types.ReconedProgram,
            httpsOnly bool,
        ) {
            defer throttler.Done()
            eng.Logger.Info("Probing %s", prog.Program.Code())
            httprobe.HttpProbe(prog, httpsOnly)
        } (eng.masterThrottler, prog, !eng.Options.ProbeHttp)
    }

    eng.masterThrottler.Wait()
    eng.Logger.Info("Finished probing programs")
}

func (eng *ReconEngine) ScanSubdomains() error {
    if eng.scanner == nil {
        return errors.New("ReconEngine.NewReconEngine: could not create nmapgo.Scanner")
    }

    for _, prog := range eng.Programs {
        eng.masterThrottler.RequestThread()
        go func(throttler *types.ThreadThrottler, prog *types.ReconedProgram,
            scanner *nmapgo.Scanner,
        ) {
            defer throttler.Done()
            eng.Logger.Info("Scanning subdomains for %s", prog.Program.Code())
            portsrecon.ScanSubdomains(scanner, prog)
        }(eng.masterThrottler, prog, eng.scanner)
    }

    eng.masterThrottler.Wait()
    eng.Logger.Info("Finished scanning subdomains")
    return nil
}

func (eng *ReconEngine) FetchRobots() {
    for _, prog := range eng.Programs {
        eng.masterThrottler.RequestThread()
        go func(throttler *types.ThreadThrottler, prog *types.ReconedProgram) {
            eng.Logger.Info("Fetching robots for %s", prog.Program.Code())
            urls.FetchUrls(prog)
            throttler.Done()
        }(eng.masterThrottler, prog)
    }
    eng.masterThrottler.Wait()
    eng.Logger.Info("Finished fetching robots")
}


func convertScope(scope *ptypes.Scope) *gocrawler.Scope {
    return gocrawler.NewScope(scope.Include, scope.Exclude)
}

func (eng *ReconEngine) CrawlPages() {
    for _, prog := range eng.Programs {
        scope := prog.Program.GetScope(_webRegex)
        cr:= gocrawler.New(convertScope(scope), eng.Options.CrawlerConfig)
        eng.masterThrottler.RequestThread()
        go func (throttler *types.ThreadThrottler, prog *types.ReconedProgram, crawler *gocrawler.Crawler) {
            defer throttler.Done()
            eng.Logger.Info("Starting crawling for %s", prog.Program.Code())
            urls.CrawlProgram(prog, crawler)
        }(eng.masterThrottler, prog, cr)
    }

    eng.masterThrottler.Wait()
    eng.Logger.Info("Finished crawling")
}

func (eng *ReconEngine) SavePrograms() ([]datamanager.ProgramDiff, error) {
    res := make([]datamanager.ProgramDiff, 0, len(eng.Programs))
    for _, prog := range eng.Programs {
        eng.Logger.Info("Saving %s", prog.Program.Code())
        diff, err := eng.DataManager.SaveProgram(prog)
        if err != nil {
            return nil, err
        }
        if !diff.IsEmpty() {
            err = eng.bot.SendDiffs(prog.Program.Code(), diff)
            if err != nil {
                eng.Logger.Error("%s", err)
                // log.Fatal(err)
            }
            res = append(res, diff)
        }
    }
    return res, nil
}


func (eng *ReconEngine) reconProgram(prog *types.ReconedProgram, stages *Stages) (datamanager.ProgramDiff, error) {

    eng.Logger.Info("[START] Starting recon for %s", prog.Program.Code())
    exists, err := eng.DataManager.ProgramExists(prog.Program)
    if err != nil {
        eng.Logger.Error("%s", err)
        return datamanager.ProgramDiff{}, err
    }

    _useStage := func(val StageScope) bool {
        return val == Always || (val == IfNew && !exists)
    }

    if _useStage(stages.ScopeExtraction) {
        eng.Logger.Info("Extraction Scope info for %s", prog.Program.Code())
        prog.ExtractScopeInfo()
    }

    if _useStage(stages.FindSubdomains) {
        eng.Logger.Info("Fetching subdomains for %s", prog.Program.Code())
        err := subdomains.FetchSubdomains(prog)
        if err != nil {
            return datamanager.ProgramDiff{}, errors.New(fmt.Sprintf("Error when fetching subdomains: %s", err))
        }
    }

    if _useStage(stages.ScanSubdomains) {
        eng.Logger.Info("Scanning subdomains for %s", prog.Program.Code())
        portsrecon.ScanSubdomains(eng.scanner, prog)
        eng.Logger.Info("Finished scanning subdomains for %s", prog.Program.Code())
    }

    if _useStage(stages.HttpProbe) {
        eng.Logger.Info("Probing %s", prog.Program.Code())
        httprobe.HttpProbe(prog, !eng.Options.ProbeHttp)
        eng.Logger.Info("Finished probing %s", prog.Program.Code())
    }

    if _useStage(stages.FetchRobots) {
        eng.Logger.Info("Fetching robots for %s", prog.Program.Code())
        urls.FetchUrls(prog)
        eng.Logger.Info("Finished fetching robots for %s", prog.Program.Code())
    }

    if _useStage(stages.Crawl) {
        eng.Logger.Info("Starting crawling for %s", prog.Program.Code())
        scope := prog.Program.GetScope(_webRegex)
        cr:= gocrawler.New(convertScope(scope), eng.Options.CrawlerConfig)
        urls.CrawlProgram(prog, cr)
        eng.Logger.Info("Finished crawling for %s", prog.Program.Code())
    }

    if _useStage(stages.Save) {
        eng.Logger.Info("Saving %s", prog.Program.Code())
        diff, err := eng.DataManager.SaveProgram(prog)
        if err != nil {
            return datamanager.ProgramDiff{}, err
        }

        eng.Logger.Info("Saved %s", prog.Program.Code())
        if !diff.IsEmpty() {

            eng.Logger.Info("Sending diffs for %s", prog.Program.Code())
            err = eng.bot.SendDiffs(prog.Program.Code(), diff)
            if err != nil {
                eng.Logger.Error("Discord bot error while sending diffs: %s", err)
                return datamanager.ProgramDiff{}, nil
            }
            eng.Logger.Info("Sent diffs for %s", prog.Program.Code())
            return diff, nil
        }
    }

    eng.Logger.Info("[END] Recon for %s done !", prog.Program.Code())
    return datamanager.ProgramDiff{}, nil
}

func (eng *ReconEngine) reconProgramWorker(prog *types.ReconedProgram, stages *Stages, ch chan datamanager.ProgramDiff) {
    diff, err := eng.reconProgram(prog, stages)
    eng.masterThrottler.Done()
    if err != nil {
        eng.Logger.Error("%s", err)
        return
    }
    ch <- diff
}

func (eng *ReconEngine) Recon(stages *Stages) ([]datamanager.ProgramDiff, error) {
    eng.Logger.Info("Starting recon")
    ch := make(chan datamanager.ProgramDiff)
    for _, prog := range eng.Programs {
        eng.masterThrottler.RequestThread()
        go eng.reconProgramWorker(prog, stages, ch)
    }

    var res []datamanager.ProgramDiff = make([]datamanager.ProgramDiff, 0)
    go func(r *[]datamanager.ProgramDiff) {
        for v := range ch {
            *r = append(*r, v)
        }
    }(&res)
    eng.masterThrottler.Wait()
    close(ch)

    eng.Logger.Info("Finished recon")
    return res, nil
}

func (eng *ReconEngine) Close() {
    eng.writer.Close()
    eng.DataManager.Close()
    eng.bot.Close()
}

