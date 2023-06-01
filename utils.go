package main

import (
	"fmt"
	"strconv"
)

// Calculate entry price for the long order.
func calcEntryPrice(price string) string {
	pricef, _ := strconv.ParseFloat(price, 64)
	// pricef = pricef - 10
	pricef = pricef + 1
	return fmt.Sprintf("%.2f", pricef)
}

// Calculate stop loss price for the short order.
func calcShortOrderStopLoss(price string, slPercentage float64) (string, error) {
	pricef, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", err
	}
	pricef = pricef * (1 + slPercentage)
	return fmt.Sprintf("%2.f", pricef), nil
}

func calcShortOrderTakeProfit(price string, tpPercentage float64) (string, error) {
	pricef, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", err
	}
	pricef = pricef * (1 - tpPercentage)
	return fmt.Sprintf("%2.f", pricef), nil
}
