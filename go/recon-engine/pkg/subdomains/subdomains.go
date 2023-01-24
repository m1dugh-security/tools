package subdomains

import (
    "bytes"
    "io"
    "context"
    "github.com/projectdiscovery/subfinder/v2/pkg/resolve"
    "github.com/projectdiscovery/subfinder/v2/pkg/runner"
    "github.com/projectdiscovery/subfinder/v2/pkg/passive"
    btypes "github.com/m1dugh/recon-engine/pkg/programs/types"
    . "github.com/m1dugh/recon-engine/pkg/types"
    types "github.com/m1dugh/recon-engine/internal/types"
    "strings"
    "regexp"
)

func GetSubdomains(domains []string) (*types.StringSet, error) {

    runnerInstance, err := runner.NewRunner(&runner.Options{
        Threads:            10, // Thread controls the number of threads to use for active enumerations
        Timeout:            30, // Timeout is the seconds to wait for sources to respond
        MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
        Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
        Sources:            passive.DefaultSources, // Use the default list of passive sources
        AllSources:         passive.DefaultAllSources, // Use the default list of all passive sources
        Recursive:          passive.DefaultRecursiveSources,    // Use the default list of recursive sources
        Providers:          &runner.Providers{}, // Use empty api keys for all providers
    })

    buf := bytes.Buffer{}
    reader := strings.NewReader(strings.Join(domains, "\n"))
    err = runnerInstance.EnumerateMultipleDomains(context.Background(), reader, []io.Writer{&buf})
    if err != nil {
        return nil, err
    }

    data, err := io.ReadAll(&buf)
    if err != nil {
        return nil, err
    }
    arr := strings.Split(string(data), "\n")
    // removes last new line with empty value
    res := types.NewStringSet(nil)
    for _, sub := range arr {
        if len(sub) > 0 {
            res.AddWord(sub)
        }
    }
    return res, nil
}


func FetchSubdomains(prog *ReconedProgram) error {
    re := regexp.MustCompile(`web`)
    domains := getDomains(prog.Program.GetScope(re))
    subdomains, err := GetSubdomains(domains.ToArray()) 
    if err != nil {
        prog.Subdomains.AddAll(domains)
        return err
    }
    prog.Subdomains.AddAll(subdomains)
    prog.Subdomains.AddAll(domains)
    return nil
}

var _domainextractor = regexp.MustCompile(`^(\*\.)?([A-Za-z0-9\-]+\.)+[a-z]{2,8}$`)
func getDomains(scope *btypes.Scope) *types.StringSet {
    res := types.NewStringSet(nil)
    for _, v := range scope.Include {
        s := _domainextractor.FindString(v)
        if strings.HasPrefix(s, "*.") {
            s = s[2:]
        }
        if len(s) > 0 {
            res.AddWord(s)
        }
    }

    return res
}

