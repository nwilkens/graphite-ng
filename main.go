package main

import (
	"bytes"
	"fmt"
	"github.com/graphite-ng/graphite-ng/functions"
	"go/scanner"
	"go/token"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"text/template"
)

type Token struct {
	pos token.Position
	tok token.Token
	lit string
}

func NewToken(pos token.Position, tok token.Token, lit string) *Token {
	return &Token{pos: pos, tok: tok, lit: lit}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s\t%s\t%q", t.pos, t.tok, t.lit)
}

func generateCommand(target string) string {
	src := []byte(target)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	tokens := make([]Token, 50)
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		tokens = append(tokens, *NewToken(fset.Position(pos), tok, lit))
	}
	cmd := ""
	for i, t := range tokens {
		switch t.tok {
		case token.IDENT:
			// a function is starting
			if tokens[i+1].tok == token.LPAREN {
				cmd += "functions." + functions.Functions[t.lit]
				// this is the beginning of a target string
			} else if tokens[i+1].tok == token.PERIOD && tokens[i-1].tok != token.PERIOD {
				cmd += "ReadMetric(\"" + t.lit
				// this is the end of a target string
			} else if tokens[i-1].tok == token.PERIOD && tokens[i+1].tok != token.PERIOD {
				cmd += t.lit + "\")"
			} else {
				cmd += t.lit
			}
		case token.LPAREN:
			cmd += "(\n"
		case token.RPAREN:
			cmd += ")"
		case token.PERIOD:
			cmd += "."
		case token.COMMA:
			cmd += ",\n"
		case token.INT:
			cmd += t.lit
		case token.FLOAT:
			cmd += t.lit
		}
	}
	return cmd
}
func renderJson(targets_list []string, from int32, until int32) string {
	type Target struct {
		Query string
		Cmd   string
	}
	type Params struct {
		From    int32
		Until   int32
		Targets []Target
	}
	targets := make([]Target, 0)
	for _, target := range targets_list {
		targets = append(targets, Target{target, generateCommand(target)})
	}
	params := Params{from, until, targets}
	t, err := template.ParseFiles("executor.go.tpl")
	if err != nil {
		panic(err)
	}
	fname := fmt.Sprintf("executor-%d.go", rand.Int())
	fo, err := os.Create(fname)
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("writing to template", params)
	t.Execute(fo, params)
	// TODO: timeout, display errors, etc
	fmt.Printf("executing: go run %s data.go\n", fname)
	cmd_exec := exec.Command("go", "run", fname, "data.go")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd_exec.Stdout = &stdout
	cmd_exec.Stderr = &stderr
	err = cmd_exec.Run()
	if err != nil {
		fmt.Printf("error:", err)
	}
	if stderr.Len() > 0 {
		fmt.Printf("sterr:", stderr.String())
		return stderr.String()
	}
	return stdout.String()
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	from := int32(0)
	until := int32(360)
	r.ParseForm()
	from_list := r.Form["from"]
	if len(from_list) > 0 {
		from_i64, err := strconv.ParseInt(from_list[0], 10, 32)
		if err != nil {
			fmt.Fprintf(w, "Error: invalid 'from' spec: "+from_list[0])
			return
		}
		from = int32(from_i64)
	}
	until_list := r.Form["until"]
	if len(until_list) > 0 {
		until_i64, err := strconv.ParseInt(until_list[0], 10, 32)
		if err != nil {
			fmt.Fprintf(w, "Error: invalid 'until' spec: "+until_list[0])
			return
		}
		until = int32(until_i64)
	}
	targets_list := r.Form["target"]
	for _, target := range targets_list {
		if target == "" {
			fmt.Fprintf(w, "invalid request: one or more empty targets")
			return
		}
	}
	if len(targets_list) < 1 {
		fmt.Fprintf(w, "invalid request: no targets requested")
	} else {
		fmt.Fprintf(w, renderJson(targets_list, from, until))
	}
}

func main() {
	fmt.Println("registered functions:")
	for k, v := range functions.Functions {
		fmt.Printf("%-20s -> %s\n", k, v)
	}
	http.HandleFunc("/render/", renderHandler)
	http.ListenAndServe(":8080", nil)
}
