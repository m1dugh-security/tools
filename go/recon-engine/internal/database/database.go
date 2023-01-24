package database

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "errors"
    "github.com/m1dugh/recon-engine/pkg/types"
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
    db *sql.DB
    Config *Config
}

func New(config *Config) *DataManager {
    if config == nil {
        config = DefaultConfig()
    }

    return &DataManager{
        Config: config,
    }
}

func (data *DataManager) Init() error {
    psqlinfo := fmt.Sprintf("host=%s user=%s password=%s " +
    "dbname=%s sslmode=disable", data.Config.Host, data.Config.User, data.Config.Password, data.Config.DBName)

    db, err := sql.Open(data.Config.Driver, psqlinfo)
    if err != nil {
        return errors.New(
                fmt.Sprintf("DataManager.Init: Could not open %s at %s:%d",
                    data.Config.User, data.Config.Host, data.Config.Port))
    }
    data.db = db

    err = db.Ping()
    if err != nil {
        return errors.New(fmt.Sprintf("DataManager.Init: Error while reaching db"))
    }

    return nil
}


func (data *DataManager) Close() {
    data.db.Close()
    data.db = nil
}

func (data *DataManager) SaveProgram(prog *types.ReconedProgram) (ProgramDiff, error) {
    res := ProgramDiff{}
    diff, progId, err := data.insertProgram(prog.Program)
    if err != nil {
        return res, errors.New("DataManager.SaveProgram: could not save program")
    }
    res.ProgramDiff = diff
    res.ProgId = int(progId)

    res.ServiceDiff, err = data.insertHosts(progId, prog.Hosts)
    if err != nil {
        return res, errors.New("DataManager.SaveProgram: could not save services in program")
    }
    res.SubdomainDiff, err = data.insertSubdomains(progId, prog.Subdomains.ToArray())
    if err != nil {
        return res, errors.New("DataManager.SaveProgram: could not save subdomains in program")
    }

    res.UrlDiff, err = data.insertUrls(progId, prog.Urls.ToArray())
    if err != nil {
        return res, errors.New("DataManager.SaveProgram: could not save urls in program")
    }
    return res, nil
}


