package xlsxformula

import (
	"testing"
)

func TestParseSingleValue(t *testing.T) {
	node, err := Parse("10")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != SingleToken {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if node.Token.Text != "10" {
		t.Errorf("node.Token.Text is wrong: %s", node.Token.Text)
	}
}

func TestParseNegativeSingleValue(t *testing.T) {
	node, err := Parse("-10")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != SingleToken {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if node.Token.Text != "-10" {
		t.Errorf("node.Token.Text is wrong: %s", node.Token.Text)
	}
}

func TestParseSimpleExpression(t *testing.T) {
	node, err := Parse("10 + 20")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != Expression {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if len(node.Children) != 3 {
		t.Errorf("child count should be 3, but %d", len(node.Children))
	} else if node.Children[0].Type != SingleToken || node.Children[1].Type != SingleToken || node.Children[2].Type != SingleToken {
		t.Errorf("node Types are wrong: %s, %s, %s", node.Children[0].Type.String(), node.Children[1].Type.String(), node.Children[2].Type.String())
	}
}

func TestParseLongExpression(t *testing.T) {
	node, err := Parse("10 + 20 / 40 * 50 ^ 2")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != Expression {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if len(node.Children) != 9 {
		t.Errorf("child count should be 9, but %d", len(node.Children))
	} else if node.Children[0].Type != SingleToken || node.Children[1].Type != SingleToken || node.Children[2].Type != SingleToken {
		t.Errorf("node Types are wrong: %s, %s, %s", node.Children[0].Type.String(), node.Children[1].Type.String(), node.Children[2].Type.String())
	}
}

func TestParseSimpleCompare(t *testing.T) {
	node, err := Parse(`FALSE <> "20"`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != Expression {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if len(node.Children) != 3 {
		t.Errorf("child count should be 3, but %d", len(node.Children))
	} else if node.Children[0].Type != SingleToken || node.Children[1].Type != SingleToken || node.Children[2].Type != SingleToken {
		t.Errorf("node Types are wrong: %s, %s, %s", node.Children[0].Type.String(), node.Children[1].Type.String(), node.Children[2].Type.String())
	}
}

func TestParseSingleValueWithParen(t *testing.T) {
	node, err := Parse("(((10)))")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != SingleToken {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if node.Token.Text != "10" {
		t.Errorf("node.Token.Text is wrong: %s", node.Token.Text)
	}
}

func TestParseNestedExpressionWithParen(t *testing.T) {
	node, err := Parse("((((10+20))))")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != Expression {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if len(node.Children) != 3 {
		t.Errorf("child count should be 3, but %d", len(node.Children))
	} else if node.Children[0].Type != SingleToken {
		t.Errorf("node Types are wrong: %s", node.Children[0].Type.String())
	}
}

func TestParseFunctionWithoutArguments(t *testing.T) {
	node, err := Parse("TODAY()")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != Function {
		t.Errorf("function Types are wrong: %s", node.Type.String())
	} else if len(node.Children) != 0 {
		t.Errorf("function should not have arguments but it has %d arguments", len(node.Children))
	} else if node.Token.Text != "TODAY" {
		t.Errorf("function name is wrong: %s", node.Token.Text)
	}
}

func TestParseFunctionWithArguments(t *testing.T) {
	node, err := Parse(`IF(10 < E2, "bigger", "smaller")`)
	t.Log(node.String())
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != Function {
		t.Errorf("function Types are wrong: %s", node.Type.String())
	} else if len(node.Children) != 3 {
		t.Errorf("function should not have arguments but it has %d arguments", len(node.Children))
	} else if node.Token.Text != "IF" {
		t.Errorf("function name is wrong: %s", node.Token.Text)
	} else if node.Children[0].Type != Expression || node.Children[1].Type != SingleToken || node.Children[2].Type != SingleToken {
		t.Errorf("argument node types are wrong: %s %s %s", node.Children[0].Type.String(), node.Children[1].Type.String(), node.Children[2].Type.String())
	}
}

/*func TestParseFunctionWithArguments(t *testing.T) {
	node, err := Parse("IF()")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	} else if node.Type != Expression {
		t.Errorf("node.Type is wrong: %s", node.Type.String())
	} else if len(node.Children) != 1 {
		t.Errorf("child count should be 3, but %d", len(node.Children))
	} else {
		function := node.Children[0]
		if function.Type != Function {
			t.Errorf("function Types are wrong: %s", node.Children[0].Type.String())
		} else if len(function.Children) != 0 {
			t.Errorf("function should not have arguments but it has %d arguments", len(function.Children))
		}
	}
}*/

func TestParseErrorSequentialValues(t *testing.T) {
	_, err := Parse("10 20")
	if err == nil {
		t.Errorf("err should not be nil")
	}
}

func TestParseErrorSequentialOperators(t *testing.T) {
	_, err := Parse("10 + +")
	if err == nil {
		t.Errorf("err should not be nil")
	}
}

func TestParseErrorTooManyLParen(t *testing.T) {
	_, err := Parse("((10")
	if err == nil {
		t.Errorf("err should not be nil")
	}
}

func TestParseErrorTooManyRParen(t *testing.T) {
	_, err := Parse("10)")
	if err == nil {
		t.Errorf("err should not be nil")
	}
}
