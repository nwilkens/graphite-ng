package main

import (
	"fmt"
	"go/scanner"
	"go/token"
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

func main() {
	test1 := "sum(stats.web1.bytes_received,scale(stats.web2.bytes_received,5))&from=60&until=300"
	src := []byte(test1)

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
	for _, t := range tokens {
		fmt.Println(t.String())
	}

	from := 60
	until := 300
	out := Functions["sum"](from, until, readMetric("stats.web1.bytes_received"), Functions["scale"](from, until, readMetric("stats.web2.bytes_received"), 5))
	for {
		d := <-out
		fmt.Println(d)
		if d.ts >= until {
			break
		}
	}

}
