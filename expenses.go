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
	purchase_slice := strings.Split(purchase, " ")
	switch len(purchase_slice) {
	case 0:
		return errors.New("New entry is empty")
	case 2:
		fmt.Println("YES")
		price, err := strconv.Atoi(purchase_slice[1])
		if err != nil {
			return err
		}
		newItem := item{
			Name:     purchase_slice[0],
			Count:    1,
			Price:    price,
			Total:    price,
			Category: "",
		}
		*e = append(*e, newItem)
	case 3:
		price, err := strconv.Atoi(purchase_slice[1])
		count, err := strconv.Atoi(purchase_slice[2])
		if err != nil {
			return err
		}
		newItem := item{
			Name:     purchase_slice[0],
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
