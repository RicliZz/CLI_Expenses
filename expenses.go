package CLI_Expenses

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"io"
	"os"
	"strconv"
	"strings"
)

const Balancefile = "balance.txt"

type item struct {
	Name     string
	Count    int
	Price    int
	Total    int
	Category string
}

type Expenses []item

func (e *Expenses) Open(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = CreateBalance(Balancefile)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	err = json.Unmarshal(file, e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Expenses) Save(filename string) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (e *Expenses) Add(purchase string) error {
	var newBalance int
	PurchaseSlice := strings.Split(purchase, " ")
	file, err := os.OpenFile(Balancefile, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	balance := scanner.Text()
	file.Truncate(0)
	file.Seek(0, 0)
	balanceInt, err := strconv.Atoi(balance)
	if err != nil {
		return err
	}
	price, err := strconv.Atoi(PurchaseSlice[1])
	if err != nil {
		return err
	}
	if len(PurchaseSlice) > 2 {
		count, _ := strconv.Atoi(PurchaseSlice[2])
		newBalance = balanceInt - price*count
	} else {
		newBalance = balanceInt - price
	}
	newBalanceString := strconv.Itoa(newBalance)
	_, err = file.WriteString(newBalanceString)
	if err != nil {
		return err
	}
	switch len(PurchaseSlice) {
	case 0:
		return errors.New("New entry is empty")

	case 1:
		return errors.New("Insufficient data")

	case 2:
		newItem := item{
			Name:     PurchaseSlice[0],
			Count:    1,
			Price:    price,
			Total:    price,
			Category: "",
		}
		*e = append(*e, newItem)

	case 3:
		count, _ := strconv.Atoi(PurchaseSlice[2])
		newItem := item{
			Name:     PurchaseSlice[0],
			Count:    count,
			Price:    price,
			Total:    count * price,
			Category: "",
		}
		*e = append(*e, newItem)
	}
	return nil
}

func (e *Expenses) List() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Price"},
			{Align: simpletable.AlignCenter, Text: "Count"},
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Total"},
		},
	}
	subtotal := 0
	for idx, row := range *e {
		r := []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx+1)},
			{Text: row.Name},
			{Text: strconv.Itoa(row.Price)},
			{Text: strconv.Itoa(row.Count)},
			{Text: row.Category},
			{Text: strconv.Itoa(row.Total)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		subtotal += row.Total
	}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: "Subtotal"},
			{Text: strconv.Itoa(subtotal)},
		},
	}
	table.SetStyle(simpletable.StyleUnicode)
	fmt.Println(table.String())
}

func InputName(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if len(scanner.Text()) == 0 {
		return "", errors.New("Empty input")
	}

	return scanner.Text(), nil
}

func FullClear() {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("All data will be DELETED. You really want this?[Y/n] ")
	input, _ := r.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" {
		err := os.Remove("expenses.json")
		if err != nil {
			fmt.Println("Error deleting file:", err)
		} else {
			fmt.Println("File deleted successfully.")
		}
	} else {
		fmt.Println("Operation cancelled.")
	}
}

func (e *Expenses) AddCategory(id int, category *string) error {
	if id < 0 || id > len(*e) {
		return errors.New("No such item")
	}
	(*e)[id-1].Category = *category
	return nil
}

func (e *Expenses) Del(id int) error {
	if id < 0 || id > len(*e) {
		return errors.New("No such item")
	}
	*e = append((*e)[:id-1], (*e)[id:]...)
	return nil
}

func CreateBalance(filepath string) error {
	var number string
	fmt.Println("Please enter your current balance\n")
	_, err := fmt.Scan(&number)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return errors.New("Failed to create file with balance")
	}
	defer file.Close()

	_, err = file.WriteString(number)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created file with balance!")
	return nil
}

func WriteBalance() (string, error) {
	var line string
	file, err := os.Open(Balancefile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line = scanner.Text()
	}
	return line, nil
}

func AddBalance(newBalance int) error {
	var line string
	file, err := os.OpenFile(Balancefile, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
	}
	IntLine, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	IntLine += newBalance
	file.Seek(0, 0)
	StringLine := strconv.Itoa(IntLine)
	_, err = file.WriteString(StringLine)
	if err != nil {
		return err
	}
	return nil
}
