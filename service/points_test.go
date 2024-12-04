package service

import (
	"testing"

	"github.com/jontitorr/receipt-processor/models"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int64
	}{
		{
			name: "Target Receipt",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			expected: 28, // 6 (retailer) + 10 (2 pairs of items) + 6 (odd day) + 6 (Klarbrunn price * 0.2)
		},
		{
			name: "M&M Corner Market Receipt",
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			expected: 109, // 14 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) + 10 (2 pairs of items) + 10 (time between 2-4)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := CalculatePoints(&tt.receipt)
			if points != tt.expected {
				t.Errorf("CalculatePoints() = %v, want %v", points, tt.expected)
			}
		})
	}
}

func TestCalculatePointsRules(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		rule     string
		expected int64
	}{
		{
			name: "Retailer name alphanumeric characters",
			receipt: models.Receipt{
				Retailer:     "Target123",
				PurchaseDate: "",
				PurchaseTime: "",
				Items:        []models.Item{},
				Total:       "",
			},
			rule:     "retailer alphanumeric",
			expected: 9,
		},
		{
			name: "Round dollar amount",
			receipt: models.Receipt{
				Total: "35.00",
				Retailer:     "",
				PurchaseDate: "",
				PurchaseTime: "",
				Items:        []models.Item{},
			},
			rule:     "round dollar",
			expected: 75, // 50 for round dollar + 25 for multiple of 0.25
		},
		{
			name: "Multiple of 0.25",
			receipt: models.Receipt{
				Total: "35.25",
			},
			rule:     "quarter multiple",
			expected: 25,
		},
		{
			name: "Items pairs",
			receipt: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "1.00"},
					{ShortDescription: "Item2", Price: "2.00"},
					{ShortDescription: "Item3", Price: "3.00"},
				},
			},
			rule:     "items pairs",
			expected: 5,
		},
		{
			name: "Description length multiple of 3",
			receipt: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "ABC", Price: "2.00"},        // length 3
					{ShortDescription: "ABCDEF", Price: "3.00"},     // length 6
					{ShortDescription: "A B C", Price: "5.00"},      // length 5 (not multiple of 3)
				},
				Retailer:     "",
				PurchaseDate: "",
				PurchaseTime: "",
				Total:       "",
			},
			rule:     "description length",
			expected: 7, // (2.00 * 0.2 = 0.4 rounds to 1) + (3.00 * 0.2 = 0.6 rounds to 1) + 5 points for 2 items
		},
		{
			name: "Odd day purchase",
			receipt: models.Receipt{
				PurchaseDate: "2022-01-01",
			},
			rule:     "odd day",
			expected: 6,
		},
		{
			name: "Time between 2:00pm and 4:00pm",
			receipt: models.Receipt{
				PurchaseTime: "14:30",
			},
			rule:     "time range",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := CalculatePoints(&tt.receipt)
			if points != tt.expected {
				t.Errorf("%s: got %v points, want %v", tt.rule, points, tt.expected)
			}
		})
	}
}
