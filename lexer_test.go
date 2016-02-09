package xlsxformula

import (
	"testing"
)

func TestOperatorsAndNumbers(t *testing.T) {
	tokens, err := Tokenize(" 10 + 20 ")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[0].Type != Integer || tokens[1].Type != Operator || tokens[2].Type != Integer {
		t.Errorf("Node type is wrong: %s %s %s", tokens[0].Type.String(), tokens[1].Type.String(), tokens[2].Type.String())
	}
}

func TestOperatorsAndFloats(t *testing.T) {
	tokens, err := Tokenize(" 10.5 - 20.6 * 30.1 / 40.4")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 7 {
		t.Errorf("Parse() should return 7 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[0].Type != Float || tokens[2].Type != Float || tokens[4].Type != Float || tokens[6].Type != Float {
		t.Errorf("Node type is wrong: %s %s %s %s", tokens[0].Type.String(), tokens[2].Type.String(), tokens[4].Type.String(), tokens[6].Type.String())
	}
	if tokens[1].Type != Operator || tokens[3].Type != Operator || tokens[5].Type != Operator {
		t.Errorf("Node type is wrong: %s %s %s", tokens[1].Type.String(), tokens[3].Type.String(), tokens[5].Type.String())
	}
}

func TestJoinOperatorAndString(t *testing.T) {
	tokens, err := Tokenize(`"hello"&"world"`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[0].Type != String || tokens[1].Type != Operator || tokens[2].Type != String {
		t.Errorf("Node type is wrong: %s %s %s", tokens[0].Type.String(), tokens[1].Type.String(), tokens[2].Type.String())
	}
}

func TestRangeAndName(t *testing.T) {
	tokens, err := Tokenize(`A1 ^ VARIABLE`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[0].Type != Range || tokens[1].Type != Operator || tokens[2].Type != Name {
		t.Errorf("Node type is wrong: %s %s %s", tokens[0].Type.String(), tokens[1].Type.String(), tokens[2].Type.String())
	}
}

func TestCompare_1(t *testing.T) {
	tokens, err := Tokenize(`1 = 10`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Type != Comparator || tokens[1].Text != "=" {
		t.Errorf("Compare operator is wrong: %s %s", tokens[1].Type.String(), tokens[1].Text)
	}
}

func TestCompare_2(t *testing.T) {
	tokens, err := Tokenize(`1 <> 10`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Type != Comparator || tokens[1].Text != "<>" {
		t.Errorf("Compare operator is wrong: %s %s", tokens[1].Type.String(), tokens[1].Text)
	}
}

func TestCompare_3(t *testing.T) {
	tokens, err := Tokenize(`1 < 10`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Type != Comparator || tokens[1].Text != "<" {
		t.Errorf("Compare operator is wrong: %s %s", tokens[1].Type.String(), tokens[1].Text)
	}
}

func TestCompare_4(t *testing.T) {
	tokens, err := Tokenize(`1 > 10`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Type != Comparator || tokens[1].Text != ">" {
		t.Errorf("Compare operator is wrong: %s %s", tokens[1].Type.String(), tokens[1].Text)
	}
}

func TestCompare_5(t *testing.T) {
	tokens, err := Tokenize(`1 <= 10`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Type != Comparator || tokens[1].Text != "<=" {
		t.Errorf("Compare operator is wrong: %s %s", tokens[1].Type.String(), tokens[1].Text)
	}
}

func TestCompare_6(t *testing.T) {
	tokens, err := Tokenize(`1 >= 10`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 3 {
		t.Errorf("Parse() should return 3 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Type != Comparator || tokens[1].Text != ">=" {
		t.Errorf("Compare operator is wrong: %s %s", tokens[1].Type.String(), tokens[1].Text)
	}
}

func TestBoolAndParenAndComma(t *testing.T) {
	tokens, err := Tokenize(`IF(TRUE, "test")`)
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 6 {
		t.Errorf("Parse() should return 6 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Type != LParen || tokens[2].Type != Bool || tokens[3].Type != Comma || tokens[4].Type != String || tokens[5].Type != RParen {
		t.Errorf("Node types are wrong: %s %s %s %s %s %s",
			tokens[0].Type.String(), tokens[1].Type.String(), tokens[2].Type.String(),
			tokens[3].Type.String(), tokens[4].Type.String(), tokens[5].Type.String())
	}
}

func TestMultiline(t *testing.T) {
	tokens, err := Tokenize("IF(\rTRUE,\r\n\"test\")")
	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
	if len(tokens) != 6 {
		t.Errorf("Parse() should return 6 tokens, but %d tokens", len(tokens))
		return
	}
	if tokens[1].Line != 1 || tokens[2].Line != 2 || tokens[3].Line != 2 || tokens[4].Line != 3 || tokens[5].Line != 3 {
		t.Errorf("Nodes' lines are wrong: %d %d %d %d %d %d",
			tokens[0].Line, tokens[1].Line, tokens[2].Line,
			tokens[3].Line, tokens[4].Line, tokens[5].Line)
	}
	if tokens[1].Col != 3 || tokens[2].Col != 1 || tokens[3].Col != 5 || tokens[4].Col != 1 || tokens[5].Col != 7 {
		t.Errorf("Nodes' columns are wrong: %d %d %d %d %d %d",
			tokens[0].Col, tokens[1].Col, tokens[2].Col,
			tokens[3].Col, tokens[4].Col, tokens[5].Col)
	}
}
