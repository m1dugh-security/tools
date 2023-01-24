package httprobe

import (
    "net/http"
    "time"
    "net"
    "crypto/tls"
    "sync"
    "io/ioutil"
    "io"

    "github.com/m1dugh/recon-engine/pkg/types"
    . "github.com/m1dugh/recon-engine/internal/types"
)

/// This code is greatly inspired by @tomnomnom httprobe cli repo.
/// https://github.com/tomnomnom/httprobe.git

func httpProbeWorker(res *ComparableSet[types.ReconedUrl], client *http.Client, throttler *types.ThreadThrottler, url string, mut *sync.Mutex) {
    if status, length := ping(client, url); status > 0 {
        mut.Lock()
        res.AddElement(types.ReconedUrl{
            Endpoint: url,
            Status: status,
            ResponseLength: length,
        })
        mut.Unlock()
    }
    throttler.Done()
}

func HttpProbe(prog *types.ReconedProgram, httpsOnly bool) {
    var to int = 1
    timeout := time.Duration(to * int(time.Second))
    var tr = &http.Transport{
        MaxIdleConns: 30,
        IdleConnTimeout: time.Second,
        DisableKeepAlives: true,
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
        DialContext: (&net.Dialer{
            Timeout: timeout,
            KeepAlive: time.Second,
        }).DialContext,
    }

    checkRedirect := func(req *http.Request, _ []*http.Request) error {
        return http.ErrUseLastResponse
    }

    client := &http.Client{
        Transport: tr,
        CheckRedirect: checkRedirect,
        Timeout: timeout,
    }

    size := prog.Subdomains.Length()
    if !httpsOnly {
        size *= 2
    }
    var mut sync.Mutex
    for _, v := range *prog.Subdomains {
        prog.Throttler.RequestThread()
        go httpProbeWorker(prog.Urls, client, prog.Throttler, "https://" + v, &mut)
    }

    if !httpsOnly {
        for _, v := range *prog.Subdomains {
            prog.Throttler.RequestThread()
            go httpProbeWorker(prog.Urls, client, prog.Throttler, "http://" + v, &mut)
        }
    }

    prog.Throttler.Wait()
}

// returns status and response length
func ping(client *http.Client, url string) (int, int64) {

    method := "GET"
    req, err := http.NewRequest(method, url, nil)
    if err != nil {
        return -1, 0
    }

    req.Header.Add("Connection", "close")
    req.Close = true

    resp, err := client.Do(req)
    if resp != nil {
        io.Copy(ioutil.Discard, resp.Body)
        resp.Body.Close()
    }

    if err != nil {
        return -1, 0
    }

    return resp.StatusCode, resp.ContentLength
}
