# db
### CSV-based DB
The CSV-based DB tool is a database based on CSV, which means it stores and manipulates data in CSV (Comma-Separated Values) format.

The goal of CSV-based DB is to provide a simple and lightweight solution for data management using CSV files. This approach can be useful in scenarios where a full database is not required or when working with CSV files is preferred for ease of portability and integration with other tools.

The tool allows users to perform basic database operations such as data insertion, querying, updating, and deletion in CSV files.
Overall, CSV-based DB offers a straightforward and efficient alternative for data management using CSV files, enabling users to store and manipulate data easily and flexibly.

## Usage

### Command-line
```
Usage of db:
  -add
        Add entry to CSV
  -column string
        Column name (for get command)
  -create
        Create CSV file
  -edit
        Edit entry
  -f string
        Select CSV file by its path
  -get
        Get entry or value
  -id int
        ID (for get command)
  -list
        List CSV files or entries in CSV
  -n int
        Select CSV file by its number
  -search
        Search entries in one or more CSV
```

### Shell
The shell allows users to perform various operations on CSV files, including listing CSV files, creating new CSV files, adding records, editing records, listing records, searching records, and deleting records. Here's an overview of its functionality:

  1.  **List CSV Files**: Lists all the CSV files in the current directory and its subdirectories.

  2.  **Select CSV File**: Allows the user to select a specific CSV file for further operations.

  3.  **Create CSV File**: Helps in creating a new CSV file by specifying column headers.

  4.  **Search Records**: Searches for records containing a specific search term in one or more CSV files.

  5.  **Exit**: Exits the program.

Once a CSV file is selected, the user can perform the following operations on it:

  1.  **Add Record**: Adds a new record to the selected CSV file, prompting the user to input values for each column.

  2.  **List Records**: Lists all the records in the selected CSV file, displaying them one record per line.

  3.  **List Records as Table**: Lists all the records in a tabular format, aligning columns for better readability.

  4.  **Search Record**: Searches for a specific record by providing a search term.

  5.  **Edit Record**: Allows the user to edit an existing record in the selected CSV file.

  6.  **Delete Record**: Deletes a specific record from the selected CSV file.

The tool also handles some features like automatically assigning IDs to records, handling date columns, and managing CSV files.

Overall, it provides a command-line interface for managing CSV-based databases with various CRUD (Create, Read, Update, Delete) operations on records within those CSV files.
