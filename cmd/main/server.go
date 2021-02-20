package main

import (
    "io"
    "log"
    "net/http"
    "net/url"
    "sync"
)

func runServer(wg *sync.WaitGroup) {
    defer wg.Done()

    println("Start Crypto ...")
    http.HandleFunc("/price", CryptoServer)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func CryptoServer(writer http.ResponseWriter, request *http.Request) {
    if request.Method != "GET" {
        writer.WriteHeader(http.StatusMethodNotAllowed)
        io.WriteString(writer, "{\"result\":\"fail\",\"data\":\"GET request expected\"}")
        return
    }

    fsyms, tsyms := parseParamsFromUri(request.URL.Query())
    if fsyms == "" || tsyms == "" {
        io.WriteString(writer, "{\"result\":\"fail\",\"data\":\"Incorrect params, expected uri like fsyms=BTC&tsyms=USD\"}")
        return
    }

    if !checkSyms(fsyms, tsyms) {
        io.WriteString(writer, "{\"result\":\"fail\",\"data\":\"Incorrect params, fsyms must be like BTC ..., tsyms like USD ...\"}")
        return
    }

    data := getDataFromCrypto(fsyms, tsyms)
    if data == "" {
        cryptoData, err := readData(db, fsyms, tsyms)
        if err != nil {
            writer.WriteHeader(http.StatusBadRequest)
            io.WriteString(writer, "{\"result\":\"fail\",\"data\":\""+err.Error()+"\"}")
            return
        }
        data = cryptoData
    }

    io.WriteString(writer, data)
}

func parseParamsFromUri(values url.Values) (fsyms string, tsyms string) {
    fsymsAny, ok := values["fsyms"]
    if !ok || len(fsymsAny) < 1 {
        return "", ""
    }
    fsyms = fsymsAny[0]

    tsymsAny, ok := values["tsyms"]
    if !ok || len(tsymsAny) < 1 {
        return "", ""
    }
    tsyms = tsymsAny[0]

    return fsyms, tsyms
}
