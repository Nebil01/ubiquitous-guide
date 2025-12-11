package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Transaction struct { // hold each Expense details
	ID       int
	Amount   float64
	Category string
	Date     time.Time
	Expense  float64
}

type ExpenseTracker struct { // manage Expense
	transactions []Transaction
	nextID       int
}

func (t Transaction) GetAmount() float64 {
	return t.Amount
}

func (bt *ExpenseTracker) AddExpense(amount float64, category string) {
	newTransaction := Transaction{
		ID:       bt.nextID,
		Amount:   amount,
		Category: category,
		Date:     time.Now(),
	}
	bt.transactions = append(bt.transactions, newTransaction)
	bt.nextID++
}

func (bt *ExpenseTracker) UpdateExpense(id int, amount float64, category string) {
	if amount < 0 {
		fmt.Println("Amount cannot be negative")
	}
	for i, transaction := range bt.transactions {
		if transaction.ID == id {
			bt.transactions[i].Amount = amount
			bt.transactions[i].Category = category
			bt.transactions[i].Date = time.Now()
			break
		}
	}
}

func (bt *ExpenseTracker) DeleteExpense(id int) {
	for i, transaction := range bt.transactions {
		if transaction.ID == id {
			bt.transactions = append(bt.transactions[:i], bt.transactions[i+1:]...)
			break
		}
	}
}

func (bt ExpenseTracker) DisplayExpense() {
	fmt.Printf("%s %5s %-12s %-10s\n", "ID", "Amount", "Category", "Date")
	for _, transaction := range bt.transactions {
		fmt.Printf("%d %5.2f %-12s %-10s\n",
			transaction.ID,
			transaction.Amount,
			transaction.Category,
			transaction.Date.Format("2006-01-02"))
	}
}

func (bt ExpenseTracker) ShowSpecificMonthExpense(date time.Time) []Transaction {
	monthName := time.Month.String(date.Month())
	total := bt.TotalExpense(0)
	fmt.Printf("Total Expenses for %s: %.2f\n", monthName, total)
	var monthlyExpenses []Transaction
	for _, t := range bt.transactions {
		if t.Date.Month() == date.Month() && t.Date.Year() == date.Year() {
			monthlyExpenses = append(monthlyExpenses, t)
		}
	}
	return monthlyExpenses
}

// Get total expense
func (bt ExpenseTracker) TotalExpense(tType float64) float64 {
	var total float64
	for _, transaction := range bt.transactions {
		if transaction.Expense == tType {
			total += transaction.Amount
		}
	}
	return total
}

// save Expenses to csv file
func (bt ExpenseTracker) SaveToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file) // creating new csv file (csv writer buffer its output for efficiency)
	defer writer.Flush()          // Flush: ensures any buffer data is written to the underlying file befor file is closed

	// write csv header
	writer.Write([]string{"ID", "Amount", "Category", "Date"})

	// write Data
	for _, t := range bt.transactions {
		record := []string{
			strconv.Itoa(t.ID),
			fmt.Sprintf("%.2f", t.Amount),
			t.Category,
			t.Date.Format("2006-01-02"),
		}
		writer.Write(record)
	}
	fmt.Printf("Tansactions saved to %s", filename)
	return nil
}

func main() {
	// Instantiation of ExpenseTracker struct
	bt := ExpenseTracker{}
	for {
		fmt.Println("\n--- Personal Expense Tracker ---")
		fmt.Println("1. Add Expense")
		fmt.Println("2. Update Expense")
		fmt.Println("3. Delete Expense")
		fmt.Println("4. Display All Expenses")
		fmt.Println("5. Show Total Expense")
		fmt.Println("6. Show Specific Month Expense")
		fmt.Println("7. Save Expense to csv file")
		fmt.Println("8. Exit")
		fmt.Println("choose an option: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("Enter Amount: ")
			var amount float64
			fmt.Scan(&amount)

			fmt.Println("Enter Category: ")
			var category string
			fmt.Scan(&category)

			bt.AddExpense(amount, category)
			fmt.Println("Expense Added!")
		case 2:
			fmt.Println("Enter Expense ID to update: ")
			var id int
			fmt.Scan(&id)

			fmt.Println("Enter Amount: ")
			var amount float64
			fmt.Scan(&amount)

			fmt.Println("Enter Category: ")
			var category string
			fmt.Scan(&category)

			bt.UpdateExpense(id, amount, category)
			fmt.Println("Expense Updated!")
		case 3:
			fmt.Println("Enter Expense ID to delete: ")
			var id int
			fmt.Scan(&id)
			bt.DeleteExpense(id)
			fmt.Println("Expense Deleted!")
		case 4:
			bt.DisplayExpense()
		case 5:
			fmt.Printf("Total Expenses: %.2f\n", bt.TotalExpense(0))
		case 6:
			fmt.Println("Enter Specific Month (YYYY-MM): ")
			var month string
			fmt.Scan(&month)
			parsedMonth, err := time.Parse("2006-01", month)
			if err != nil {
				fmt.Println("Invalid month format. Please use YYYY-MM format.")
			}
			bt.ShowSpecificMonthExpense(parsedMonth)
		case 7:
			fmt.Println("Enter filename for csv file: ")
			var filename string
			fmt.Scan(&filename)
			if err := bt.SaveToCSV(filename); err != nil {
				fmt.Println("Error saving Expense:", err)
			}
		case 8:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
