package main

import (
    "database/sql"
    "log"
    "sync"
    "time"
)

func runScheduler(wg *sync.WaitGroup, db *sql.DB) {
    defer wg.Done()

    ticker := time.NewTicker(time.Minute)
    stop := make(chan struct{})
    for {
        select {
        case <-ticker.C:
            storeDataFromCrypto(db)
        case <-stop:
            ticker.Stop()
            return
        }
    }
}

func storeDataFromCrypto(db *sql.DB) {
    log.Println("storeDataFromCrypto")

    body := getDataFromCrypto("", "")
    writeData(db, body)
}


