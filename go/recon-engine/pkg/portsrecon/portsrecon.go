package portsrecon

import (
    "github.com/m1dugh-security/tools/go/recon-engine/pkg/types"
    "github.com/m1dugh/nmapgo/pkg/nmapgo"
    "sync"
)

func scanSubdomainWorker(prog *types.ReconedProgram, scanner *nmapgo.Scanner,
url string, mut *sync.Mutex) {
    defer prog.Throttler.Done()
    host, err := scanner.ScanHost(url)
    if err != nil {
        return
    }

    if host != nil {
        mut.Lock()
        prog.Hosts = append(prog.Hosts, host)
        mut.Unlock()
    }
}

func ScanSubdomains(scanner *nmapgo.Scanner, prog *types.ReconedProgram) error {
    var mut sync.Mutex
    prog.Hosts = make([]*nmapgo.Host, 0, prog.Subdomains.Length())
    for _, url := range *prog.Subdomains {
        prog.Throttler.RequestThread()
        go scanSubdomainWorker(prog, scanner, url, &mut)
    }
    prog.Throttler.Wait()
    return nil
}

