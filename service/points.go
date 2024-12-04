package service

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jontitorr/receipt-processor/models"
)

func CalculatePoints(receipt *models.Receipt) int64 {
	var points int64

	// Rule 1: One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points++
		}
	}

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if math.Mod(total*100, 25) == 0 {
			points += 25
		}
	}

	// Rule 4: 5 points for every two items on the receipt
	points += int64((len(receipt.Items) / 2) * 5)

	// Rule 5: Points for items with description length multiple of 3
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int64(math.Ceil(price * 0.2))
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	if purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if purchaseDate.Day()%2 == 1 {
			points += 6
		}
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		hour := purchaseTime.Hour()
		minute := purchaseTime.Minute()
		if (hour == 14) || (hour == 15 && minute == 0) {
			points += 10
		}
	}

	return points
}
