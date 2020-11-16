package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/runeanielsen/go-todo/internal/todo"
)

var todoFileName = ".todo.json"

func main() {
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	add := flag.Bool("a", false, "Add item to the Todo list")
	list := flag.Bool("l", false, "List all tasks")
	complete := flag.Int("c", 0, "Item to be completed")
	delete := flag.Int("d", 0, "Item to be deleted")
	verbose := flag.Bool("v", false, "Verbose mode")
	hideCompleted := flag.Bool("hc", false, "Hide completed items")

	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(l.Display(*verbose, *hideCompleted))
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		tasks, err := getTask(os.Stdin, flag.Args()...)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, t := range tasks {
			l.Add(t)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) ([]string, error) {
	var t []string

	if len(args) > 0 {
		t = append(t, strings.Join(args, " "))
		return t, nil
	}

	s := bufio.NewScanner(r)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return t, err
		}

		t = append(t, s.Text())
	}

	if len(t) == 0 {
		return t, fmt.Errorf("Task cannot be blank")
	}

	return t, nil
}
