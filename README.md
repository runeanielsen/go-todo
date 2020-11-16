# Go-todo

## Usage

Simple CLI todo-app

### Add

#### From argument

``` sh
./todo -a "My new todo one"
```

#### From STDIN

``` sh
echo "My new todo one" | ./todo -a
```

Multi-line example resulting in multiple todos

``` sh
echo "My new todo one\nMy new todo two" | ./todo -a
```

### List

``` sh
./todo -l
```

#### Verbose

``` sh
./todo -l -v
```

#### Show only todos that are not completed

``` sh
./todo -l -hc
```

### Complete

```sh
# example ./todo -c 1
./todo -c <task-number>
```

### Delete

``` sh
# example ./todo -d 1
./todo -d <task-number>
```

## Dev

* Requires [Go-task](https://taskfile.dev/)
* Requires Go version 1.15

### Build

``` sh
task build
```

### Test

``` sh
task test
```
