package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
)

var (
    apiUrl      = "https://min-api.cryptocompare.com/data/pricemultifull"
    fsymsValues = []string{"BTC", "XRP", "ETH", "BCH", "EOS", "LTC", "XMR", "DASH"}
    tsymsValues = []string{"USD", "EUR", "GBP", "JPY", "RUR"}
)

func getDataFromCrypto(fsyms string, tsyms string) string {
    url := ""
    if fsyms != "" && tsyms != "" {
        url = makeUrl([]string{fsyms}, []string{tsyms})
    } else {
        url = makeUrl(fsymsValues, tsymsValues)
    }
    
    fmt.Printf("url: %#v\n\n", url)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("get by url error: ", err.Error())
        return ""
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        fmt.Println("Status code = " + strconv.Itoa(resp.StatusCode))
        return ""
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("read body error: ", err.Error())
        return ""
    }
    
    fmt.Printf("bosy: %#v\n\n", string(body))

    respData := filterCryptoData(string(body), fsyms, tsyms)

    return respData
}

func makeUrl(f []string, t []string) string {
    return apiUrl + "?fsyms=" + strings.Join(f, ",") + "&tsyms=" + strings.Join(t, ",")
}
