package main

import (
	"bufio"
	"flag"
	"fmt"
	CLI_Expenses "github.com/RicliZz/CLI_expenses"
	"os"
	"strconv"
	"strings"
)

const (
	jsonfile = "expenses.json"
)

func main() {
	add := flag.Bool("add", false, "add a new expense")
	list := flag.Bool("list", false, "list all expenses")
	rm := flag.Bool("rm", false, "delete a expense")
	category := flag.String("category", "", "category")
	del := flag.Int("del", 0, "delete a expense")
	balance := flag.Bool("balance", false, "balance a expense")
	balance_add := flag.Int("balance_add", 0, "add a new expense")
	flag.Parse()

	expenses := &CLI_Expenses.Expenses{}

	if err := expenses.Open(jsonfile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		fmt.Println("Please provide a new expense in format:\n1)NAME AND PRICE\n2)NAME, PRICE AND COUNT")
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
	case *category != "":
		r := bufio.NewReader(os.Stdin)
		fmt.Println("Select the ID of the entry you want to change")
		idStr, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		idStr = strings.TrimSpace(idStr)
		intId, _ := strconv.Atoi(idStr)
		err = expenses.AddCategory(intId, category)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = expenses.Save(jsonfile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case *del > 0:
		err := expenses.Del(*del)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = expenses.Save(jsonfile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case *balance:
		line, err := CLI_Expenses.WriteBalance()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(line)
	case *balance_add > 0:
		err := CLI_Expenses.AddBalance(*balance_add)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
