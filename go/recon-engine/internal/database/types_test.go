package database

import (
    "testing"
    "reflect"
)

func TestCalculateDiffSubdomains(t *testing.T) {
    v1 := subdomainDao{
        Subdomain: "test1",
    }

    v2 := subdomainDao{
        Subdomain: "test2",
    }

    res := calculateDiff(&v1, &v2, "Subdomain")
    expected := make(Diffs)
    expected["Subdomain"] = diff{
        "test1",
        "test2",
    }

    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Error when comparing: \n%s\nis not matching:\n%s", res, expected)
    }
}

func TestCalculateDiffNil(t *testing.T) {

    v2 := subdomainDao{
        Subdomain: "test2",
    }

    res := calculateDiff(nil, &v2, "Subdomain")
    expected := make(Diffs)
    expected["Subdomain"] = diff{
        nil,
        "test2",
    }

    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Error when comparing: \n%s\nis not matching:\n%s", res, expected)
    }
}

func TestCalculateDiffUrl(t *testing.T) {
    v1 := urlDao{
        Endpoint: "http://example.com",
        Status: 200,
        ResponseLength: 42,
    }

    v2 := urlDao{
        Endpoint: "http://example.com",
        Status: 403,
        ResponseLength: 42,
    }

    res := calculateDiff(&v1, &v2, "Status")
    expected := make(Diffs)
    expected["Status"] = diff{
        200,
        403,
    }

    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Error when comparing: \n%s\nis not matching:\n%s", res, expected)
    }
}
