package xlsxformula

import (
	"bytes"
	"fmt"
)

type NodeType int

const (
	Function NodeType = iota
	Expression
	Compare
	SingleToken
)

func (nt NodeType) String() string {
	switch nt {
	case Function:
		return "Function"
	case Expression:
		return "Expression"
	case Compare:
		return "Compare"
	case SingleToken:
		return "SingleToken"
	}
	return "Unknown"
}

type Node struct {
	Children []*Node
	Token    *Token
	Type     NodeType
}

func (node Node) String() string {
	switch node.Type {
	case Function:
		var buffer bytes.Buffer
		buffer.WriteString(node.Token.Text)
		buffer.WriteByte('(')
		for i, child := range node.Children {
			if i != 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(child.String())
		}
		buffer.WriteByte(')')
		return buffer.String()
	case Expression:
		var buffer bytes.Buffer
		buffer.WriteByte('(')
		for i, child := range node.Children {
			if i != 0 {
				buffer.WriteByte(' ')
			}
			buffer.WriteString(child.String())
		}
		buffer.WriteByte(')')
		return buffer.String()
	case Compare:
		var buffer bytes.Buffer
		buffer.WriteString(node.Children[0].String())
		buffer.WriteByte(' ')
		buffer.WriteString(node.Token.Text)
		buffer.WriteByte(' ')
		buffer.WriteString(node.Children[1].String())
		return buffer.String()
	case SingleToken:
		return node.Token.Text
	}
	return ""
}

var nullToken *Token = &Token{
	Type: Null,
}

func get(tokens []*Token, index int) *Token {
	if index < len(tokens) {
		return tokens[index]
	}
	return nullToken
}

func Parse(formula string) (*Node, error) {
	tokens, err := Tokenize(formula)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("Formula is empty")
	}

	i := 0

	currentNode := &Node{
		Type: Expression,
	}
	stack := []*Node{currentNode}

	acceptValue := true

	for i < len(tokens) {
		token := tokens[i]
		switch token.Type {
		case Name:
			if !acceptValue {
				return nil, fmt.Errorf("Unexpected name '%s' appears at %d:%d", token.Text, token.Line, token.Col)
			}
			if get(tokens, i+1).Type == LParen {
				next := &Node{
					Token: token,
					Type:  Function,
				}
				currentNode.Children = append(currentNode.Children, next)
				if get(tokens, i+2).Type == RParen {
					acceptValue = false
					i += 3
				} else {
					param := &Node{
						Type: Expression,
					}
					next.Children = append(next.Children, param)
					stack = append(stack, next, param)
					currentNode = param
					i += 2
				}
			} else {
				currentNode.Children = append(currentNode.Children, &Node{
					Type:  SingleToken,
					Token: token,
				})
				acceptValue = false
				i++
			}
		case Comma:
			if len(stack) < 2 {
				return nil, fmt.Errorf("Unexpected comma ',' appears at %d:%d", token.Text, token.Line, token.Col)
			}
			parentFunction := stack[len(stack)-2]
			if parentFunction.Type != Function {
				return nil, fmt.Errorf("Unexpected comma ',' appears outside of function arguments at %d:%d", token.Text, token.Line, token.Col)
			}
			nextParam := &Node{
				Type: Expression,
			}
			stack[len(stack)-1] = nextParam
			parentFunction.Children = append(parentFunction.Children, nextParam)
			currentNode = nextParam
			acceptValue = true
			i++
		case Range:
			if !acceptValue {
				return nil, fmt.Errorf("Unexpected range '%s' appears at %d:%d", token.Text, token.Line, token.Col)
			}
			currentNode.Children = append(currentNode.Children, &Node{
				Type:  SingleToken,
				Token: token,
			})
			acceptValue = false
			i++
		case Bool:
			if !acceptValue {
				return nil, fmt.Errorf("Unexpected boolean value '%s' appears at %d:%d", token.Text, token.Line, token.Col)
			}
			currentNode.Children = append(currentNode.Children, &Node{
				Type:  SingleToken,
				Token: token,
			})
			acceptValue = false
			i++
		case Integer:
			if !acceptValue {
				return nil, fmt.Errorf("Unexpected integer '%s' appears at %d:%d", token.Text, token.Line, token.Col)
			}
			currentNode.Children = append(currentNode.Children, &Node{
				Type:  SingleToken,
				Token: token,
			})
			acceptValue = false
			i++
		case String:
			if !acceptValue {
				return nil, fmt.Errorf("Unexpected string '%s' appears at %d:%d", token.Text, token.Line, token.Col)
			}
			currentNode.Children = append(currentNode.Children, &Node{
				Type:  SingleToken,
				Token: token,
			})
			acceptValue = false
			i++
		case Operator:
			if acceptValue {
				if len(currentNode.Children) == 0 && isNumber(get(tokens, i+1)) {
					nextToken := tokens[i+1]
					currentNode.Children = append(currentNode.Children, &Node{
						Type: SingleToken,
						Token: &Token{
							Type: nextToken.Type,
							Text: "-" + nextToken.Text,
							Line: token.Line,
							Col:  token.Col,
						},
					})
					i += 2
				} else {
					return nil, fmt.Errorf("Unexpected operator '%s' appears at %d:%d", token.Text, token.Line, token.Col)
				}
			} else {
				currentNode.Children = append(currentNode.Children, &Node{
					Type:  SingleToken,
					Token: token,
				})
				acceptValue = true
				i++
			}
		case Comparator:
			if acceptValue {
				return nil, fmt.Errorf("Unexpected comparator '%s' appears at %d:%d", token.Text, token.Line, token.Col)
			}
			currentNode.Children = append(currentNode.Children, &Node{
				Type:  SingleToken,
				Token: token,
			})
			acceptValue = true
			i++
		case LParen:
			if !acceptValue {
				return nil, fmt.Errorf("Unexpected left paren '(' appears at %d:%d", token.Line, token.Col)
			}
			if i+1 == len(tokens) {
				return nil, fmt.Errorf("Right paren ')' is missing. Orphan left paren appears at %d:%d, ", token.Line, token.Col)
			}
			nextToken := tokens[i+1]
			if nextToken.Type == RParen {
				return nil, fmt.Errorf("Empty paren blocks '()' at %d:%d", token.Line, token.Col)
			}
			next := &Node{
				Token: token,
				Type:  Expression,
			}
			currentNode.Children = append(currentNode.Children, next)
			stack = append(stack, next)
			currentNode = next
			i++
		case RParen:
			if acceptValue || len(stack) == 1 {
				return nil, fmt.Errorf("Unexpected right paren ')' appears at %d:%d", token.Line, token.Col)
			}
			if len(stack) > 1 && stack[len(stack)-2].Type == Function {
				function := stack[len(stack)-2]
				for i, param := range function.Children {
					function.Children[i] = clean(param)
				}
				stack = stack[:len(stack)-2]
				currentNode = stack[len(stack)-1]
			} else {
				lastNode := stack[len(stack)-1]
				currentNode = stack[len(stack)-2]
				currentNode.Children[len(currentNode.Children)-1] = clean(lastNode)
				stack = stack[:len(stack)-1]
			}
			i++
		default:
			return nil, fmt.Errorf("Unknown token: %s", token.Type.String())
		}
	}
	if len(stack) > 1 {
		return nil, fmt.Errorf("The following nest defined at %d:%d is not closed yet: %s", stack[1].Token.Line, stack[1].Token.Col, stack[1].String())
	}
	return clean(stack[0]), nil
}

func clean(node *Node) *Node {
	for {
		if node.Type == Expression && len(node.Children) == 1 {
			node = node.Children[0]
		} else {
			break
		}
	}
	return node
}

func isNumber(token *Token) bool {
	return token.Type == Integer || token.Type == Float
}

var isValue map[TokenType]bool = map[TokenType]bool{
	Integer: true,
	Float:   true,
	String:  true,
	Bool:    true,
	Range:   true,
	Name:    true,
}
