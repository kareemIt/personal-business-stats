package printops

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	// "github.com/kareemItani/personal-business-stats/services/printops/model"
	"github.com/kareemItani/personal-business-stats/util"
)

type PrintOpsAPI struct {
	APIKey           string
	StoreID          string
	Printops_api_url string
}

func GETorders() (interface{}, error) {
	APIKey := os.Getenv("PRINTOPS_API_KEY")
	storeID := os.Getenv("PRINTOPS_STORE_ID")
	baseURL := os.Getenv("PRINTOPS_API_URL")
	url := fmt.Sprintf("%s/%s/orders?page=1&limit=500", baseURL, storeID)

	responseStr, err := util.MakeGetAPICall(APIKey, url)
	if err != nil {
		return nil, err
	}

	var payload interface{}
	if err := json.Unmarshal([]byte(responseStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON from PrintOps: %w", err)
	}

	return payload, nil
}

func NormalizingAPIResponse(response interface{}) (map[string]int, error) {
	top := response
	if m, ok := response.(map[string]interface{}); ok {
		if d, ok := m["data"]; ok {
			top = d
		}
	}

	var rawOrders []interface{}
	switch v := top.(type) {
	case map[string]interface{}:
		if o, ok := v["orders"].([]interface{}); ok {
			rawOrders = o
		}
	case []interface{}:
		rawOrders = v
	}
	if len(rawOrders) == 0 {
		return nil, fmt.Errorf("no orders found in response; did you unwrap data?")
	}

	result := make(map[string]int)
	for _, item := range rawOrders {
		ord, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		if state, _ := ord["state"].(string); state == "cancelled" || state == "production" {
			continue
		}

		cents := 0
		if ot, ok := ord["order_total"].(map[string]interface{}); ok {
			if gt, ok := ot["grand_total"].(float64); ok {
				cents = int(gt)
			}
		}
		price := cents

		// Find first non-voided shipment & grab its tracking number
		tracking := ""
		if shipWrapper, ok := ord["shipments"].(map[string]interface{}); ok {
			if ships, ok := shipWrapper["shipments"].([]interface{}); ok {
				for _, s := range ships {
					sMap, ok := s.(map[string]interface{})
					if !ok {
						continue
					}
					if status, _ := sMap["status"].(string); status == "cancelled" || status == "pending" {
						continue
					}
					if vOn := sMap["voided_on"]; vOn != nil {
						continue
					}
					if tn, ok := sMap["tracking_number"].(string); ok {
						tracking = tn
						break
					}
				}
			}
		}

		if tracking != "" {
			result[tracking] = -(price)
		}
	}

	return result, nil
}
