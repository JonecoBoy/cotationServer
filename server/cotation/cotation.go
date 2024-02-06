package cotation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type errorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AwesomeAPICotation struct {
	Cotation `json:"-"`
}

type Cotation struct {
	From      string `json:"code"`
	To        string `json:"codein"`
	Name      string `json:"name"`
	High      string `json:"high"`
	Low       string `json:"low"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	CreatedAt string `json:"-"`
}

const requestExpirationTime = 200 * time.Millisecond

func GetCotation(from string, to string) (Cotation, error) {
	cotationPair := strings.ToUpper(from + "-" + to)

	ctx := context.Background()
	// o contexto expira em 1 segundo!
	ctx, cancel := context.WithTimeout(ctx, requestExpirationTime)
	defer cancel() // de alguma forma nosso contexto será cancelado
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/"+cotationPair, nil)
	if err != nil {
		return Cotation{}, err
	}

	// faz a request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return Cotation{}, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Api fetch timeout exceeed.")
		return Cotation{}, errors.New("Api fetch timeout exceeed.")
	}

	// depois de tudo termina e faz o body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Cotation{}, err
	}

	if resp.StatusCode >= 300 {
		var businessError errorResponse
		err = json.Unmarshal(body, &businessError)
		if err != nil {
			return Cotation{}, err
		}
		return Cotation{}, errors.New(businessError.Message)
	} else {

		c, err := parseJSONToCotation(string(body))
		if err != nil {
			return Cotation{}, err
		}
		//cotation := cotationRaw.Cotation
		return c, err
	}

}

func parseJSONToCotation(jsonData string) (Cotation, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return Cotation{}, err
	}

	// Assuming there is only one key in the map
	for _, value := range data {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			return Cotation{
				From:      nestedMap["code"].(string),
				To:        nestedMap["codein"].(string),
				Name:      nestedMap["name"].(string),
				High:      nestedMap["high"].(string),
				Low:       nestedMap["low"].(string),
				VarBid:    nestedMap["varBid"].(string),
				PctChange: nestedMap["pctChange"].(string),
				Bid:       nestedMap["bid"].(string),
				Ask:       nestedMap["ask"].(string),
				CreatedAt: "", // You may need to handle the timestamp field accordingly
			}, nil
		}
		break
	}

	return Cotation{}, fmt.Errorf("no valid data found in JSON")
}
