package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/leekchan/accounting"
)

type FinancialsAsReported struct {
	Cik  string `json:"cik"`
	Data []struct {
		AcceptedDate string `json:"acceptedDate"`
		AccessNumber string `json:"accessNumber"`
		Cik          string `json:"cik"`
		EndDate      string `json:"endDate"`
		FiledDate    string `json:"filedDate"`
		Form         string `json:"form"`
		Quarter      int64  `json:"quarter"`
		Report       struct {
			Bs []struct {
				Concept string `json:"concept"`
				Label   string `json:"label"`
				Unit    string `json:"unit"`
				Value   int64  `json:"value"`
			} `json:"bs"`
			Cf []struct {
				Concept string `json:"concept"`
				Label   string `json:"label"`
				Unit    string `json:"unit"`
				Value   int64  `json:"value"`
			} `json:"cf"`
			Ic []struct {
				Concept string `json:"concept"`
				Label   string `json:"label"`
				Unit    string `json:"unit"`
				Value   int64  `json:"value"`
			} `json:"ic"`
		} `json:"report"`
		StartDate string `json:"startDate"`
		Symbol    string `json:"symbol"`
		Year      int64  `json:"year"`
	} `json:"data"`
	Symbol string `json:"symbol"`
}

func main() {
	url := "https://finnhub.io/api/v1/stock/financials-reported?symbol=INTC&from=2022-01-01"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-finnhub-token", "d6g8ec9r01qt4931ub00d6g8ec9r01qt4931ub0g")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func FmtMoney(number float64) {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	money := ac.FormatMoney(number)
	fmt.Println(money)
}
