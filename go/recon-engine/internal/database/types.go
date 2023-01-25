package database

import (
	"time"

	"github.com/m1dugh-security/tools/go/recon-engine/pkg/types"
	"github.com/uptrace/bun"
)

type urlDao struct {
    bun.BaseModel       `bun:"table:urls"`
    ID int64            `bun:",pk,autoincrement"`
    ProgramId int64     `bun:"program_id"`
    Endpoint string     `bun:"endpoint"`
    Status int          `bun:"status"`
    ResponseLength int  `bun:"response_length"`
}

type subdomainDao struct {
    ProgramID int64
    ID int64            `bun:",pk,autoincrement"`
    Subdomain string
}

type targetDao struct {
    bun.BaseModel       `bun:"table:targets"`
    ID int64            `bun:",pk,autoincrement"`
    ProgramId int64     `bun:"program_id"`
    Category string     `bun:"target"`
}

type serviceDao struct {
    bun.BaseModel       `bun:"table:services"`
    ID int64            `bun:",pk,autoincrement"`
    ProgramId int64     `bun:"program_id"`
    Subdomain string    `bun:"subdomain"`
    Addr string         `bun:"address"`
    Port int            `bun:"port"`
    Protocol string     `bun:"protocol"`
    Name string         `bun:"name"`
    Product string      `bun:"product"`
    Version string      `bun:"version"`
    Additionals string  `bun:"additinals"`
}

type programDao struct {
    bun.BaseModel       `bun:"table:programs"`
    ID int64            `bun:",pk,autoincrement"`
    Code string         `bun:"code,notnull,unique"`
    Name string         `bun:"name,notnull"`
    Platform string     `bun:"platform,notnull"`
    PlatformUrl string  `bun:"platform_url,notnull"`
    Status string       `bun:"status"`
    SafeHarbor string   `bun:"safe_harbor"`
    Managed bool        `bun:"managed"`
    Category string     `bun:"category"`
    ReconDate time.Time `bun:"recon_date,default:current_timestamp"`
    Urls []*urlDao      `bun:"urls,rel:has-many,join:id=program_id"`
    Services []*serviceDao  `bun:"services,rel:has-many,join:id=program_id"`     
    Targets []*targetDao    `bun:"targets,rel:has-many,join:id=program_id"`
    Subdomains []*subdomainDao  `bun:"rel:has-many,join:id=program_id"`
}

func toProgramDao(rc *types.ReconedProgram) *programDao {
    prog := rc.Program
    res := &programDao{
        Code: prog.Code(),
        Name: prog.Name,
        Platform: prog.Platform,
        PlatformUrl: prog.PlatformUrl,
        Status: prog.Status,
        Managed: prog.Managed,
        Category: prog.Category,
        Urls: make([]*urlDao, 0, rc.Urls.Length()),
        Subdomains: make([]*subdomainDao, 0, rc.Subdomains.Length()),
    }

    for _, url := range rc.Urls.UnderlyingArray() {
        res.Urls = append(res.Urls, &urlDao{
            Endpoint: url.Endpoint,
            Status: url.Status,
            ResponseLength: int(url.ResponseLength),
        })
    }

    for _, subdomain := range rc.Subdomains.UnderlyingArray() {
        res.Subdomains = append(res.Subdomains, &subdomainDao{
            Subdomain: subdomain,
        })
    }

    return res
}

// func fromProgDao(dao *programDao)

