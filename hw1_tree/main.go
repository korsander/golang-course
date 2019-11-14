package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var rootInfo, err = os.Stat(path)
	if err != nil {
		return err
	}

	if !rootInfo.IsDir() {
		return fmt.Errorf("path must be a dir")
	}

	return printDir(out, []bool{}, path, printFiles)
}

func printDir(out io.Writer, lasts []bool, root string, printFiles bool) error {
	var nodes []node
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, info := range files {
		if !printFiles && !info.IsDir() {
			continue
		}

		name := info.Name()
		if !info.IsDir() {
			var size = strconv.FormatInt(info.Size(), 10) + "b"
			if info.Size() == 0 {
				size = "empty"
			}
			name = fmt.Sprintf("%s (%s)", name, size)
		}

		nodes = append(nodes, node{path.Join(root, info.Name()), name, info.IsDir()})
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].name < nodes[j].name
	})

	size := len(nodes)
	for i, node := range nodes {
		var bb bytes.Buffer
		isLast := i == size-1

		for _, v := range lasts {
			if !v {
				bb.WriteString("│")
			}
			bb.WriteString("\t")
		}

		if isLast {
			bb.WriteString("└")
		} else {
			bb.WriteString("├")
		}
		bb.WriteString("───")
		bb.WriteString(node.name)
		bb.WriteString("\n")

		out.Write(bb.Bytes())
		if node.isDir {
			printDir(out, append(lasts, isLast), node.path, printFiles)
		}
	}

	return nil
}

type node struct {
	path  string
	name  string
	isDir bool
}
