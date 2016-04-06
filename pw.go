package pw

import (
	"fmt"
	"unicode"
)

const (
	Version = "Version: 0.9.5"
)

type ParseWords struct {
	buf string // Current string
	pos int    // where we are in string
	st  int    // current state for DFA
	db  bool   // Debuging Flat
	qf  string // "C"		== '" with \
	// "SQL"	== "" or ''				// xyzzy - TBD
	keep_quote     bool // Keep ' and " in output
	keep_backslash bool // Keep \\ in output
}

func NewParseWords() (pw *ParseWords) {
	return &ParseWords{qf: "C", db: false, keep_quote: false, keep_backslash: false}
}

func (this *ParseWords) SetOptions(qf string, kq, kb bool) {
	this.qf = qf
	this.keep_quote = kq
	this.keep_backslash = kb
}

func (this *ParseWords) SetDebug(b bool) {
	this.db = b
}

func (this *ParseWords) AppendLine(s string) {
	this.buf += s
}

func (this *ParseWords) SetLine(s string) {
	this.buf = s
}

// xyzzy - TBD
/*

From: http://blog.golang.org/strings

	const nihongo = "日本語"
    for i, w := 0, 0; i < len(nihongo); i += w {
        runeValue, width := utf8.DecodeRuneInString(nihongo[i:])
        fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
        w = width
    }
*/

func (this *ParseWords) GetWords() []string {
	i := this.pos
	l := len(this.buf)
	rv := make([]string, 0, 10)
	cs := ""
	c := ""
	wf := false
	for i < l {
		c = this.buf[i : i+1]
		if this.db {
			fmt.Printf("top st=%d c->%s<-\n", this.st, c)
		}
		switch this.st {
		default:
			fmt.Printf("Error(): Invalid state(%d)\n", this.st)
			i++
		case 0: // scan across to blank
			if c == "\"" {
				this.st = 1
				cs = ""
				if this.keep_quote {
					cs += c
				}
				wf = true
				i++
			} else if c == "'" {
				this.st = 11
				cs = ""
				if this.keep_quote {
					cs += c
				}
				wf = true
				i++
			} else if unicode.IsSpace(rune(c[0])) {
				this.st = 2
				if wf {
					rv = append(rv, cs)
					cs = ""
					wf = false
				}
				i++
			} else {
				wf = true
				cs += c
				i++
			}
		case 1: // Start of "
			if c == "\\" {
				this.st = 3
				if this.keep_backslash {
					cs += c
				}
				i++
			} else if c == "\"" {
				this.st = 0
				if this.keep_quote {
					cs += c
				}
				if wf {
					rv = append(rv, cs)
					cs = ""
					wf = false
				}
				i++
			} else {
				cs += c
				i++
			}
		case 11: // Start of "
			if c == "\\" {
				this.st = 13
				if this.keep_backslash {
					cs += c
				}
				i++
			} else if c == "'" {
				this.st = 0
				if this.keep_quote {
					cs += c
				}
				if wf {
					rv = append(rv, cs)
					cs = ""
					wf = false
				}
				i++
			} else {
				cs += c
				i++
			}
		case 2: // Found blank
			// Scan across blanks until non-blank
			if unicode.IsSpace(rune(c[0])) {
				i++
			} else {
				wf = true
				if c == "\"" {
					this.st = 4
					wf = true
					cs = ""
					if this.keep_quote {
						cs += c
					}
				} else if c == "'" {
					this.st = 14
					wf = true
					cs = ""
					if this.keep_quote {
						cs += c
					}
				} else {
					this.st = 0
					cs = c
				}
				i++
			}
		case 3: // \" processing
			this.st = 1
			cs += c
			i++
		case 13: // \' processing
			this.st = 11
			cs += c
			i++
		case 4: // scan across to blank
			if c == "\"" {
				this.st = 0
				if this.keep_quote {
					cs += c
				}
				rv = append(rv, cs)
				cs = ""
				wf = false
				i++
			} else if c == "\\" {
				this.st = 5
				if this.keep_backslash {
					cs += c
				}
				i++
			} else {
				wf = true
				cs += c
				i++
			}
		case 14: // scan across to blank
			if c == "'" {
				this.st = 0
				if this.keep_quote {
					cs += c
				}
				rv = append(rv, cs)
				cs = ""
				wf = false
				i++
			} else if c == "\\" {
				this.st = 15
				if this.keep_backslash {
					cs += c
				}
				i++
			} else {
				wf = true
				cs += c
				i++
			}
		case 5: // \" processing
			this.st = 4
			cs += c
			i++
		case 15: // \' processing
			this.st = 14
			cs += c
			i++
		}
	}
	if wf {
		rv = append(rv, cs)
	}
	return rv
}
