package types

import (
    "github.com/m1dugh-security/tools/go/utils/pkg/utils"
    "github.com/m1dugh/program-browser/pkg/types"
    "github.com/m1dugh/nmapgo/pkg/nmapgo"
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
    Program     *types.Program
    Subdomains  *utils.StringSet
    Urls        *utils.ComparableSet[ReconedUrl]
    Hosts       []*nmapgo.Host
    Throttler   *utils.ThreadThrottler
}

func (prog *ReconedProgram) GetUrls() []string {
    res := make([]string, prog.Urls.Length())
    for i, v := range *prog.Urls {
        res[i] = v.Endpoint
    }

    return res
}

func (prog *ReconedProgram) ExtractScopeInfo() {
    scope := prog.Program.GetScope(types.Website, types.API)
    urls, subdomains := scope.ExtractInfo()
    for _, u := range urls.UnderlyingArray() {
        prog.Urls.AddElement(ReconedUrl{
            Endpoint: u,
            Status: -1,
            ResponseLength: -1,
        })
    }

    for _, subdomain := range subdomains.UnderlyingArray() {
        prog.Subdomains.AddWord(subdomain)
    }
}

func NewReconedProgram(prog *types.Program, throttler *utils.ThreadThrottler) *ReconedProgram {
    res := &ReconedProgram{}
    res.Program = prog
    if throttler == nil {
        throttler = utils.NewThreadThrottler(10)
    }
    res.Subdomains = utils.NewStringSet(nil)
    res.Urls = utils.NewComparableSet[ReconedUrl](nil)
    res.Throttler = throttler
    return res
}
