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
	list := flag.Bool("list", false, "list all expenses")
	rm := flag.Bool("rm", false, "delete a expense")
	flag.Parse()

	expenses := &CLI_Expenses.Expenses{}

	if err := expenses.Open(jsonfile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		line, err := CLI_Expenses.InputName(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = expenses.Add(line)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = expenses.Save(jsonfile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case *list:
		expenses.List()
	case *rm:
		CLI_Expenses.FullClear()
	}
}
