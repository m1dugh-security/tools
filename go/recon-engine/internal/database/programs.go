package database

import (
    "database/sql"
    "errors"
    "time"
    ptypes "github.com/m1dugh/recon-engine/pkg/programs/types"
    "fmt"
)

const (
    findProgramQuery = `SELECT id,code,name,platform,platform_url,status,safe_harbor,managed,category from programs WHERE platform=$1 AND name=$2;`
    countProgramQuery = `SELECT COUNT(*) FROM programs WHERE platform=$1 AND name=$2;`
    updateProgramQuery = `UPDATE programs SET status=$1,safe_harbor=$2,managed=$3,category=$4,recon_date=to_timestamp($5) WHERE id=$6;`
    updateProgramTimestamp = `UPDATE programs SET recon_date=to_timestamp($1) WHERE id=$2;`
    insertProgramQuery = `INSERT INTO programs (code,name,platform,platform_url,status,safe_harbor,managed,category,recon_date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,to_timestamp($9));`
    deleteProgramQuery = `DELETE from programs WHERE id=$1;`
)

// returns the diff, the progId, and an error if applicable
func (data *DataManager) insertProgramDao(value programDao,
    find *sql.Stmt,
    del *sql.Stmt,
    insert *sql.Stmt,
    update *sql.Stmt,
    updateTimestamp *sql.Stmt,
) (Diffs, int64, error) {

    timestamp := time.Now().Unix()
    var dao programDao
    err := find.QueryRow(value.Platform, value.Name).Scan(
        &dao.Id,
        &dao.Code,
        &dao.Name,
        &dao.Platform,
        &dao.PlatformUrl,
        &dao.Status,
        &dao.SafeHarbor,
        &dao.Managed,
        &dao.Category,
    )


    if err != nil {
        _, err := insert.Exec(
            value.Code,
            value.Name,
            value.Platform,
            value.PlatformUrl,
            value.Status,
            value.SafeHarbor,
            value.Managed,
            value.Category,
            timestamp,
        )
        if err != nil {
            return nil, -1,errors.New("DataManager.insertProgram: Error while inserting program")
        }
        var id int64
        st, err := data.db.Prepare(`SELECT id FROM programs WHERE code=$1;`)
        if err != nil {
            return nil, -1, errors.New("DataManager.insertProgramDao: error when retrieving index for inserted element")
        }
        defer st.Close()
        err = st.QueryRow(value.Code).Scan(&id)
        if err != nil {
            return nil, -1, errors.New("DataManager.insertProgramDao: error when readind index for inserted element")
        }
        return (&value).Differentiate(nil), id, nil
    }
    d := (&value).Differentiate(&dao)
    if len(d) > 0 {
        _, err := update.Exec(
            value.Status,
            value.SafeHarbor,
            value.Managed,
            value.Category,
            timestamp,
            dao.Id,
        )
        if err != nil {
            return nil, dao.Id, errors.New("DataManager.insertProgram: Error while updating program")
        }
    } else {
        _, err := updateTimestamp.Exec(timestamp, dao.Id)
        if err != nil {
            return nil, dao.Id, errors.New("DataManager.insertProgram: Error while updating timestamp")
        }
    }

    return d, dao.Id, nil
}


func (data *DataManager) insertProgram(p *ptypes.Program) (Diffs, int64, error) {
    if data.db == nil {
        return nil, -1, errors.New("DataManager has not been Init()")
    } else if err := data.db.Ping(); err != nil {
        return nil, -1, errors.New(fmt.Sprintf("DataManager.insertPrograms: Could not connect to db")) 
    }
    insert, err := data.db.Prepare(insertProgramQuery)
    if err != nil {
        return nil, -1, errors.New("DataManager.insertPrograms: could not create insert prepared statement")
    }
    del, err := data.db.Prepare(deleteProgramQuery)
    if err != nil {
        return nil, -1, errors.New("DataManager.insertPrograms: could not create delete prepared statement")
    }
    update, err := data.db.Prepare(updateProgramQuery)
    if err != nil {
        return nil, -1, errors.New("DataManager.insertPrograms: could not create update prepared statement")
    }
    updateTimestamp, err := data.db.Prepare(updateProgramTimestamp)
    if err != nil {
        return nil, -1, errors.New("DataManager.insertPrograms: could not create update prepared statement")
    }
    find, err := data.db.Prepare(findProgramQuery)
    if err != nil {
        return nil, -1, errors.New("DataManager.insertPrograms: could not create select prepared statement")
    }

    dao := programDao{
        Code: p.Code(),
        Name: p.Name,
        Platform: p.Platform,
        PlatformUrl: p.PlatformUrl,
        Status: p.Status,
        SafeHarbor: p.SafeHarborStatus,
        Managed: p.Managed,
    }

    defer insert.Close()
    defer del.Close()
    defer update.Close()
    defer find.Close()
    return data.insertProgramDao(dao, find, del, insert, update, updateTimestamp)
}

func (data *DataManager) ProgramExists(p *ptypes.Program) (bool, error) {
    st, err := data.db.Prepare(countProgramQuery)
    if err != nil {
        return false, errors.New(fmt.Sprintf("Error on creating prepared statement: %s", err))
    }
    defer st.Close()

    var count int
    err = st.QueryRow(p.Platform, p.Name).Scan(&count)
    if err != nil {
        return false, errors.New(fmt.Sprintf("Error on querying programs: %s", err))
    }

    return count > 0, nil
}
