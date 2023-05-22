package main

import (
	"encoding/json"
        "flag"
        "fmt"
        "io/ioutil"
        "net/http"
        "net/url"
        "strconv"
        "time"
)

type TransactionResp struct {
        Data    []Data `json:"data"`
        Success bool   `json:"success"`
        Meta    Meta   `json:"meta"`
        RawResp []byte `json:"-"`
}
type TokenInfo struct {
        Symbol   string `json:"symbol"`
        Address  string `json:"address"`
        Decimals int    `json:"decimals"`
        Name     string `json:"name"`
}
type Data struct {
        TransactionID  string    `json:"transaction_id"`
        TokenInfo      TokenInfo `json:"token_info"`
        BlockTimestamp int64     `json:"block_timestamp"`
        From           string    `json:"from"`
        To             string    `json:"to"`
        Type           string    `json:"type"`
        Value          string    `json:"value"`
}
type Links struct {
        Next string `json:"next"`
}
type Meta struct {
        At          int64  `json:"at"`
        Fingerprint string `json:"fingerprint"`
        Links       Links  `json:"links"`
        PageSize    int    `json:"page_size"`
}

var from string

func main() {
        flag.StringVar(&from, "a", "", "")
        flag.Parse()
        ti := time.Tick(3 * time.Second)
	for ;true;<-ti {
                to := GetTrc20Trx(from, time.Now().UnixMilli()-90000, time.Now().UnixMilli(), true)
                if to != "" {
                        GetTrc20Trx(to, time.Now().UnixMilli()-90000, time.Now().UnixMilli(), false)
                }
		fmt.Println("-----------------------------------------------------")
        }
}
func GetTrc20Trx(account string, minTs, maxTs int64, tf bool) string {
        baseUrl, err := url.Parse(fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions/trc20", account))
        if err != nil {
                return ""
        }
        vals := url.Values{}
        //vals.Add("only_confirmed", "true")
        if tf {
                vals.Add("only_from", "true")
        } else {
                vals.Add("only_to", "true")
        }
        vals.Add("contract_address", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
        if minTs > 0 {
                vals.Add("min_timestamp", strconv.Itoa(int(minTs)))
        }
        if maxTs > 0 {
                vals.Add("max_timestamp", strconv.Itoa(int(maxTs)))
        }
        //vals.Add("order_by", "block_timestamp,asc")
        //vals.Add("limit", "200")
        baseUrl.RawQuery = vals.Encode()
        req, _ := http.NewRequest("GET", baseUrl.String(), nil)
        req.Header.Add("accept", "application/json")
        req.Header.Add("TRON_PRO_API_KEY", "b731a248-0c4d-4add-b396-62a1b639b5fe")
        fmt.Println(baseUrl.String())
        res, err := http.DefaultClient.Do(req)
        //fmt.Println(req.RequestURI)
        if err != nil {
                return ""
        }
        defer res.Body.Close()
        body, _ := ioutil.ReadAll(res.Body)
        fmt.Println(string(body))
        transactions := &TransactionResp{}
        transactions.RawResp = body
        json.Unmarshal(body, transactions)
        //fmt.Println(string(body))
	fmt.Println(len(transactions.Data))
        if len(transactions.Data) > 0 {
                return transactions.Data[0].To
        }
        return ""
}

