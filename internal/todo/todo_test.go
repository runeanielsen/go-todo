package todo_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/runeanielsen/go-todo/internal/todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, l[0].Task)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("New task should not be completed")
	}

	l.Complete(1)

	if !l[0].Done {
		t.Errorf("New task should be completed.")
	}
}

func TestDelete(t *testing.T) {
	l := todo.List{}

	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}

	for _, t := range tasks {
		l.Add(t)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("Expected %q, got %q instead", tasks[0], l[0].Task)
	}

	l.Delete(2)

	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q, got %q instead.", tasks[2], l[1].Task)
	}
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "New Task"
	l1.Add(taskName)

	if l1[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, l1[0].Task)
	}

	tf, err := ioutil.TempFile("", "")

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to list: %s", err)
	}

	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should match %q task.", l1[0].Task, l2[0].Task)
	}
}

func TestVerbose(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}

	expected := fmt.Sprintf("%s%d: %s, %s\n", "  ", 1, l[0].Task, l[0].CreatedAt.Format("02 Jan 2006 15:04:05"))

	out := l.Display(true, false)

	if expected != out {
		t.Errorf("Expected %q, got %q instead\n", expected, out)
	}
}

func TestHideCompletedVerbose(t *testing.T) {
	l := todo.List{}

	taskNameOne := "New Task One"
	taskNameTwo := "New Task Two"
	l.Add(taskNameOne)
	l.Add(taskNameTwo)

	l[1].Done = true

	if l[0].Task != taskNameOne {
		t.Errorf("Expected %q, got %q instead", taskNameOne, l[0].Task)
	}

	if l[1].Task != taskNameTwo {
		t.Errorf("Expected %q, got %q instead", taskNameTwo, l[0].Task)
	}

	expected := fmt.Sprintf("%s%d: %s, %s\n", "  ", 1, l[0].Task, l[0].CreatedAt.Format("02 Jan 2006 15:04:05"))

	out := l.Display(true, true)

	if expected != out {
		t.Errorf("Expected %q, got %q instead\n", expected, out)
	}
}
