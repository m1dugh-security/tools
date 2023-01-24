package database

import (
    "time"
    "reflect"
    "fmt"
    "encoding/json"
)

type urlDao struct {
    Id int64
    ProgId int64
    Endpoint string
    Status int
    ResponseLength int
}

type subdomainDao struct {
    Id int64
    ProgId int64
    Subdomain string
}

type targetDao struct {
    Id int64
    ProgId int64
    Category string
}

type serviceDao struct {
    Id int64
    ProgId int64
    Subdomain string
    Addr string
    Port int
    Protocol string
    Name string
    Product string
    Version string
    Additionals string
}

type programDao struct {
    Id int64
    Code string
    Name string
    Platform string
    PlatformUrl string
    Status string
    SafeHarbor string
    Managed bool
    Category string
    ReconDate time.Time
}

func (value *serviceDao) Differentiate(old *serviceDao) Diffs {
    keys := []string{"Protocol", "Name", "Product", "Version", "Additionals"}
    return calculateDiff(old, value, keys...)
}
func (value *programDao) Differentiate(old *programDao) Diffs {
    keys := []string{"Status", "SafeHarbor", "Managed", "Category"}
    return calculateDiff(old, value, keys...)
}
func (value *subdomainDao) Differentiate(old *subdomainDao) Diffs {
    keys := []string{"Subdomain"}
    return calculateDiff(old, value, keys...)
}

func (value *urlDao) Differentiate(old *urlDao) Diffs {
    keys := []string{"Endpoint", "Status", "Length"}
    return calculateDiff(old, value, keys...)
}
func (value *targetDao) Differentiate(old *targetDao) Diffs {
    return calculateDiff(old, value, "Category")
}

func calculateDiff(old interface{}, value interface{}, fields ...string) Diffs  {
    res := make(Diffs)

    var oldType, newType reflect.Value
    var oldNumField int = -1
    var newNumField int = -1
    if old != nil {
        oldType = reflect.ValueOf(old).Elem()
        if oldType.Kind() == reflect.Struct {
            oldNumField = oldType.NumField()
        }
    }
    if value != nil {
        newType = reflect.ValueOf(value).Elem()
        if newType.Kind() == reflect.Struct {
            newNumField = newType.NumField()
        }
    }
    for _, field := range fields {

        var oldIndex, newIndex int
        for oldIndex = 0; oldIndex < oldNumField;oldIndex++ {
            f := oldType.Type().Field(oldIndex)
            if f.Name == field {
                break
            }
        }
        for newIndex = 0; newIndex < newNumField;newIndex++ {
            f := newType.Type().Field(newIndex)
            if f.Name == field {
                break
            }
        }
        if oldIndex < oldNumField {
            oldValue := oldType.Field(oldIndex).Interface()
            if newIndex < newNumField {
                newValue := newType.Field(newIndex).Interface()
                if newValue != oldValue {
                    res[field] = diff{
                        Old: oldValue,
                        New: newValue,
                    }
                }
            } else {
                res[field] = diff{
                    Old: oldValue,
                    New: nil,
                }
            }
        } else if newIndex < newNumField { 

            newValue := newType.Field(newIndex).Interface()
            res[field] = diff{
                Old: nil,
                New: newValue,
            }
        }
    }
    return res
}

type Diffs map[string]diff

func (d Diffs) String() string {
    var res string
    for k, v := range d {
        res += fmt.Sprintf("%s: '%s' -> '%s'\n", k, v.Old, v.New)
    }

    return res
}

type Differentiable interface {
    Differentiate(old *interface{}) Diffs
}

type diff struct {
    Old interface{}
    New interface{}
}

type DiffsList map[string]Diffs

type ProgramDiff struct {
    ProgId int                  `json:"prog_id"`
    ProgramDiff Diffs           `json:"program_diff"`
    TargetDiff  DiffsList       `json:"target_diff"`
    SubdomainDiff   DiffsList   `json:"subdomain_diff"`
    UrlDiff     DiffsList       `json:"url_diff"`
    ServiceDiff DiffsList       `json:"service_diff"`
}

func (d ProgramDiff) String() string {
    content, err := json.Marshal(d)
    if err != nil {
        return ""
    }

    return string(content)
}

func (d ProgramDiff) IsEmpty() bool {
    return len(d.ProgramDiff) == 0 && len(d.TargetDiff) == 0 && len(d.SubdomainDiff) == 0 && len(d.UrlDiff) == 0 && len(d.ServiceDiff) == 0
}

