package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/m1dugh-security/tools/go/recon-engine/pkg/types"
	"github.com/m1dugh-security/tools/go/utils/pkg/utils"
	ptypes "github.com/m1dugh/program-browser/pkg/types"
)

func mockReconedProgram() *types.ReconedProgram {
    prog := &ptypes.Program{
        Name: "example.com",
        Platform: "none",
        Category: "website",
    }

    urls := utils.NewComparableSet[types.ReconedUrl](nil)

    urls.AddElement(types.ReconedUrl{
        Endpoint: "https://www.example.com",       
    })

    rc := &types.ReconedProgram{
        Program: prog,
        Subdomains: utils.NewStringSet([]string{"test.example.com", "example.com"}),
        Urls: urls,
    }

    return rc
}

func TestInsertDB(t *testing.T) {
    dao := toProgramDao(mockReconedProgram())

    manager := New(nil)
    manager.Init()
    defer manager.Close()

    fmt.Println(dao.Subdomains)

    ctx := context.Background()
    manager.db.NewInsert().
        Model(dao).
        On("CONFLICT (code) DO UPDATE").
        Exec(ctx)
}

func TestOpenDB(t *testing.T) {

    manager := New(nil)

    manager.Init()
    defer manager.Close()

    ctx := context.Background()

    progDao := new(programDao)

    err := manager.db.NewSelect().Model(progDao).Scan(ctx)
    if err != nil {
        t.Error(err)
    }

    fmt.Println(progDao.Subdomains)
}
