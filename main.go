package main

import (
	"bytes"
	"fmt"
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

func renderJson(command string, from int32, until int32) string {
	src := []byte(command)

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
	fmt.Println("incoming request:")
	fmt.Println("command:", command)
	fmt.Println("from:   ", from)
	fmt.Println("until:  ", until)
	cmd := ""
	skip := false
	for i, t := range tokens {
		if skip {
			skip = false
			continue
		}
		switch t.tok {
		case token.IDENT:
			// a function is starting
			if tokens[i+1].tok == token.LPAREN {
				cmd += Functions[t.lit]
				skip = true // skip the next LPAREN, we already included one
				// this is the beginning of a target string
			} else if tokens[i+1].tok == token.PERIOD && tokens[i-1].tok != token.PERIOD {
				cmd += "ReadMetric(\"" + t.lit
				// this is the end of a target string
			} else if tokens[i-1].tok == token.PERIOD && tokens[i+1].tok != token.PERIOD {
				cmd += t.lit + "\", from, until)"
			} else {
				cmd += t.lit
			}
		case token.LPAREN:
			cmd += "("
		case token.RPAREN:
			cmd += ")"
		case token.PERIOD:
			cmd += "."
		case token.COMMA:
			cmd += ", "
		case token.INT:
			cmd += t.lit
		}
	}
	type Params struct {
		From  int32
		Until int32
		Cmd   string
	}
	params := Params{from, until, cmd}
	fmt.Println("writing to template", params)
	t.Execute(fo, params)
	// TODO: timeout, display errors, etc
	fmt.Printf("executing: go run %s functions.go data.go\n", fname)
	cmd_exec := exec.Command("go", "run", fname, "functions.go", "data.go")
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
	const lenPath = len("/render/")
	from := int32(0)
	until := int32(360)
	command := r.URL.Path[lenPath:]
	r.ParseForm()
	fmt.Println("FORM", r.Form)
	map_from := r.Form["from"]
	if len(map_from) > 0 {
		from_i64, err := strconv.ParseInt(map_from[0], 10, 32)
		if err != nil {
			fmt.Fprintf(w, "Error: invalid 'from' spec: "+map_from[0])
			return
		}
		from = int32(from_i64)
	}
	map_until := r.Form["until"]
	if len(map_until) > 0 {
		until_i64, err := strconv.ParseInt(map_until[0], 10, 32)
		if err != nil {
			fmt.Fprintf(w, "Error: invalid 'until' spec: "+map_until[0])
			return
		}
		until = int32(until_i64)
	}
	fmt.Fprintf(w, renderJson(command, from, until))
}

func main() {
	http.HandleFunc("/render/", renderHandler)
	http.ListenAndServe(":8080", nil)
}
