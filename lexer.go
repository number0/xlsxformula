package xlsxformula

import (
	"fmt"
	"regexp"
	"strconv"
)

type TokenType int

const (
	Number     TokenType = iota // number
	String                      // double quoted string
	Bool                        // TRUE/FALSE
	Operator                    // +, -, *, /, ^, &
	LParen                      // (
	RParen                      // )
	Comma                       // ,
	Comparator                  // =, <>, <, >, <=, >=
	Name                        // function name, named range etc
	Range                       // A2:B3
	Null
)

func (tt TokenType) String() string {
	switch tt {
	case Number:
		return "Number"
	case String:
		return "String"
	case Bool:
		return "Bool"
	case Operator:
		return "Operator"
	case LParen:
		return "LParen"
	case RParen:
		return "RParen"
	case Comma:
		return "Comma"
	case Comparator:
		return "Comparator"
	case Name:
		return "Name"
	case Range:
		return "Range"
	}
	return "Unknown"
}

type Token struct {
	Type TokenType
	Text string
	Line int
	Col  int
}

var rangePattern *regexp.Regexp = regexp.MustCompile(`^(\$?[A-Z]+\$?[1-9][0-9]*)(:(\$?[A-Z]+|\$?[1-9][0-9]*|\$?[A-Z]+\$?[1-9][0-9]*))?$`)

var symbolSeparator map[rune]bool = map[rune]bool{
	' ':  true,
	'+':  true,
	'-':  true,
	'*':  true,
	'/':  true,
	'^':  true,
	'&':  true,
	'(':  true,
	')':  true,
	',':  true,
	'<':  true,
	'>':  true,
	'=':  true,
	'\r': true,
	'\n': true,
}

var singleCharNode map[rune]TokenType = map[rune]TokenType{
	'+': Operator,
	'-': Operator,
	'*': Operator,
	'/': Operator,
	'^': Operator,
	'&': Operator,
	'(': LParen,
	')': RParen,
	'=': Comparator,
	',': Comma,
}

func Tokenize(formula string) ([]*Token, error) {
	tokens := []*Token{}

	source := []rune(formula)
	index := 0
	line := 1
	lineHead := 0
	for index < len(source) {
		ch := source[index]
		if nodeType, ok := singleCharNode[ch]; ok {
			tokens = append(tokens, &Token{
				Type: nodeType,
				Text: string(ch),
				Line: line,
				Col:  index - lineHead + 1,
			})
			index++
			continue
		}
		switch ch {
		case ' ':
			index++
			continue
		case '\r':
			if index+1 < len(source) && source[index+1] == '\n' {
				index += 2
			} else {
				index++
			}
			line++
			lineHead = index
			continue
		case '<':
			if index+1 < len(source) {
				switch source[index+1] {
				case '>':
					tokens = append(tokens, &Token{
						Type: Comparator,
						Text: "<>",
						Line: line,
						Col:  index - lineHead + 1,
					})
					index += 2
					continue
				case '=':
					tokens = append(tokens, &Token{
						Type: Comparator,
						Text: "<=",
						Line: line,
						Col:  index - lineHead + 1,
					})
					index += 2
					continue
				}
			}
			tokens = append(tokens, &Token{
				Type: Comparator,
				Text: "<",
				Line: line,
				Col:  index - lineHead + 1,
			})
			index++
			continue
		case '>':
			if index+1 < len(source) && source[index+1] == '=' {
				tokens = append(tokens, &Token{
					Type: Comparator,
					Text: ">=",
					Line: line,
					Col:  index - lineHead + 1,
				})
				index += 2
			} else {
				tokens = append(tokens, &Token{
					Type: Comparator,
					Text: ">",
					Line: line,
					Col:  index - lineHead + 1,
				})
				index++
			}
			continue
		case '"':
			last := index + 1
			found := false
			for last < len(source) {
				if source[last] == '"' {
					tokens = append(tokens, &Token{
						Type: String,
						Text: string(source[index+1 : last]),
						Line: line,
						Col:  index - lineHead + 1,
					})
					index = last + 1
					found = true
					break
				}
				last++
			}
			if !found {
				return tokens, fmt.Errorf(`closing double quoatation is missing: %s`, string(source[index:]))
			}
		default:
			start := index
			last := index
			found := false
			var text string
			for last < len(source) {
				if symbolSeparator[source[last]] {
					text = string(source[index:last])
					index = last
					found = true
					break
				}
				last++
			}
			if !found {
				text = string(source[index:])
				index = len(source)
			}
			if _, err := strconv.ParseFloat(text, 64); err == nil {
				tokens = append(tokens, &Token{
					Type: Number,
					Text: text,
					Line: line,
					Col:  start - lineHead + 1,
				})
			} else if rangePattern.MatchString(text) {
				tokens = append(tokens, &Token{
					Type: Range,
					Text: text,
					Line: line,
					Col:  start - lineHead + 1,
				})
			} else if text == "TRUE" || text == "FALSE" {
				tokens = append(tokens, &Token{
					Type: Bool,
					Text: text,
					Line: line,
					Col:  start - lineHead + 1,
				})
			} else {
				tokens = append(tokens, &Token{
					Type: Name,
					Text: text,
					Line: line,
					Col:  start - lineHead + 1,
				})
			}
		}
	}

	return tokens, nil
}
