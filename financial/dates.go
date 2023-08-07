package financial

import (
	"fmt"
	"strconv"
	"time"

	xmock "github.com/mt1976/appFrame/mock"
)

type Date struct {
	Code     string
	Name     string
	Date     time.Time
	Sort     int
	Simple   string
	External string
	Human    string
}

func getTenorDate(tenor Tenor, tradeDate time.Time, ccy ...string) (time.Time, error) {

	fmt.Printf("Tenor [%v] Trade Date [%v] Ccys %v [%v]\n", tenor.String(), tradeDate.Format("2006-01-02"), ccy, len(ccy))

	// Calculate the settlement days, and adjust the date based on the term string provided i.e. 1D, 1W, 1M, 1Y
	// loop thgouth currencies

	if !xmock.IsValidPeriod(tenor.String()) {
		return time.Now(), fmt.Errorf("invalid tenor [%s]", tenor.String())
	}

	spotDays := 0

	for _, c := range ccy {
		ccySpot, spotError := getSettlementDaysCCY(c)
		if spotError != nil {
			return time.Now(), spotError
		}
		if ccySpot > spotDays {
			spotDays = ccySpot
		}
	}

	switch tenor.term {
	case "SP":
		return tradeDate.AddDate(0, 0, spotDays), nil
	case "ON":
		return tradeDate.AddDate(0, 0, spotDays-1), nil
	case "TN":
		return tradeDate.AddDate(0, 0, 1), nil
	case "TD":
		return tradeDate, nil
	}

	tenorPeriod, err := tenorToDuration(tenor)
	if err != nil {
		return time.Now(), err
	}

	rtn := tradeDate.AddDate(0, 0, spotDays).Add(tenorPeriod)
	return rtn, nil
}

func getTenorDateTEST(tenor Tenor, tradeDate time.Time, ccy ...xmock.Currency) (time.Time, error) {
	return time.Now(), nil
}

func tenorToDuration(tenor Tenor) (time.Duration, error) {
	term := tenor.String()
	if len(term) < 2 {
		return 0, fmt.Errorf("invalid term length [%s]", term)
	}

	valueStr := term[:len(term)-1]
	unit := term[len(term)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid term prefix [%s]", term)
	}

	switch unit {
	case 'D':
		return time.Duration(value) * 24 * time.Hour, nil
	case 'W':
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case 'M':
		return time.Duration(value) * 30 * 24 * time.Hour, nil // Assuming 30 days per month
	case 'Y':
		return time.Duration(value) * 365 * 24 * time.Hour, nil // Assuming 365 days per year
	default:
		return 0, fmt.Errorf("invalid term unit: %c", unit)
	}
}

func getLadderCCY(pivotDate time.Time, ccy string) []Date {
	var DateList []Date
	rateLadder := xmock.Ladder
	xmock.RateValueToString(rateLadder)
	fmt.Printf("rateLadder: %v\n", rateLadder)
	fmt.Printf("ccy: %v\n", ccy)
	fmt.Printf("pivotDate: %v\n", pivotDate.Format("2006-01-02"))
	// range over the ladder
	for _, ladder := range rateLadder {
		thisTenor, err := NewTenor(ladder.Code)
		if err != nil {
			fmt.Printf("Error [%v]\n", err)
		}
		date, err := getTenorDate(thisTenor, pivotDate, ccy)
		if err != nil {
			fmt.Printf("Error [%v]\n", err)
		}
		fmt.Printf("thisTenor: [%v] [%v] -> [%v]\n", ladder.Code, thisTenor.String(), date.Format("2006-01-02"))
		di := Date{}
		di.Code = ladder.Code
		di.Name = ladder.Name
		di.Date = date
		di.Sort = ladder.Index
		di.Simple = date.Format("01/02/2006")
		di.External = date.Format("2006-01-02")
		di.Human = date.Format("Mon 02 Jan 2006")
		DateList = append(DateList, di)
	}
	return DateList
}

func getLadder(pivotdate time.Time, ccy ...string) []Date {
	// TODO - this is a stub, please write logic
	return nil
}

func getTenorFromDateCCY(inDate, pivotDate time.Time, ccy string) (Tenor, error) {
	// TODO - this is a stub, please write logic
	return Tenor{}, nil
}

func getTenorFromDateCCYPAIR(inDate, pivotDate time.Time, ccy1 string, ccy2 string) (Tenor, error) {
	// TODO - this is a stub, please write logic
	return Tenor{}, nil
}

func getTenorFromDateCCYCROSS(inDate, pivotDate time.Time, ccy1 string, via string, ccy2 string) (Tenor, error) {
	// TODO - this is a stub, please write logic
	return Tenor{}, nil
}

func getTenorFromDate(inDate, pivotDate time.Time, ccy ...string) (Tenor, error) {
	// TODO - this is a stub, please write logic
	return Tenor{}, nil
}
