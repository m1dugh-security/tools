package database

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/m1dugh/nmapgo/pkg/nmapgo"
)

const (
    findServiceQuery = `SELECT * FROM services WHERE prog_id=$1 AND subdomain=$2 AND port=$3;`
    updateServiceQuery = `UPDATE services SET ip_addr=$2,protocol=$3,name=$4,product=$5,version=$6,additionals=$7 WHERE id=$1;`
    insertServiceQuery = `INSERT INTO services (prog_id,subdomain,ip_addr,port,protocol,name,product,version,additionals) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);`
)

func (data *DataManager) insertService(value serviceDao, find *sql.Stmt, insert *sql.Stmt, update *sql.Stmt) (Diffs, error) {
    var dao serviceDao
    err := find.QueryRow(value.ProgId, value.Subdomain, value.Port).Scan(
        &dao.Id,
        &dao.ProgId,
        &dao.Subdomain,
        &dao.Addr,
        &dao.Port,
        &dao.Protocol,
        &dao.Name,
        &dao.Product,
        &dao.Version,
        &dao.Additionals,
    )

    if err != nil {
        _, err := insert.Exec(
            value.ProgId,
            value.Subdomain,
            value.Addr,
            value.Port,
            value.Protocol,
            value.Name,
            value.Product,
            value.Version,
            value.Additionals,
        )

        if err != nil {
            return nil, errors.New("DataManager.insertService: could not insert data in services")
        }

        return (&value).Differentiate(nil), nil
    }
    d := (&value).Differentiate(&dao)
    if len(d) > 0 {
        _, err := update.Exec(
            dao.Id,
            value.Addr,
            value.Protocol,
            value.Name,
            value.Product,
            value.Version,
            value.Additionals,
        )

        if err != nil {
            return nil, errors.New(fmt.Sprintf("DataManager.insertService: could not update data at index %d", dao.Id))
        }

        return d, nil
    }
    return nil, nil
}

func portToServiceDao(progId int64, subdomain string, addr string, val nmapgo.Port) serviceDao {
    return serviceDao{
        ProgId: progId,
        Subdomain: subdomain,
        Addr: addr,
        Port: val.Port,
        Protocol: val.Protocol,
        Name: val.Service.Name,
        Product: val.Service.Product,
        Version: val.Service.Version,
        Additionals: val.Service.Additionals,
    }
}

func (data *DataManager) insertServices(progId int64, host *nmapgo.Host, res *DiffsList) error {
    
    if data.db == nil {
        return errors.New("DataManager has not been Init()")
    } else if err := data.db.Ping(); err != nil {
        return errors.New(fmt.Sprintf("DataManager.insertServices: Could not connect to db")) 
    }
    insert, err := data.db.Prepare(insertServiceQuery)
    if err != nil {
        return errors.New("DataManager.insertServices: could not create insert prepared statement")
    }
    defer insert.Close()
    update, err := data.db.Prepare(updateServiceQuery)
    if err != nil {
        return errors.New("DataManager.insertServices: could not create update prepared statement")
    }
    defer update.Close()
    find, err := data.db.Prepare(findServiceQuery)
    if err != nil {
        return errors.New("DataManager.insertServices: could not create select prepared statement")
    }
    defer find.Close()

    for _, p := range host.Ports {
        dao := portToServiceDao(progId, host.Hostname, host.Address, p)
        d, err := data.insertService(dao, find, insert, update)
        if err != nil {
            continue
        }
        if len(d) > 0 {
            (*res)[fmt.Sprintf("%s:%d", host.Hostname, p.Port)] = d
        }
    }
    return nil
}

func (data *DataManager) insertHosts(progId int64, hosts []*nmapgo.Host) (DiffsList, error) {
    var res DiffsList = make(DiffsList)
    for _, h := range hosts {
        err := data.insertServices(progId, h, &res)
        if err != nil {
            return res, err
        }
    }
    return res, nil
}
