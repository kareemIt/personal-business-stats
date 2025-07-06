package util

import (
	"strconv"
)

func GetProfit(sold string, priceToMake string) float64 {
	soldFloat, err := strconv.ParseFloat(sold, 64)
	if err != nil {
		return 0
	}
	priceToMakeFloat, err := strconv.ParseFloat(priceToMake, 64)
	if err != nil {
		return 0
	}
	return soldFloat - priceToMakeFloat
}

func GetPercentageOfProfit(Profit string, priceToMake string) float64 {
	profitFloat, err := strconv.ParseFloat(Profit, 64)
	if err != nil {
		return 0
	}
	priceToMakeFloat, err := strconv.ParseFloat(priceToMake, 64)
	if err != nil || priceToMakeFloat == 0 {
		return 0
	}
	return profitFloat / priceToMakeFloat
}
