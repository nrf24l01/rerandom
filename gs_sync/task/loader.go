package task

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/nrf24l01/rerandom/gs_sync/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func DownloadSheet(spreadsheetId string, targetRange string, cfg config.GSConfig) *Sheet {
	// Convert GSConfig to JSON for JWTConfigFromJSON
	configJSON, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("Unable to marshal config: %v", err)
	}

	config, err := google.JWTConfigFromJSON(configJSON, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse service account config: %v", err)
	}

	ctx := context.Background()
	client := config.Client(ctx)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to create Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, targetRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		log.Println("No data found.")
		return &Sheet{}
	}

	var rows []SheetRow

	// Parse each row from the Google Sheets response
	for i, row := range resp.Values {
		if i == 0 {
			// Skip header row if exists
			continue
		}

		// Ensure the row has enough columns
		if len(row) < 8 {
			continue
		}

		// Parse the data according to CSV structure
		var sheetRow SheetRow

		// Column 0: Fraction (доля)
		if fraction, err := strconv.ParseUint(toString(row[0]), 10, 32); err == nil {
			sheetRow.Fraction = uint(fraction)
		}

		// Column 2: FractionFrom (от)
		if fractionFrom, err := strconv.ParseUint(toString(row[2]), 10, 32); err == nil {
			sheetRow.FractionFrom = uint(fractionFrom)
		}

		// Column 3: FractionTo (до)
		if fractionTo, err := strconv.ParseUint(toString(row[3]), 10, 32); err == nil {
			sheetRow.FractionTo = uint(fractionTo)
		}

		// Column 4: Alive (жив ли)
		if alive, err := strconv.ParseBool(toString(row[4])); err == nil {
			sheetRow.Alive = alive
		}

		// Column 5: ID
		if id, err := strconv.ParseUint(toString(row[5]), 10, 32); err == nil {
			sheetRow.Id = uint(id)
		}

		// Column 6: LastName (фамилия)
		sheetRow.LastName = toString(row[6])

		// Column 7: FirstName (имя)  
		if len(row) > 7 {
			sheetRow.FirstName = toString(row[7])
		}

		rows = append(rows, sheetRow)
	}

	// Return the Sheet with parsed data
	return &Sheet{
		ClearUsers:    rows,
		ModifiedUsers: rows, // Initialize as empty
		Actions:       []Action{},
	}
}

// Helper function to safely convert interface{} to string
func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}