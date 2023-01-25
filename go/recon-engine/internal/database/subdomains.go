package database

/*
import (
    "database/sql"
    "errors"
    "fmt"
)

const (
    findSubdomainQuery = `SELECT id FROM subdomains WHERE prog_id=$1 AND subdomain=$2;`
    insertSubdomainQuery = `INSERT INTO subdomains (prog_id, subdomain) VALUES ($1, $2);`
)

func (data *DataManager) insertSubdomain(value subdomainDao, find *sql.Stmt, insert *sql.Stmt) (Diffs, error) {

    var dao subdomainDao
    err := find.QueryRow(value.ProgId, value.Subdomain).Scan(&dao.Id)

    if err != nil {
        _, err := insert.Exec(value.ProgId, value.Subdomain)
        if err != nil {
            return nil, errors.New("DataManager.insertSubdomains: could not insert data in subdomains")
        }
        return (&value).Differentiate(nil), nil
    }
    return nil, nil
}

func (data *DataManager) insertSubdomains(progId int64, subdomains []string) (DiffsList, error) {
    if data.db == nil {
        return nil, errors.New("DataManager has not been Init()")
    } else if err := data.db.Ping(); err != nil {
        return nil, errors.New(fmt.Sprintf("DataManager.insertSubdomains: Could not connect to db")) 
    }
    insert, err := data.db.Prepare(insertSubdomainQuery)

    if err != nil {
        return nil, errors.New("DataManager.insertSubdomains: could not create insert prepared statement")
    }
    defer insert.Close()
    find, err := data.db.Prepare(findSubdomainQuery)
    if err != nil {
        return nil, errors.New("DataManager.insertSubdomains: could not create select prepared statement")
    }

    defer find.Close()
    res := make(DiffsList, 0)

    for _, sub := range subdomains {
        dao := subdomainDao{
            ProgId: progId,
            Subdomain: sub,
        }
        d, err := data.insertSubdomain(dao, find, insert)
        if err != nil {
            continue
        }
        if len(d) > 0 {
            res[sub] = d
        }
    }


    return res, nil
}
*/
