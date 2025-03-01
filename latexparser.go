package gosymbol

import (
	"errors"
	"slices"
	"strings"
)

type treeParser struct {
	operator string
	sons     []treeParser
}

func findMatchingBrackets(formula string) (int, error) {
	var matchingBracket rune
	bracket := rune(formula[0])
	switch bracket {
	case '(':
		matchingBracket = ')'
	case '[':
		matchingBracket = ']'
	case '{':
		matchingBracket = '}'
	}
	count := 0
	for i, c := range formula[1:] {
		if c == bracket {
			count++
		}
		if c == matchingBracket && count > 0 {
			count--
		}
		if c == matchingBracket && count == 0 {
			return i + 1, nil
		}
	}
	return 0, errors.New("no matching bracket")
}

func ParseLatex(formula string) (treeParser, error) {
	strings.ReplaceAll(formula, " ", "")
	var parseTree treeParser
	operators := []byte{'+', '-', '*', '/', '^'}
	openBrackets := []byte{'(', '[', '{'}
	i := 0
	var variable string
	for i < len(formula) {
		if slices.Contains(operators, formula[i]) {
			if variable != string("") {
				parseTree.sons = append(parseTree.sons, treeParser{variable, nil})
				variable = string("")
			}
			parseTree.operator = string(formula[i])
			variable = string("")
		} else if slices.Contains(openBrackets, formula[i]) {
			if variable != string("") {
				parseTree.sons = append(parseTree.sons, treeParser{variable, nil})
				variable = string("")
			}
			variable = string("")
			finalIndex, err := findMatchingBrackets(formula[i:])
			if err != nil {
				return treeParser{}, err
			}
			parsedSon, err := ParseLatex(formula[i+1 : i+finalIndex])
			if err != nil {
				return treeParser{}, err
			}
			parseTree.sons = append(parseTree.sons, parsedSon)
			i = finalIndex
		} else {
			variable = variable + string(formula[i])
		}
		i++
	}
	if variable != string("") {
		parseTree.sons = append(parseTree.sons, treeParser{variable, nil})
	}
	return parseTree, nil
}
