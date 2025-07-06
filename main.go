package main

import (
	"fmt"

	"github.com/kareemItani/personal-business-stats/services/etsy"
	"github.com/kareemItani/personal-business-stats/services/printops"
)

func main() {
	OrdersMap := SetupPrintOps()
	SetupEtsy(OrdersMap)
}

func SetupPrintOps() map[string]int {
	PrintOpsOrders, err := printops.GETorders()
	if err != nil {
		fmt.Println("Error getting orders:", err)
		return nil
	}
	PrintOpsOrdersNormalized, err := printops.NormalizingAPIResponse(PrintOpsOrders)
	fmt.Println("printops orders normalize ")
	for i, job := range PrintOpsOrdersNormalized {
		fmt.Printf("Index: %s, PrintJob: %+v\n", i, job)
	}

	return PrintOpsOrdersNormalized
}

func SetupEtsy(EtsyObject map[string]int) {
	feesMap, err := etsy.GetTrackingToFeesMap(EtsyObject)
	if err != nil {
		fmt.Println("Error fetching tracking to fees map:", err)
		return
	}
	fmt.Println("Etsy Fees Map: ")
	totalProfit := 0
	for tracking, total := range feesMap {
		fmt.Printf("Tracking: %s, Net Profit: %+v\n", tracking, total)
		totalProfit += total
	}
	totalProfit = totalProfit / 100
	totalFees := len(feesMap) * 3
	fmt.Println("Total Orders:", len(feesMap))
	fmt.Println("Profit per order average:", totalProfit/len(feesMap))
	fmt.Println("Total Estimated Profited:", totalProfit-totalFees)

}
