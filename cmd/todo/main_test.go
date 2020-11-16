package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Build tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests....")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	task := "test task number 1"
	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListsTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("ListsTasksVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		matchPattern := "  \\d{1,}: [a-zA-Z 0-9]*, \\d{1,} *[a-zA-z]{0,3} \\d{4} \\d{1,2}:\\d{1,2}:\\d{1,2}\\\n"
		regExp := fmt.Sprintf("^%s%s$", matchPattern, matchPattern)

		matches, err := regexp.MatchString(regExp, string(out))
		if err != nil {
			t.Fatal(err)
		}

		if !matches {
			t.Errorf("Pattern %q, did not match output %q\n", matchPattern, string(out))
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListsTasksVerboseHideCompleted", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-verbose", "-hide-completed")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		matchPattern := "  \\d{1,}: [a-zA-Z 0-9]*, \\d{1,} *[a-zA-z]{0,3} \\d{4} \\d{1,2}:\\d{1,2}:\\d{1,2}\\\n"
		regExp := fmt.Sprintf("^%s$", matchPattern)

		matches, err := regexp.MatchString(regExp, string(out))
		if err != nil {
			t.Fatal(err)
		}

		if !matches {
			t.Errorf("Pattern %q, did not match output %q\n", matchPattern, string(out))
		}
	})

	t.Run("DeleteTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", "1")

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
}
