package database

import (
    "testing"
)

func TestOpenDB(t *testing.T) {
    manager := new(DataManager)

    err := manager.Init()
    if err != nil {
        t.Errorf("Error: %s", err)
    }
    defer manager.Close()
}
