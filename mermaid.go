//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Jobs map[string]Job `yaml:"jobs"`
}

type Job struct {
	Needs interface{} `yaml:"needs,omitempty"`
	If    string      `yaml:"if,omitempty"`
}

type Dict struct {
	Nodes []Node
}

type Node struct {
	If    string
	Needs []string
	Ids   []string
}

func (n Node) dumpMermaidNode() string {
	if len(n.Needs) == 0 {
		return fmt.Sprintf("    %s(%s)", strings.Join(n.Ids, ""), strings.Join(n.Ids, "<br>"))
	}
	l := make([]string, len(n.Needs))
	for _, v := range n.Needs {
		l = append(l, fmt.Sprintf("    %s(%s) --> %s(%s)", v, v, strings.Join(n.Ids, ""), strings.Join(n.Ids, "<br>")))
	}
	return strings.Join(l, "\n")
}

func (n *Node) convNeedstoSlice(needs interface{}) {
	switch t := needs.(type) {
	case nil:
		return
	case []interface{}:
		l := reflect.ValueOf(needs)
		for i := 0; i < l.Len(); i++ {
			s := fmt.Sprintf("%s", l.Index(i))
			n.Needs = append(n.Needs, s)
		}
	case string:
		s := fmt.Sprintf("%s", needs)
		n.Needs = append(n.Needs, s)
	default:
		log.Fatalf("%T is not string or slice", t)
	}
}

func main() {
	b, err := os.ReadFile("gha.yml")
	if err != nil {
		log.Fatal(err)
	}
	config := Workflow{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	keys := sortedKeys(config)
	fmt.Println(dump(config, keys))
}

func dump(w Workflow, keys []string) string {
	first := "```mermaid"
	second := "flowchart LR"
	last := "```"
	lines := make([]string, 0, len(w.Jobs)+3)
	lines = append(lines, first)
	lines = append(lines, second)
	for _, v := range keys {
		switch t := w.Jobs[v].Needs.(type) {
		case nil:
			lines = append(lines, nodeLink(v, ""))
		case []interface{}:
			l := reflect.ValueOf(w.Jobs[v].Needs)
			for i := 0; i < l.Len(); i++ {
				s := fmt.Sprintf("%s", l.Index(i))
				lines = append(lines, nodeLink(s, v))
			}
		case string:
			s := fmt.Sprintf("%s", w.Jobs[v].Needs)
			lines = append(lines, nodeLink(s, v))
		default:
			log.Fatalf("%T is not string or slice", t)
		}
	}
	lines = append(lines, last)
	return strings.Join(lines, "\n")
}

func sortedKeys(w Workflow) []string {
	keys := make([]string, 0, len(w.Jobs))
	for k := range w.Jobs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func nodeLink(from string, to string) string {
	// TODO: needsとif文が同じだったらsubgraphでまとめる
	if to == "" {
		return fmt.Sprintf("    %s(%s)", from, from)
	}
	return fmt.Sprintf("    %s(%s) --> %s(%s)", from, from, to, to)
}
