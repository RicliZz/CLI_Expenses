package CLI_Expenses

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

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
	PurchaseSlice := strings.Split(purchase, " ")
	switch len(PurchaseSlice) {
	case 0:
		return errors.New("New entry is empty")
	case 2:
		price, err := strconv.Atoi(PurchaseSlice[1])
		if err != nil {
			return err
		}
		newItem := item{
			Name:     PurchaseSlice[0],
			Count:    1,
			Price:    price,
			Total:    price,
			Category: "",
		}
		*e = append(*e, newItem)
	case 3:
		price, err := strconv.Atoi(PurchaseSlice[1])
		count, err := strconv.Atoi(PurchaseSlice[2])
		if err != nil {
			return err
		}
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
	for idx, item := range *e {
		fmt.Printf("%d - %s\n", idx+1, item.Name)
	}
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
