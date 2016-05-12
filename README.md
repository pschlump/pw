pw: parse words
===============
 
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/pschlump/Go-FTL/master/LICENSE)


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


