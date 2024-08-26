package CLI_Expenses

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
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

func (e *Expenses) Add(purchase ...string) error {
	switch len(purchase) {
	case 0:
		return errors.New("New entry is empty")
	case 2:
		price, err := strconv.Atoi(purchase[1])
		if err != nil {
			return err
		}
		newItem := item{
			Name:     purchase[0],
			Count:    1,
			Price:    price,
			Total:    price,
			Category: "",
		}
		*e = append(*e, newItem)
	case 3:
		price, err := strconv.Atoi(purchase[1])
		count, err := strconv.Atoi(purchase[2])
		if err != nil {
			return err
		}
		newItem := item{
			Name:     purchase[0],
			Count:    count,
			Price:    price,
			Total:    count * price,
			Category: "",
		}
		*e = append(*e, newItem)
	}
	return nil
}
