# Receipt Processor Solution

This is a Go-based implementation of the Receipt Processor challenge. The service provides two endpoints for processing receipts and calculating points based on the specified rules.

## Requirements

- Go 1.22 or higher

## Running the Application

1. Clone the repository
2. Navigate to the project directory
3. Run the application:
   ```bash
   go run main.go
   ```
   The server will start on port 8080.

## API Endpoints

### Process Receipt
- **POST** `/receipts/process`
- Accepts a receipt JSON in the request body
- Returns a unique ID for the receipt

Example request using PowerShell:
```powershell
# Using a JSON file
curl -X POST -H "Content-Type: application/json" -d "@examples/target-receipt.json" http://localhost:8080/receipts/process

# Using inline JSON
curl -X POST -H "Content-Type: application/json" -d '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    }
  ],
  "total": "6.49"
}' http://localhost:8080/receipts/process
```

### Get Points
- **GET** `/receipts/{id}/points`
- Returns the points awarded for the receipt

Example request:
```powershell
curl http://localhost:8080/receipts/{id}/points
```

## Example Receipts

The `examples` directory contains two sample receipts with their expected points:

1. `target-receipt.json`: A receipt from Target with 5 items
   - Expected points: 28
   - Breakdown:
     - 6 points for retailer name (Target = 6 alphanumeric chars)
     - 10 points for 2 pairs of items (5 items = 2 pairs)
     - 6 points for odd day (January 1st)
     - 6 points for "Klarbrunn 12-PK 12 FL OZ" (length 20, multiple of 3, price $12.00 * 0.2 = 2.4, rounded up to 3)

2. `mm-receipt.json`: A receipt from M&M Corner Market with 4 Gatorades
   - Expected points: 109
   - Breakdown:
     - 14 points for retailer name (M&M Corner Market = 14 alphanumeric chars)
     - 50 points for round dollar amount ($9.00)
     - 25 points for total being a multiple of 0.25
     - 10 points for 2 pairs of items (4 items = 2 pairs)
     - 10 points for time between 2:00 PM and 4:00 PM (14:33)

## Testing

The project includes comprehensive test coverage:

1. Unit Tests for Points Calculation (`service/points_test.go`):
   - Tests all individual point calculation rules
   - Validates points for example receipts
   - Run with: `go test ./service`

2. API Integration Tests (`handlers/handlers_test.go`):
   - Tests receipt processing endpoint
   - Tests points retrieval endpoint
   - Tests error handling
   - Run with: `go test ./handlers`

Run all tests with:
```bash
go test ./...
```

## Implementation Details

The application is structured into several packages:
- `models`: Contains data structures for receipts and responses
- `service`: Contains the points calculation logic
- `handlers`: Contains HTTP handlers for the API endpoints

Data is stored in memory using a thread-safe map structure.
