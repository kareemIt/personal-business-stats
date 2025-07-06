package etsy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/kareemItani/personal-business-stats/util"
)

type etsyAPI struct {
	APIKey       string
	storeID      string
	etsy_api_url string
}

type FeeResponse struct {
	TotalFees float64 `json:"total_fees"`
}

// ExchangeCodeForToken exchanges the authorization code for an access token
func ExchangeCodeForToken() (*EtsyTokenResponse, error) {
	clientID := os.Getenv("ETSY_API_KEY")
	redirectURI := os.Getenv("ETSY_REDIRECT_URI")
	code := os.Getenv("ETSY_AUTH_ID")
	tokenURL := "https://api.etsy.com/v3/public/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)
	data.Set("code_verifier", "vvkdljkejllufrvbhgeiegrnvufrhvrffnkvcknjvfid")

	headers := map[string]string{}

	respStr, err := util.MakePostFormAPICall(tokenURL, data, headers)
	if err != nil {
		return nil, err
	}

	var tokenResp EtsyTokenResponse
	if err := json.Unmarshal([]byte(respStr), &tokenResp); err != nil {
		return nil, fmt.Errorf("invalid JSON from Etsy token exchange: %w", err)
	}
	return &tokenResp, nil
}

// GetEtsyOrders fetches receipts (orders) from the Etsy API using OAuth2 access token.
func GetEtsyOrders(accessToken, shopID string, clientID string, EtsyMap map[string]int) (map[string]int, error) {
	baseURL := "https://api.etsy.com/v3/application/shops"
	limit := 100
	offset := 0
	url := fmt.Sprintf("%s/%s/receipts?includes=Transactions,Shippings&limit=%d&offset=%d&client_id=%s", baseURL, shopID, limit, offset, clientID)

	responseStr, err := util.MakeGetAPICall(accessToken, url)
	if err != nil {
		return nil, err
	}

	var payload interface{}
	if err := json.Unmarshal([]byte(responseStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON from Etsy: %w", err)
	}

	parsed := ParsingEtsyOrders(EtsyMap, payload)

	return parsed, nil
}

func ParsingEtsyOrders(EtsyMap map[string]int, payload interface{}) map[string]int {
	fmt.Println("Parsing Etsy Orders...")

	payloadMap, ok := payload.(map[string]interface{})
	if !ok {
		return EtsyMap
	}
	results, ok := payloadMap["results"].([]interface{})
	if !ok {
		return EtsyMap
	}
	fmt.Println("Number of results:", len(results))

	for _, order := range results {
		orderMap, ok := order.(map[string]interface{})
		if !ok {
			continue
		}

		// Get the tracking number (shipping number)
		trackingNumber := ""
		if shipping, ok := orderMap["shipments"].([]interface{}); ok && len(shipping) > 0 {
			shipMap, ok := shipping[0].(map[string]interface{})
			if ok {
				if tn, ok := shipMap["tracking_code"].(string); ok {
					trackingNumber = tn
				}
			}
		} else if tn, ok := orderMap["tracking_code"].(string); ok {
			trackingNumber = tn
		}

		if trackingNumber == "" {
			continue
		}

		// Check if tracking number is in the map
		if _, exists := EtsyMap[trackingNumber]; !exists {
			fmt.Println("Tracking number not found in map, skipping:", trackingNumber)
			continue
		}

		// Get grand total
		grandTotal := 0
		if gtRaw, ok := orderMap["grandtotal"].(map[string]interface{}); ok {
			// Etsy always returns numbers as float64
			if amtF, ok := gtRaw["amount"].(float64); ok {
				// Etsy’s “amount” is already in cents (e.g. 1999 == $19.99)
				grandTotal = int(amtF)
			}
		}

		// Subtract fee from grand total and update the map
		net := grandTotal + EtsyMap[trackingNumber]
		fmt.Println("net: ", net, " grandTotal: ", grandTotal, " EtsyMap[trackingNumber]: ", EtsyMap[trackingNumber])
		EtsyMap[trackingNumber] = net
	}

	return EtsyMap
}

// GetTrackingToFeesMap fetches all orders and maps tracking numbers to their total fees.
func GetTrackingToFeesMap(EtsyMap map[string]int) (map[string]int, error) {
	clientID := os.Getenv("ETSY_API_KEY")
	shopID := os.Getenv("ETSY_STORE_ID")

	// Step 1: Get the bearer token
	tokenResp, err := ExchangeCodeForToken()
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	accessToken := tokenResp.AccessToken

	// Step 2: Get all orders
	orders, err := GetEtsyOrders(accessToken, shopID, clientID, EtsyMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	// Convert map[string]float64 to map[string]int
	intFeesMap := make(map[string]int)
	for k, v := range orders {
		intFeesMap[k] = int(v)
	}

	return intFeesMap, nil
}
