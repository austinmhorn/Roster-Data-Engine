package main

import (
	"Roster-Data-Engine/notionapi" // Matches `go.mod`
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

func writeLastUpdated() {
	// Get current timestamp
	timestamp := time.Now().Format("01/02/2006 - 15:04:05")

	// Write timestamp to a second CSV file
	lastUpdatedFile := "last_updated.csv"
	csvFile, err := os.Create(lastUpdatedFile)
	if err != nil {
		fmt.Println("❌ ERROR: Creating Last Updated CSV file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"Last Updated", timestamp})

	fmt.Printf("✅ Timestamp successfully written to %s!\n", lastUpdatedFile)
}

func main() {
	// Load configuration first
	err := notionapi.LoadConfig()
	if err != nil {
		fmt.Println("❌ Error loading config:", err)
		return
	}

	fmt.Println("🚀 Fetching Notion Data...")
	data, err := notionapi.FetchNotionData()
	if err != nil {
		fmt.Println("❌ Error fetching data from Notion:", err)
		return
	}

	totalEntries := len(data)
	if totalEntries == 0 {
		fmt.Println("❌ No data found in Notion.")
		return
	}

	fmt.Printf("✅ Successfully retrieved %d entries from Notion.\n", totalEntries)

	// Save data to CSV
	csvFileName := "notion_data.csv"
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		fmt.Println("❌ ERROR: Creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	headers := []string{
		"Employee ID",
		"Name",
		"Status",
		"Property",
		"Title",
		"Hire Date",
		"Term Date",
		"Tags",
		"EMP Lease",
		"Entrata",
		"PHQ",
		"Outlook",
		"IT",
		"Entrata Group Routing",
		"Knock",
		"AIRM",
		"POPIC",
		"Leo247",
		"Snappt",
		"AptList Dashboard",
		"NetVendor",
		"ApartmentIQ",
		"Spruce",
		"The Guarantors",
		"Community Rewards",
		"BillSpend",
		"TBD: Pay Ready",
		"TBD: Online Print Marketing Store",
		"Regional Manager",
		"Department",
		"Company",
	}
	writer.Write(headers)

	// Process each entry
	for i, entry := range data {
		props, ok := entry["properties"].(map[string]interface{})
		if !ok {
			fmt.Printf("⚠️ Skipping entry %d: Invalid properties format.\n", i+1)
			continue
		}

		employeeIDStr := notionapi.GetCleanPlainTextValue(props, "Employee ID")
		nameStr := notionapi.GetName(props, "Name")
		statusStr := notionapi.GetStatus(props, "Status")
		propertyStr := notionapi.GetRollupFormulaString(props, "Property (As Text)")
		titleStr := notionapi.GetRollupFormulaString(props, "Title (As Text)")
		hireDateStr := notionapi.GetDateValue(props, "Hire Date")
		termDateStr := notionapi.GetDateValue(props, "Term Date")
		tagsStr := notionapi.GetMultiSelectStrings(props, "Tags")
		// Convert tags slice to a comma-separated string
		tagsStrFormatted := strings.Join(tagsStr, ", ")
		empLeaseStr := notionapi.GetSelectValue(props, "EMP Lease")
		entrataStr := notionapi.GetDateValue(props, "Entrata")
		phqStr := notionapi.GetDateValue(props, "PHQ")
		outlookStr := notionapi.GetDateValue(props, "Outlook")
		itStr := notionapi.GetDateValue(props, "IT")
		entrataGroupRoutingStr := notionapi.GetCleanPlainTextValue(props, "Entrata Group Routing")
		knockStr := notionapi.GetCleanPlainTextValue(props, "Knock")
		airmStr := notionapi.GetSelectValue(props, "AIRM")
		popicStr := notionapi.GetSelectValue(props, "POPIC")
		leo247Str := notionapi.GetSelectValue(props, "Leo247")
		snapptStr := notionapi.GetCleanPlainTextValue(props, "Snappt")
		aptListDashboardStr := notionapi.GetSelectValue(props, "AptList Dashboard")
		netVendorStr := notionapi.GetSelectValue(props, "NetVendor")
		aparmentIQStr := notionapi.GetSelectValue(props, "ApartmentIQ")
		spruceStr := notionapi.GetSelectValue(props, "Spruce")
		theGuarantorsStr := notionapi.GetSelectValue(props, "The Guarantors")
		communityRewardsStr := notionapi.GetSelectValue(props, "Community Rewards")
		billSpendStr := notionapi.GetSelectValue(props, "BillSpend")
		tbdPayReadyStr := notionapi.GetSelectValue(props, "TBD: Pay Ready")
		tbdOnlinePrintMarketingStoreStr := notionapi.GetSelectValue(props, "TBD: Online Print Marketing Store")
		regionalManagerStr := notionapi.GetRollupFormulaString(props, "Regional Manager (As Text)")
		departmentStr := notionapi.GetRollupSelectValue(props, "Department")
		companyStr := notionapi.GetSelectValue(props, "Company")

		row := []string{
			employeeIDStr,
			nameStr,
			statusStr,
			propertyStr,
			titleStr,
			hireDateStr,
			termDateStr,
			tagsStrFormatted,
			empLeaseStr,
			entrataStr,
			phqStr,
			outlookStr,
			itStr,
			entrataGroupRoutingStr,
			knockStr,
			airmStr,
			popicStr,
			leo247Str,
			snapptStr,
			aptListDashboardStr,
			netVendorStr,
			aparmentIQStr,
			spruceStr,
			theGuarantorsStr,
			communityRewardsStr,
			billSpendStr,
			tbdPayReadyStr,
			tbdOnlinePrintMarketingStoreStr,
			regionalManagerStr,
			departmentStr,
			companyStr}
		writer.Write(row)

		// Print progress
		fmt.Printf("📊 Progress: %d/%d (%.2f%%)\n", i+1, totalEntries, float64(i+1)/float64(totalEntries)*100)
	}

	fmt.Printf("✅ Data successfully written to %s!\n", csvFileName) // Save data to CSV

	// Write last updated timestamp
	writeLastUpdated()
}
