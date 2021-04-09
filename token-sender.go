package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ququzone/binance-scripts/utils"
	"github.com/ququzone/go-common/env"
)

func main() {
	contentBytes, err := ioutil.ReadFile("./token-sender-list.csv")
	if err != nil {
		log.Fatalf("read csv file error: %v", err)
	}
	r := csv.NewReader(bytes.NewReader(contentBytes))
	rows, _ := r.ReadAll()

	signer := utils.NewHmacSigner(env.GetNonEmpty("SECRET_KEY"))

	client := &http.Client{Transport: &http.Transport{}}

	for i, row := range rows {
		req, err := http.NewRequest("POST", "https://api.binance.com/wapi/v3/withdraw.html", nil)
		if err != nil {
			log.Fatalf("new http request error: %v", err)
		}
		q := req.URL.Query()
		q.Add("asset", env.GetNonEmpty("ASSET"))
		q.Add("network", env.GetNonEmpty("NETWORK"))
		q.Add("address", row[0])
		q.Add("amount", row[1])
		q.Add("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
		q.Add("signature", signer.Sign([]byte(q.Encode())))
		req.URL.RawQuery = q.Encode()

		req.Header.Add("X-MBX-APIKEY", env.GetNonEmpty("API_KEY"))

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("submit %d withdraw request error: %v", i, err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("read %d withdraw response body error: %v", i, err)
		}
		if resp.StatusCode != 200 {
			log.Fatalf("submit withdraw request error: %v", string(body))
		}
		log.Printf("Request %d withdraw successful: %s\n", i, string(body))
	}
}
