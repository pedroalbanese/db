package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	selectedFile := flag.String("f", "", "Select CSV file")
	listFlag := flag.Bool("list", false, "List CSV files")
	createFlag := flag.Bool("create", false, "Create CSV file")
	searchFlag := flag.Bool("search", false, "Search entries in one or more CSV")
	idFlag := flag.Int("id", 0, "ID (for get command)")
	columnFlag := flag.String("column", "", "Column name (for get command)")
	getFlag := flag.Bool("get", false, "Get value")
	editFlag := flag.Bool("edit", false, "Edit entry")

	var headers []string
	for {
		flag.Parse()

		if len(os.Args) > 1 {
			if *listFlag && *selectedFile != "" {
				var err error
				headers, err = readHeaders(*selectedFile)
				if err != nil {
					fmt.Println("Error reading headers:", err)
				} else {
					listRecords2(*selectedFile, headers)
					os.Exit(0)
				}
				os.Exit(0)
			} else if *listFlag {
				listCSVFiles()
				os.Exit(0)
			} else if *createFlag {
				createNewCSV()
				os.Exit(0)
			} else if *searchFlag && *selectedFile != "" {
				var err error
				headers, err = readHeaders(*selectedFile)
				if err != nil {
					fmt.Println("Error reading headers:", err)
				} else {
					searchRecord(*selectedFile, headers)
					os.Exit(0)
				}
			} else if *searchFlag && *selectedFile == "" {
				searchAndDisplayRecords()
				os.Exit(0)
			} else if *getFlag && *idFlag > 0 {
				value, err := getValue(*idFlag, *selectedFile, *columnFlag)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				if *columnFlag != "" {
					fmt.Printf("Value for ID '%d' in column '%s': %s\n", *idFlag, *columnFlag, value)
				} else {
					fmt.Printf("Value for ID '%d': %s\n", *idFlag, value)
				}
				os.Exit(0)
			} else if *editFlag && *selectedFile != "" {
				var err error
				headers, err = readHeaders(*selectedFile)
				if err != nil {
					fmt.Println("Error reading headers:", err)
				} else {
					editRecord(*selectedFile, headers)
					os.Exit(0)
				}
				os.Exit(0)
			}
		} else {
			fmt.Println("CSV-Based Database Manager")
			fmt.Println(" 1. List CSV Files")
			fmt.Println(" 2. Select CSV File")
			fmt.Println(" 3. Create CSV File")
			fmt.Println(" 4. Search Records")
			fmt.Println(" 5. Exit")
			fmt.Print("Choose an option: ")

			var choice int
			fmt.Scanln(&choice)

			switch choice {
			case 1:
				listCSVFiles()
			case 2:
				selectCSVFile()
			case 3:
				createNewCSV()
			case 4:
				searchAndDisplayRecords()
			case 5:
				os.Exit(0)
			default:
				fmt.Println("Invalid option.")
			}
		}
	}
}

func listCSVFiles() {
	files, err := listFiles(".csv")
	if err != nil {
		fmt.Println("Error listing CSV files:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No CSV files found in the current directory and its subdirectories.")
		return
	}

	fmt.Println("CSV Files:")
	for i, file := range files {
		fmt.Printf(" %d. %s\n", i+1, file)
	}
}

func getValue(id int, file, column string) (string, error) {
	csvFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	headers, err := reader.Read()
	if err != nil {
		return "", err
	}

	colIndex := -1
	if column != "" {
		for i, header := range headers {
			if header == column {
				colIndex = i
				break
			}
		}

		if colIndex == -1 {
			return "", fmt.Errorf("Column '%s' not found in the CSV file", column)
		}
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
		} else if err != nil {
			return "", err
		}

		rowID, err := strconv.Atoi(row[0])
		if err != nil {
			return "", err
		}

		if rowID == id {
			if column != "" {
				if colIndex < len(row) {
					return row[colIndex], nil
				}
				return "", fmt.Errorf("Column '%s' not found in the row with ID %d", column, id)
			}
			var formattedRow strings.Builder
			maxKeyWidth := 0
			for i := 1; i < len(headers); i++ {
				if len(headers[i]) > maxKeyWidth {
					maxKeyWidth = len(headers[i])
				}
			}
			for i := 1; i < len(headers); i++ {
				key := headers[i]
				value := row[i]
				formattedRow.WriteString(fmt.Sprintf("%-*s : %s\n", maxKeyWidth, key, value))
			}
			return fmt.Sprintf("\n%s\n%s%s", strings.Repeat("=", 80), formattedRow.String(), strings.Repeat("-", 80)), nil
		}
	}

	return "", fmt.Errorf("No row found with ID %d", id)
}

func selectCSVFile() {
	files, err := listFiles(".csv")
	if err != nil {
		fmt.Println("Error listing CSV files:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No CSV files found in the current directory and its subdirectories.")
		return
	}

	fmt.Println("Select a CSV file to manage:")
	for i, file := range files {
		fmt.Printf(" %d. %s\n", i+1, file)
	}

	var fileChoice int
	fmt.Print("Enter the number of the CSV file to manage: ")
	fmt.Scanln(&fileChoice)

	if fileChoice < 1 || fileChoice > len(files) {
		fmt.Println("Invalid choice.")
		return
	}

	selectedFile := files[fileChoice-1]
	fmt.Printf("Selected CSV file: %s\n", selectedFile)

	headers, err := readHeaders(selectedFile)
	if err != nil {
		fmt.Println("Error reading headers:", err)
		return
	}

	for {
		fmt.Println("CSV File Management Menu")
		fmt.Println(" 1. Add Record")
		fmt.Println(" 2. List Records")
		fmt.Println(" 3. List Records as Table")
		fmt.Println(" 4. Search Record")
		fmt.Println(" 5. Edit Record")
		fmt.Println(" 6. Delete Record")
		fmt.Println(" 7. Return to Main Menu")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			addRecord(selectedFile, headers)
		case 2:
			listRecords(selectedFile, headers)
		case 3:
			listRecords2(selectedFile, headers)
		case 4:
			searchRecord(selectedFile, headers)
		case 5:
			editRecord(selectedFile, headers)
		case 6:
			deleteRecord(selectedFile, headers)
		case 7:
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func addRecord(selectedFile string, headers []string) {
	fmt.Println("Add Record")

	file, err := os.OpenFile(selectedFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	recordData := make([]string, len(headers))
	recordData[0] = getNextID(selectedFile)

	for i := 1; i < len(headers); i++ {
		fmt.Printf("%s: ", headers[i])
		if strings.ToLower(headers[i]) == "date" {
			currentDate := time.Now().Format("2006-01-02")
			fmt.Printf(" (Press Enter to use - %s): ", currentDate)
			var fieldValue string
			fmt.Scanln(&fieldValue)

			if fieldValue == "" {
				recordData[i] = currentDate
			} else {
				recordData[i] = fieldValue
			}
		} else {
			var fieldValue string
			fmt.Scanln(&fieldValue)
			recordData[i] = fieldValue
		}
	}

	writer := csv.NewWriter(file)
	if err := writer.Write(recordData); err != nil {
		fmt.Println("Error writing the record:", err)
		return
	}
	writer.Flush()

	fmt.Println("Record added successfully!")
}

func listRecords(selectedFile string, headers []string) {
	fmt.Println("List Records")
	fmt.Println(strings.Repeat("=", 80))

	file, err := os.Open(selectedFile)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the records:", err)
		return
	}

	if len(records) == 0 {
		fmt.Println("No records found.")
		return
	}

	maxHeaderWidth := 0
	for _, header := range headers {
		if len(header) > maxHeaderWidth {
			maxHeaderWidth = len(header)
		}
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		fmt.Printf("ID%*s : %s\n", maxHeaderWidth-2, "", record[0])
		for j := 1; j < len(headers); j++ {
			fmt.Printf("%-*s : %s\n", maxHeaderWidth, headers[j], record[j])
		}
		fmt.Println(strings.Repeat("-", 80))
	}
}

func listRecords2(selectedFile string, headers []string) {
	fmt.Println("List Records")
	fmt.Println(strings.Repeat("=", 80))

	file, err := os.Open(selectedFile)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the records:", err)
		return
	}

	if len(records) == 0 {
		fmt.Println("No records found.")
		return
	}

	columnWidths := make([]int, len(headers))
	for i := 0; i < len(headers); i++ {
		columnWidths[i] = len(headers[i])
	}
	for _, record := range records {
		for i := 0; i < len(record); i++ {
			if len(record[i]) > columnWidths[i] {
				columnWidths[i] = len(record[i])
			}
		}
	}

	headerRow := ""
	for i, header := range headers {
		headerRow += fmt.Sprintf("%-*s", columnWidths[i]+2, header)
	}
	fmt.Println(headerRow)

	for i, record := range records {
		if i == 0 {
			continue
		}
		recordRow := ""
		for j := 0; j < len(record); j++ {
			recordRow += fmt.Sprintf("%-*s", columnWidths[j]+2, record[j])
		}
		fmt.Println(recordRow)
	}
	fmt.Println(strings.Repeat("-", 80))
}

func editRecord(selectedFile string, headers []string) {
	fmt.Println("Edit Record")

	file, err := os.OpenFile(selectedFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the records:", err)
		return
	}

	if len(records) <= 1 {
		fmt.Println("No records to edit.")
		return
	}

	var entryNumber int
	fmt.Print("Enter the number of the entry to edit: ")
	fmt.Scanln(&entryNumber)

	if entryNumber < 1 || entryNumber >= len(records) {
		fmt.Println("Invalid entry number.")
		return
	}

	entry := records[entryNumber]

	fmt.Printf("Editing Record %s:\n", entry[0])

	editedValues := make(map[string]string)

	for i := 1; i < len(headers); i++ {
		currentHeader := headers[i]
		currentValue := entry[i]

		fmt.Printf("%s (Press Enter to keep '%s'): ", currentHeader, currentValue)
		var editedValue string
		fmt.Scanln(&editedValue)

		if editedValue == "" {
			editedValues[currentHeader] = currentValue
		} else {
			editedValues[currentHeader] = editedValue
		}
	}

	for i := 1; i < len(headers); i++ {
		entry[i] = editedValues[headers[i]]
	}

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("Error resetting the file offset:", err)
		return
	}
	if err := file.Truncate(0); err != nil {
		fmt.Println("Error truncating the file:", err)
		return
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		fmt.Println("Error writing the records:", err)
		return
	}
	writer.Flush()

	fmt.Println("Record edited successfully!")
}

func searchRecord(selectedFile string, headers []string) {
	fmt.Println("Search Record")

	file, err := os.Open(selectedFile)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	var term string
	fmt.Print("Enter a search term: ")
	fmt.Scanln(&term)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the records:", err)
		return
	}

	found := false
	maxHeaderWidth := 0

	for _, header := range headers {
		if len(header) > maxHeaderWidth {
			maxHeaderWidth = len(header)
		}
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		for j := 1; j < len(headers); j++ {
			if strings.Contains(record[j], term) {
				if !found {
					fmt.Println(strings.Repeat("=", 80))
				}
				maxHeaderWidth := 0
				maxValueWidth := 0
				for j := 0; j < len(headers); j++ {
					if len(headers[j]) > maxHeaderWidth {
						maxHeaderWidth = len(headers[j])
					}
					if len(record[j]) > maxValueWidth {
						maxValueWidth = len(record[j])
					}
				}
				fmt.Printf("%-*s : %-*s\n", maxHeaderWidth, "ID", maxValueWidth, record[0])
				for k := 1; k < len(headers); k++ {
					fmt.Printf("%-*s : %s\n", maxHeaderWidth, headers[k], record[k])
				}
				fmt.Println(strings.Repeat("-", 80))
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("No records found with the search term:", term)
	}
}

func searchAndDisplayRecords() {
	fmt.Println("Search Records")

	var term string
	fmt.Print("Enter a search term: ")
	fmt.Scanln(&term)

	files, err := listFiles(".csv")
	if err != nil {
		fmt.Println("Error listing CSV files:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No CSV files found in the current directory and its subdirectories.")
		return
	}

	for _, filePath := range files {
		headers, err := readHeaders(filePath)
		if err != nil {
			fmt.Printf("Error reading headers from %s: %v\n", filePath, err)
			continue
		}

		fileContainsTerm := false

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening the file %s: %v\n", filePath, err)
			continue
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Printf("Error reading the records from %s: %v\n", filePath, err)
			continue
		}

		for i, record := range records {
			if i == 0 {
				continue
			}

			for j, field := range record {
				if strings.Contains(field, term) {
					if !fileContainsTerm {
						fmt.Println(strings.Repeat("=", 80))
						fmt.Printf("File   : %s\n", filePath)
						fmt.Println(strings.Repeat("=", 80))
						fileContainsTerm = true
					}
					fmt.Printf("ID     : %s\n", record[0])
					fmt.Printf("Column : %s\n", headers[j])
					fmt.Println(strings.Repeat("-", 80))
				}
			}
		}
	}
}

func deleteRecord(selectedFile string, headers []string) {
	fmt.Println("Delete Record")

	file, err := os.OpenFile(selectedFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the records:", err)
		return
	}

	var idToDelete string
	fmt.Print("Enter the ID to delete: ")
	fmt.Scanln(&idToDelete)

	found := false
	for i, record := range records {
		if i == 0 || record[0] != idToDelete {
			continue
		}

		records = append(records[:i], records[i+1:]...)
		found = true
		break
	}

	if err := file.Truncate(0); err != nil {
		fmt.Println("Error clearing the file:", err)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("Error resetting the file offset:", err)
		return
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		fmt.Println("Error writing the records:", err)
		return
	}
	writer.Flush()

	if found {
		fmt.Println("Record deleted successfully!")
		renumberRecords(selectedFile)
	} else {
		fmt.Println("No record found with the specified ID:", idToDelete)
	}
}

func getNextID(selectedFile string) string {
	file, err := os.Open(selectedFile)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return "1"
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the records:", err)
		return "1"
	}

	maxID := 0
	for i, record := range records {
		if i == 0 {
			continue
		}

		recordID := record[0]
		if len(recordID) > 0 {
			id, _ := strconv.Atoi(recordID)
			if id > maxID {
				maxID = id
			}
		}
	}

	maxID++
	return strconv.Itoa(maxID)
}

func renumberRecords(selectedFile string) {
	file, err := os.OpenFile(selectedFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the records:", err)
		return
	}

	for i := 1; i < len(records); i++ {
		records[i][0] = strconv.Itoa(i)
	}

	if err := file.Truncate(0); err != nil {
		fmt.Println("Error clearing the file:", err)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("Error resetting the file offset:", err)
		return
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		fmt.Println("Error writing the records:", err)
		return
	}
	writer.Flush()

	fmt.Println("Records renumbered successfully!")
}

func listFiles(extension string) ([]string, error) {
	var files []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func readHeaders(selectedFile string) ([]string, error) {
	file, err := os.Open(selectedFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return headers, nil
}

func createNewCSV() {
	fmt.Println("Create a New CSV")

	var fileName string
	fmt.Print("Enter the name for the new CSV file (without .csv extension): ")
	fmt.Scanln(&fileName)

	fileName = fileName + ".csv"

	var headers []string
	fmt.Println("Enter the column headers (one header per line). Enter 'done' when finished.")
	for {
		var header string
		fmt.Print("Header (or 'done'): ")
		fmt.Scanln(&header)

		if header == "done" {
			break
		}

		headers = append(headers, header)
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating the new CSV file:", err)
		return
	}
	defer file.Close()

	headersWithComma := append([]string{""}, headers...)

	writer := csv.NewWriter(file)
	writer.Write(headersWithComma)
	writer.Flush()

	fmt.Println("New CSV file created successfully:", fileName)
}
