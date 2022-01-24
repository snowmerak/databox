package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	intogo := kingpin.New("intogo", "A command line tool to convert file to byte buffer in go.")
	directory := intogo.Arg("dir", "The directory to convert.").Required().String()
	packagename := intogo.Arg("pkg", "The package name.").Required().String()
	variableName := intogo.Arg("var", "The variable name.").Required().String()

	kingpin.MustParse(intogo.Parse(os.Args[1:]))

	if *directory == "" || *packagename == "" || *variableName == "" {
		intogo.Usage(os.Args)
		os.Exit(1)
	}

	dirStats, err := os.Stat(*directory)
	if os.IsNotExist(err) {
		panic("Directory does not exist.")
	}
	if err := os.MkdirAll(*packagename, 0755); err != nil {
		panic(err)
	}

	if !dirStats.IsDir() {
		file, err := os.Create(filepath.Join(*packagename, strings.ToLower(*variableName)+".go"))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		src, err := os.Open(*directory)
		if err != nil {
			panic(err)
		}
		defer src.Close()
		data, err := io.ReadAll(src)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(file, "package "+*packagename+"\n\n")
		fmt.Fprint(file, "var "+*variableName+" = []byte{")
		for i, b := range data {
			fmt.Fprint(file, b)
			if i != len(data)-1 {
				fmt.Fprint(file, ", ")
			}
		}
		fmt.Fprint(file, "}\n")

		fmt.Println("Successfully converted file to byte buffer.")
		return
	}

	fileMap := map[string][]byte{}
	queue := []string{*directory}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		stats, err := os.Stat(cur)
		if err != nil {
			panic(err)
		}
		if stats.IsDir() {
			files, err := ioutil.ReadDir(cur)
			if err != nil {
				panic(err)
			}
			for _, f := range files {
				queue = append(queue, filepath.Join(cur, f.Name()))
			}
			continue
		}
		data, err := os.ReadFile(cur)
		if err != nil {
			panic(err)
		}
		fileMap[cur] = data
	}

	file, err := os.Create(filepath.Join(*packagename, strings.ToLower(*variableName)+".go"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, "package "+*packagename+"\n\n")
	fmt.Fprint(file, "var "+*variableName+" = map[string][]byte{\n")
	for k, v := range fileMap {
		path := strings.Join(filepath.SplitList(k), "/")
		fmt.Fprintf(file, "\t`%s`: {", path)
		for i, b := range v {
			fmt.Fprint(file, b)
			if i != len(v)-1 {
				fmt.Fprint(file, ", ")
			}
		}
		fmt.Fprint(file, "},\n")
	}
	fmt.Fprint(file, "}\n")
}
