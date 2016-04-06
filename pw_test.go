

// (C) Philip Schlump, 2014.

package pw

import (
	"testing"
	"fmt"
	"encoding/json"
)

type TestCase struct {
	In			string
	Out			[]string
	kq			bool
	kb			bool
}

var testCases = []TestCase{
/*  0 */	{ " abc def ",			[]string{ "abc", "def" },			false,	false },
/*  1 */	{ "abc def ",			[]string{ "abc", "def" },			false,	false },
/*  2 */	{ "abc def",			[]string{ "abc", "def" },			false,	false },
/*  3 */	{ "abc  def",			[]string{ "abc", "def" },			false,	false },
/*  4 */	{ "  abc  def",			[]string{ "abc", "def" },			false,	false },
/*  5 */	{ "  abc  def       ",	[]string{ "abc", "def" },			false,	false },
/*  6 */	{ "a b c def",			[]string{ "a", "b", "c", "def" },	false,	false },

/*  7 */	{ ` "a" `,				[]string{ "a", },					false,	false },
/*  8 */	{ ` "a" "e f" `,		[]string{ "a", "e f" },				false,	false },
/*  9 */	{ ` "a""e f" `,			[]string{ "a", "e f" },				false,	false },
/* 10 */	{ ` "a\"b" `,			[]string{ "a\"b", },				false,	false },
/* 11 */	{ ` "a\\b" `,			[]string{ "a\\b", },				false,	false },
/* 12 */	{ `"a\"b" `,			[]string{ "a\"b", },				false,	false },
/* 13 */	{ `"a\"b"`,				[]string{ "a\"b", },				false,	false },
/* 14 */	{ `"a\""`,				[]string{ "a\"", },					false,	false },
/* 15 */	{ `"\""`,				[]string{ "\"", },					false,	false },
/* 16 */	{ `"\\"`,				[]string{ "\\", },					false,	false },

/* 17 */	{ ` 'a' `,				[]string{ "a", },					false,	false },
/* 18 */	{ ` 'a' 'e f' `,		[]string{ "a", "e f" },				false,	false },
/* 19 */	{ ` 'a''e f' `,			[]string{ "a", "e f" },				false,	false },
/* 20 */	{ ` 'a\'b' `,			[]string{ "a'b", },					false,	false },
/* 21 */	{ ` 'a\\b' `,			[]string{ "a\\b", },				false,	false },
/* 22 */	{ `'a\'x' `,			[]string{ "a'x", },					false,	false },
/* 23 */	{ `'a\'b'`,				[]string{ "a'b", },					false,	false },
/* 24 */	{ `'y\''`,				[]string{ "y'", },					false,	false },
/* 25 */	{ `'\''`,				[]string{ "'", },					false,	false },
/* 26 */	{ `'\\'`,				[]string{ "\\", },					false,	false },

/* 27 */	{ ` "a" `,				[]string{ "a", },					false,	false },

/* 28 */	{ ` "a"'e f' `,			[]string{ "a", "e f" },				false,	false },

/* 29 */	{ ` "a" `,				[]string{ `"a"`, },					true,	false },
/* 30 */	{ ` "a" "e f" `,		[]string{ "\"a\"", "\"e f\"" },		true,	false },
/* 31 */	{ ` "a""e f" `,			[]string{ "\"a\"", "\"e f\"" },		true,	false },
/* 32 */	{ ` "a\"b" `,			[]string{ "\"a\"b\"", },			true,	false },
/* 33 */	{ ` "a\\b" `,			[]string{ "\"a\\b\"", },			true,	false },
/* 34 */	{ `"a\"b" `,			[]string{ "\"a\"b\"", },			true,	false },
/* 35 */	{ `"a\"b"`,				[]string{ "\"a\"b\"", },			true,	false },
/* 36 */	{ `"a\""`,				[]string{ "\"a\"\"", },				true,	false },
/* 37 */	{ `"\""`,				[]string{ "\"\"\"", },				true,	false },
/* 38 */	{ `"\\"`,				[]string{ "\"\\\"", },				true,	false },

/* 39 */	{ ` 'a' `,				[]string{ "'a'", },					true,	false },
/* 40 */	{ ` 'a' 'e f' `,		[]string{ "'a'", "'e f'" },			true,	false },
/* 41 */	{ ` 'a''e f' `,			[]string{ "'a'", "'e f'" },			true,	false },
/* 42 */	{ ` 'a\'b' `,			[]string{ "'a'b'", },				true,	false },
/* 43 */	{ ` 'a\\b' `,			[]string{ "'a\\b'", },				true,	false },
/* 44 */	{ `'a\'x' `,			[]string{ "'a'x'", },				true,	false },
/* 45 */	{ `'a\'b'`,				[]string{ "'a'b'", },				true,	false },
/* 46 */	{ `'y\''`,				[]string{ "'y''", },				true,	false },
/* 47 */	{ `'\''`,				[]string{ "'''", },					true,	false },
/* 48 */	{ `'\\'`,				[]string{ "'\\'", },				true,	false },

/* 49 */	{ ` "a"'e f' `,			[]string{ "\"a\"", "'e f'" },		true,	false },

/* 50 */	{ ` 'a\'b' `,			[]string{ "'a\\'b'", },				true,	true },
/* 51 */	{ `'\''`,				[]string{ "'\\''", },				true,	true },
/* 52 */	{ `'\\'`,				[]string{ "'\\\\'", },				true,	true },
/* 53 */	{ `"a\"b" `,			[]string{ "\"a\\\"b\"", },			true,	true },

/* 54 */	{ ` 'a\'b' `,			[]string{ "a\\'b", },				false,	true },
/* 55 */	{ `'\''`,				[]string{ "\\'", },					false,	true },
/* 56 */	{ `'\\'`,				[]string{ "\\\\", },				false,	true },
/* 57 */	{ `"a\"b" `,			[]string{ "a\\\"b", },				false,	true },
}

func arrayEq ( a []string, b []string ) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range(a) {
		if b[i] != v {
			return false
		}
	}
	return true
}

func sVar ( v interface{} ) string {
	s, err := json.Marshal ( v )
	// s, err := json.MarshalIndent ( v, "", "\t" )
	if ( err != nil ) {
		return fmt.Sprintf ( "Error:%s", err )
	} else {
		return string(s)
	}
}


var db bool = false

func TestWf0(t *testing.T) {
	if false {
		fmt.Printf ( "keep compiler happy when we are not using fmt.\n" )
	}
	//pw := NewParseWords ()
	//pw.SetLine ( "  abc def " )
	//x := pw.GetWords()
	//for i, v := range x {
	//	fmt.Printf ( "%d: [%v]\n", i, v )
	//}

	pw := NewParseWords ()
	pw.SetDebug(db)
	for i, v := range testCases {
		if ( db ) { fmt.Printf ( "Running %d\n", i ); }
		pw.SetLine ( v.In )
		pw.SetOptions ( "C", v.kq, v.kb )
		rv := pw.GetWords()
		// result := Format(v.format, v.value)
		if ! arrayEq ( rv, v.Out ) {
			t.Fatalf("Error at [%d]: for=[%s] results=[%s] expected=[%s]", i, v.In, sVar(rv), sVar(v.Out) )
		}
	}
}

/* vim: set noai ts=4 sw=4: */
