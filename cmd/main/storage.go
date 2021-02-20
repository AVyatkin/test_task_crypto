package main

import (
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "time"
)

type response struct {
    Raw     mapItemRaw     `json:"RAW"`
    Display mapItemDisplay `json:"DISPLAY"`
}

type mapItemRaw map[string]map[string]itemRaw
type mapItemDisplay map[string]map[string]itemDisplay

type itemRaw struct {
    Change24hour    float64 `json:"CHANGE24HOUR"`
    Changepct24hour float64 `json:"CHANGEPCT24HOUR"`
    Open24hour      float64 `json:"OPEN24HOUR"`
    Volume24hour    float64 `json:"VOLUME24HOUR"`
    Volume24hourto  float64 `json:"VOLUME24HOURTO"`
    Low24hour       float64 `json:"LOW24HOUR"`
    High24hour      float64 `json:"HIGH24HOUR"`
    Price           float64 `json:"PRICE"`
    Lastupdate      int64   `json:"LASTUPDATE"`
    Supply          float64 `json:"SUPPLY"`
    Mktcap          float64 `json:"MKTCAP"`
}

type itemDisplay struct {
    Change24hour    string `json:"CHANGE24HOUR"`
    Changepct24hour string `json:"CHANGEPCT24HOUR"`
    Open24hour      string `json:"OPEN24HOUR"`
    Volume24hour    string `json:"VOLUME24HOUR"`
    Volume24hourto  string `json:"VOLUME24HOURTO"`
    Low24hour       string `json:"LOW24HOUR"`
    High24hour      string `json:"HIGH24HOUR"`
    Price           string `json:"PRICE"`
    Lastupdate      string `json:"LASTUPDATE"`
    Supply          string `json:"SUPPLY"`
    Mktcap          string `json:"MKTCAP"`
}

var mysqlSource = "mysql_user:mysql_password@tcp(db:3306)/mysql_db"

func initDb() *sql.DB {
    db, err := sql.Open("mysql", mysqlSource)
    if err != nil {
        panic(err)
    }

    db.SetConnMaxLifetime(time.Minute * 3)
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10)

    return db
}

func createTableIfNotExist(db *sql.DB) {
    _, err := db.Query(`create table if not exists crypto
            (
                timestamp bigint  not null,
                symbol    text    not null
            )`)
    if err != nil {
        panic(err.Error())
    }
}

func writeData(db *sql.DB, data string) {
    insertString := `insert crypto (timestamp, symbol) values (?, ?)`
    insert, err := db.Query(insertString, strconv.FormatInt(time.Now().Unix(), 10), data)
    if err != nil {
        panic(err.Error())
    }
    insert.Close()
}

func readData(db *sql.DB, fsyms string, tsyms string) (responseData string, err error) {

    rows, err := db.Query("select timestamp, symbol from crypto order by timestamp desc limit 1")
    defer func() {
        if rows != nil {
            rows.Close()
        }
    }()

    if err != nil {
        fmt.Printf("err: %#v\n\n", err.Error())
    }

    if rows == nil || !rows.Next() {
        return "", errors.New("no data in db")
    }

    ts := 0
    responseDataSource := ""
    err = rows.Scan(&ts, &responseDataSource)
    if err != nil {
        fmt.Println("error:" + err.Error())
    }

    fmt.Printf("ts: %#v\n\n", ts)

    respData := filterCryptoData(responseDataSource, fsyms, tsyms)

    return respData, nil
}

func filterCryptoData(data string, fsyms string, tsyms string) string {
    var respData response
    err := json.Unmarshal([]byte(data), &respData)
    if err != nil {
        fmt.Println("unmarshal error: ", err.Error())
        return ""
    }

    if !checkSyms(fsyms, tsyms) {
        return ""
    }

    dataFiltered := response{}
    if fsyms == "" || tsyms == "" {
        fmt.Println("1")
        dataFiltered = respData
    } else {
        dataFiltered = response{
            mapItemRaw{
                fsyms: map[string]itemRaw{
                    tsyms: respData.Raw[fsyms][tsyms],
                },
            },
            mapItemDisplay{
                fsyms: map[string]itemDisplay{
                    tsyms: respData.Display[fsyms][tsyms],
                },
            },
        }
    }

    dataFilteredBytes, _ := json.Marshal(dataFiltered)
    fmt.Printf("respData: %s\n\n", string(dataFilteredBytes))

    return string(dataFilteredBytes)
}

func getSet(slice []string) map[string]struct{} {
    set := make(map[string]struct{})
    for _, s := range slice {
        set[s] = struct{}{}
    }
    return set
}

func checkSyms(fsyms string, tsyms string) bool {
    setF := getSet(fsymsValues)
    setT := getSet(tsymsValues)
    if _, ok := setF[fsyms]; !ok && fsyms != "" {
        return false
    }
    if _, ok := setT[tsyms]; !ok && tsyms != "" {
        return false
    }

    return true
}
