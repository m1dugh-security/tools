package types

import (
    "github.com/m1dugh/recon-engine/internal/types"
    "github.com/m1dugh/recon-engine/internal/utils"
    programs "github.com/m1dugh/recon-engine/pkg/programs/types"
    "github.com/m1dugh/nmapgo/pkg/nmapgo"
    "regexp"
)

type ReconedUrl struct {
    Endpoint string
    Status int
    ResponseLength int64
}

func (u ReconedUrl) Compare(other interface{}) int {
    u2, ok := other.(ReconedUrl)
    if !ok {
        return -1
    } else {
        return utils.CompareStrings(u.Endpoint, u2.Endpoint)
    }
}

type ReconedProgram struct {
    Program     *programs.Program
    Subdomains  *types.StringSet
    Urls        *types.ComparableSet[ReconedUrl]
    Hosts       []*nmapgo.Host
    Throttler   *ThreadThrottler
}

func (prog *ReconedProgram) GetUrls() []string {
    res := make([]string, prog.Urls.Length())
    for i, v := range *prog.Urls {
        res[i] = v.Endpoint
    }

    return res
}

var _webRegex *regexp.Regexp = regexp.MustCompile(`web|api|http`)
func (prog *ReconedProgram) ExtractScopeInfo() {
    scope := prog.Program.GetScope(_webRegex)
    urls, subdomains := scope.ExtractInfo()
    for _, u := range urls.UnderlyingArray() {
        prog.Urls.AddElement(ReconedUrl{
            Endpoint: u,
            Status: -1,
            ResponseLength: -1,
        })
    }
    prog.Subdomains.AddAll(subdomains)
}

func NewReconedProgram(prog *programs.Program, throttler *ThreadThrottler) *ReconedProgram {
    res := &ReconedProgram{}
    res.Program = prog
    if throttler == nil {
        throttler = NewThreadThrottler(10)
    }
    res.Subdomains = types.NewStringSet(nil)
    res.Urls = types.NewComparableSet[ReconedUrl](nil)
    res.Throttler = throttler
    return res
}
