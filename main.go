package main

import (
	"fmt"
	"go/scanner"
    "math/rand"
	"go/token"
    "bytes"
	"net/http"
	"os"
    "os/exec"
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
	fo, err := os.Create("executor.go")
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
    fmt.Println("original command:")
    fmt.Println(command)
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
    fmt.Printf(params.Cmd)
	t.Execute(fo, params)
    # go run executor.go functions.go data.go
	return out
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	const lenPath = len("/render/")
	command := r.URL.Path[lenPath:]
	fmt.Fprintf(w, renderJson(command, 1, 1))
}

func main() {
	http.HandleFunc("/render/", renderHandler)

	json := renderJson("sum(stats.web1.bytes_received,scale(stats.web2.bytes_received,5))", 60, 300)
	fmt.Print(json)
	http.ListenAndServe(":8080", nil)
}
