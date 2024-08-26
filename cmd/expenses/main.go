package main

import (
	"flag"
	"fmt"
	CLI_Expenses "github.com/RicliZz/CLI_expenses"
	"os"
)

const (
	jsonfile = "expenses.json"
)

func main() {
	add := flag.Bool("add", false, "add a new expense")
	flag.Parse()

	expenses := &CLI_Expenses.Expenses{}

	if err := expenses.Open(jsonfile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		fmt.Printf("Please enter only:\n NAME AND PRICE\n or \n NAME, PRICE and COUNT(if COUNT is not 1)\n")
		err := expenses.Add("Bluuuu", "200")
		if err != nil {
			fmt.Println(err.Error())
		}
		err = expenses.Save(jsonfile)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
