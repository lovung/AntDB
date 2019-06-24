package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/btree"

	"github.com/lovung/AntDB/pkg/golib"
)

const (
	all     = 0
	exactly = 1
	gte     = 2
	lte     = 3
	glte    = 4
)

var (
	memMap  = make(map[btree.Int]int)
	slice   []int
	fi      = flag.String("i", "resources/20.txt", "file path to read from")
	command = flag.String("c", "", "command")
	degree  = flag.Int("degree", 8, "degree of btree")
	verbose = flag.Bool("v", false, "verbose")
	tr      *btree.BTree
	count   int64
)

func main() {
	flag.Parse()
	var t, vals interface{}
	// var stats runtime.MemStats
	for i := 0; i < 10; i++ {
		runtime.GC()
	}

	if *verbose {
		// fmt.Println("-------- BEFORE ----------")
		// runtime.ReadMemStats(&stats)
		// fmt.Printf("%+v\n", stats)
	}
	golib.ReadWholeFile(*fi, &slice)
	vals = slice
	start := time.Now()
	tr = btree.New(*degree)
	for _, v := range slice {
		i := btree.Int(v)
		if !tr.Has(i) {
			tr.ReplaceOrInsert(i)
			memMap[i] = 1
		} else {
			memMap[i]++
		}
	}
	t = tr // keep it around
	if *verbose {
		fmt.Printf("%v inserts in %v\n", len(slice), time.Since(start))
		// fmt.Println("-------- AFTER ----------")
		// runtime.ReadMemStats(&stats)
		// fmt.Printf("%+v\n", stats)
	}
	for i := 0; i < 10; i++ {
		runtime.GC()
	}

	if *verbose {
		// fmt.Println("-------- AFTER GC ----------")
		// runtime.ReadMemStats(&stats)
		// fmt.Printf("%+v\n", stats)
	}
	if t == vals {
		fmt.Println("to make sure vals and tree aren't GC'd")
	}
	if *command != "" {
		fmt.Printf("Process: %s\n", *command)
		processCommand(*command)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n-------ANT DB--------")
	fmt.Println("---------------------")
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, " ", "", -1)
		text = strings.Replace(text, "\t", "", -1)
		text = strings.ToLower(text)
		t := time.Now()
		processCommand(text)
		if len(text) > 0 {
			fmt.Printf("Time to process %v items: %v\n", count, time.Since(t))
		}
		count = 0
	}
}

func processCommand(cmd string) {
	if len(cmd) < 1 {
		return
	}
	if strings.Compare(cmd, "h") == 0 {
		fmt.Println("Example command: \n\t<= 1000\n\t>=20000\n\t=50\n\t>=2000 and <= 3000")
	} else if strings.Compare(cmd, "a") == 0 {
		handler(all)
	} else if cmd[0] == '=' {
		input, err := strconv.Atoi(cmd[1:])
		if err != nil {
			fmt.Println("Wrong input")
			return
		}
		handler(exactly, input)
	} else {
		if strings.Contains(cmd, "and") {
			s := strings.Split(cmd[2:], "and")
			if len(s) != 2 {
				fmt.Println("And need exactly 2 input values")
				return
			}
			fmt.Printf("i1 = %s, i2 = %s\n", s[0], s[1])
			input1, err := strconv.Atoi(s[0])
			if err != nil {
				fmt.Println("Wrong input 1")
				return
			}
			input2, err := strconv.Atoi(s[1][2:])
			if err != nil {
				fmt.Println("Wrong input 2")
				return
			}
			handler(glte, input1, input2)
			return
		}
		if strings.Contains(cmd, ">=") {
			input, err := strconv.Atoi(cmd[2:])
			if err != nil {
				fmt.Println("Wrong input")
				return
			}
			handler(gte, input)
		}
		if strings.Contains(cmd, "<=") {
			input, err := strconv.Atoi(cmd[2:])
			if err != nil {
				fmt.Println("Wrong input")
				return
			}
			handler(lte, input)
		}
	}
}

func printItem(v btree.Item) bool {
	var val = v.(btree.Int)
	for i := 0; i < memMap[val]; i++ {
		fmt.Printf("%d ", val)
		count++
	}
	return true
}

func handler(typeOfCommand int, input ...int) {
	var i, i2 btree.Int
	if len(input) > 0 {
		i = btree.Int(input[0])
		if *verbose {
			fmt.Printf("Input 1: %d\n", i)
		}
	}
	if len(input) == 2 {
		i2 = btree.Int(input[1])
		fmt.Printf("Input 2: %d\n", i2)
	}

	switch typeOfCommand {
	case all:
		fmt.Print("[ ")
		tr.Ascend(printItem)
	case exactly:
		fmt.Print("[ ")
		if tr.Has(i) {
			printItem(i)
		}
	case gte:
		fmt.Print("[ ")
		tr.AscendGreaterOrEqual(i, printItem)
	case lte:
		fmt.Print("[ ")
		tr.AscendLessThan(i, printItem)
		if tr.Has(i) {
			printItem(i)
		}
	case glte:
		fmt.Print("[ ")
		tr.AscendRange(i, i2, printItem)
		if tr.Has(i) {
			printItem(i)
		}
	default:
		fmt.Println("Wrong case")
	}
	fmt.Print("]\n")
}
