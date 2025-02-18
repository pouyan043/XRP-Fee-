package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const XRPL_API_URL = "https://responsive-blissful-snow.xrp-mainnet.quiknode.pro/8355909bcb1d4067f3184b5e37c90b05875fb30d"

type FeeResponse struct {
	Result struct {
		Drops struct {
			BaseFee       string `json:"base_fee"`
			MinimumFee    string `json:"minimum_fee"`
			MedianFee     string `json:"median_fee"`
			OpenLedgerFee string `json:"open_ledger_fee"`
		} `json:"drops"`
	} `json:"result"`
}

type rippleClient struct {
	rpcurl string
}

func (r *rippleClient) rpcRequest(command string, params []interface{}, rpcURL string) ([]byte, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  command,
		"params":  params,
		"id":      1,
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}

func (r *rippleClient) fee() (*FeeResponse, error) {
	params := []interface{}{}
	body, err := r.rpcRequest("fee", params, r.rpcurl)
	if err != nil {
		return nil, err
	}

	var responseJSON FeeResponse
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		return nil, err
	}

	return &responseJSON, nil
}

func convertDropsToXRP(drops string) string {
	dropValue, _ := strconv.Atoi(drops)
	return fmt.Sprintf("%.6f", float64(dropValue)/1000000.0)
}

func main() {
	client := &rippleClient{
		rpcurl: XRPL_API_URL,
	}

	feeResponse, err := client.fee()
	if err != nil {
		log.Fatalf("Error fetching fee: %v", err)
	}

	fmt.Println("Network Fee Details:")
	fmt.Printf("Base Fee: %s drops (~ %s XRP)\n", feeResponse.Result.Drops.BaseFee, convertDropsToXRP(feeResponse.Result.Drops.BaseFee))
	fmt.Printf("Minimum Fee: %s drops (~ %s XRP)\n", feeResponse.Result.Drops.MinimumFee, convertDropsToXRP(feeResponse.Result.Drops.MinimumFee))
	fmt.Printf("Median Fee: %s drops (~ %s XRP)\n", feeResponse.Result.Drops.MedianFee, convertDropsToXRP(feeResponse.Result.Drops.MedianFee))
	fmt.Printf("Open Ledger Fee: %s drops (~ %s XRP)\n", feeResponse.Result.Drops.OpenLedgerFee, convertDropsToXRP(feeResponse.Result.Drops.OpenLedgerFee))
}
