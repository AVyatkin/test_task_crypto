package main

import (
    "database/sql"
    "sync"
)

var (
    db *sql.DB
)

func main() {
    db = initDb()
    createTableIfNotExist(db)

    wg := sync.WaitGroup{}
    wg.Add(1)
    go runScheduler(&wg, db)

    wg.Add(1)
    go runServer(&wg)

    wg.Wait()
}
