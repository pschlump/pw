pw: parse words
===============

This is for parsing a string into a  set of words.  For Example:

	import "github.com/pschlump/pw"

	/* ... */

	func ParseLineIntoWords(line string) []string {
		Pw := pw.NewParseWords()
		Pw.SetOptions("C", true, true)
		Pw.SetLine(line)
		rv := Pw.GetWords()
		return rv
	}

Will take a line like this one and return an array of words.


