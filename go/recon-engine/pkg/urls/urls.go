package urls

import (
    "github.com/m1dugh-security/tools/go/recon-engine/pkg/types"
    "github.com/m1dugh-security/tools/go/utils/pkg/utils"
    "net/http"
    "regexp"
    "sync"
    "bufio"
    "strings"
)

var (
    urlExtractor *regexp.Regexp = regexp.MustCompile(`^https?://([\w\-]+\.)[a-z]{2,7}$`)
)

func fetchUrlsWorker(root string,
    client *http.Client,
    res *utils.StringSet,
    throttler *types.ThreadThrottler,
    mut *sync.Mutex,
    urls *utils.ComparableSet[types.ReconedUrl],
) {
    defer throttler.Done()
    resp, err := client.Get(root + "/robots.txt")
    if err != nil {
        return
    }


    var found *utils.StringSet = utils.NewStringSet(nil)
    scanner := bufio.NewScanner(resp.Body)
    defer resp.Body.Close()
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "Allow") || strings.HasPrefix(line, "Disallow") {
            index := strings.IndexByte(line, byte('/'))
            if index >= 0 {
                line = line[index:]
            }
            index = strings.IndexFunc(line, func (r rune) bool {
                return r == '*' || r == '$' || r == '^' || r == '?'
            })
            if index > 0 {
                found.AddWord(root + line[:index])
            }
        } else if strings.HasPrefix(line, "Sitemap: ") {
            found.AddWord(line[len("Sitemap: "):])
        }
    }
    mut.Lock()
    res.AddAll(found)
    urls.AddElement(types.ReconedUrl{
        Endpoint: root + "/robots.txt",
        Status: resp.StatusCode,
        ResponseLength: resp.ContentLength,
    })
    mut.Unlock()
}

func FetchUrls(prog *types.ReconedProgram) {
    endpoint := utils.NewStringSet(nil)
    for _, url := range *prog.Urls {
        root := urlExtractor.FindString(url.Endpoint)
        endpoint.AddWord(root)
    }

    client := &http.Client{}
    found := utils.NewStringSet(nil)
    var mut sync.Mutex
    for _, url := range *endpoint {
        if len(url) == 0 {
            continue
        }
        prog.Throttler.RequestThread()
        go fetchUrlsWorker(url, client, found, prog.Throttler, &mut, prog.Urls)
    }
    prog.Throttler.Wait()
    found.AddAll(endpoint)
    for _, url := range found.ToArray() {
        if len(url) > 0 {
            prog.Urls.AddElement(types.ReconedUrl{
                Endpoint: url,
                Status: -1,
            })
        }
    }
}

