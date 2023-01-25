package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m1dugh-security/tools/go/recon-engine/pkg/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)


const (
    host string = "localhost"
    port int = 5432
    user = "postgres"
    password = "postgres"
    dbname = "recon-engine"
    driver = "postgres"
)

type Config struct {
    Host string
    Port int
    User string
    Password string
    DBName string
    Driver string
}

func DefaultConfig() *Config {
    return &Config{
        Host: host,
        Port: port,
        User: user,
        Password: password,
        DBName: dbname,
        Driver: driver,
    }
}

type DataManager struct {
    db *bun.DB
    Config *Config
    ctx context.Context
}

func New(config *Config) *DataManager {
    if config == nil {
        config = DefaultConfig()
    }

    return &DataManager{
        Config: config,
        ctx: context.Background(),
    }
}

func (data *DataManager) Init() {
    dsn := fmt.Sprintf(
        "%s://%s:%s@%s:%d/%s?sslmode=disable",
        data.Config.Driver,
        data.Config.User,
        data.Config.Password,
        data.Config.Host,
        data.Config.Port,
        data.Config.DBName,
    )

    sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

    db := bun.NewDB(sqldb, pgdialect.New())

    data.db = db

    db.NewCreateTable().Model(&programDao{}).IfNotExists().Exec(data.ctx)
    db.NewCreateTable().Model(&urlDao{}).IfNotExists().Exec(data.ctx)
    db.NewCreateTable().Model(&serviceDao{}).IfNotExists().Exec(data.ctx)
    db.NewCreateTable().Model(&subdomainDao{}).IfNotExists().Exec(data.ctx)
    db.NewCreateTable().Model(&targetDao{}).IfNotExists().Exec(data.ctx)
}

func (data *DataManager) SaveProgram(program *types.ReconedProgram) {
    data.db.NewInsert()
}


func (data *DataManager) Close() {
    data.db.Close()
    data.db = nil
}
