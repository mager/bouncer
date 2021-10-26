package premint

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

// https://www.premint.xyz/api/c7408fc0de8be594bd28d0d53eae2e90b3791d96148804caaa95a66bb4ecf897

type Data struct {
	WalletAddress string `json:"wallet_address"`
}

type Resp struct {
	Data []Data `json:"data"`
}

type PremintClient struct {
	httpClient *http.Client
}

// ProvidePremint provides an HTTP client
func ProvidePremint() PremintClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	return PremintClient{
		httpClient: &http.Client{
			Transport: tr,
		},
	}
}

var Options = ProvidePremint

// GetCollectionsForAddress returns the collections for an address
func (o *PremintClient) GetWalletAddresses(address string) ([]string, error) {
	u, err := url.Parse("https://www.premint.xyz/api/c7408fc0de8be594bd28d0d53eae2e90b3791d96148804caaa95a66bb4ecf897")
	if err != nil {
		log.Fatal(err)
		return []string{}, nil
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
		return []string{}, nil
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return []string{}, nil
	}
	defer resp.Body.Close()

	var (
		r         Resp
		addresses []string
	)

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
		return []string{}, nil
	}

	for _, data := range r.Data {
		addresses = append(addresses, data.WalletAddress)
	}

	return addresses, nil
}
