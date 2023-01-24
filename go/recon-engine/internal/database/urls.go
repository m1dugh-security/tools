package database

import (
    "database/sql"
    "github.com/m1dugh/recon-engine/pkg/types"
    "errors"
    "fmt"
)

const (
    findURLQuery = `SELECT id,endpoint,status,response_length from urls WHERE prog_id=$1 AND endpoint=$2;`
    updateURLQuery = `UPDATE urls SET status=$2,response_length=$3 WHERE id=$1;`
    insertURLQuery = `INSERT INTO urls (prog_id, endpoint, status, response_length) VALUES ($1, $2, $3, $4);`
    deleteURLQuery = `DELETE from urls WHERE id=$1;`
)


func (data *DataManager) insertURL(value urlDao, find *sql.Stmt, del *sql.Stmt, insert *sql.Stmt, update *sql.Stmt) (Diffs, error) {

    var dao urlDao
    err := find.QueryRow(value.ProgId, value.Endpoint).Scan(&dao.Id, &dao.Endpoint, &dao.Status, &dao.ResponseLength)

    if err != nil {
        _, err := insert.Exec(value.ProgId, value.Endpoint, value.Status, value.ResponseLength)
        if err != nil {
            return nil, errors.New("DataManager.insertURL: could not insert data in urls")
        }
        return (&value).Differentiate(nil), nil
    }

    d := (&value).Differentiate(&dao)
    if len(d) > 0 {
        _, err := update.Exec(dao.Id, value.Status, value.ResponseLength)
        if err != nil {
            return nil, errors.New(fmt.Sprintf("DataManager.insertURL: could not update item %d", dao.Id))
        }
        return d, nil
    }
    return nil, nil
}

func (data *DataManager) insertUrls(progId int64, urls []types.ReconedUrl) (DiffsList, error) {
    if data.db == nil {
        return nil, errors.New("DataManager has not been Init()")
    } else if err := data.db.Ping(); err != nil {
        return nil, errors.New(fmt.Sprintf("DataManager.insertURLs: Could not connect to db")) 
    }
    insert, err := data.db.Prepare(insertURLQuery)
    if err != nil {
        return nil, errors.New("DataManager.insertURLs: could not create insert prepared statement")
    }
    del, err := data.db.Prepare(deleteURLQuery)
    if err != nil {
        return nil, errors.New("DataManager.insertURLs: could not create delete prepared statement")
    }
    update, err := data.db.Prepare(updateURLQuery)
    if err != nil {
        return nil, errors.New("DataManager.insertURLs: could not create update prepared statement")
    }
    find, err := data.db.Prepare(findURLQuery)
    if err != nil {
        return nil, errors.New("DataManager.insertURLs: could not create select prepared statement")
    }

    res := make(DiffsList, 0)

    for _, u := range urls {
        dao := urlDao{
            ProgId: progId,
            Endpoint: u.Endpoint,
            Status: u.Status,
            ResponseLength: int(u.ResponseLength),
        }
        d, err := data.insertURL(dao, find, del, insert, update)
        if err != nil {
            continue
        }
        if len(d) > 0 {
            res[u.Endpoint] = d
        }
    }

    insert.Close()
    del.Close()
    update.Close()
    find.Close()

    return res, nil
}
