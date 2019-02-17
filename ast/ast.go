package ast

import (
	"bytes"
	"strings"

	"github.com/tamurayoshiya/monkey/token"
)

// ASTノード共通インターフェース
type Node interface {
	TokenLiteral() string // デバッグ用
	String() string       // デバッグ用
}

// 文インターフェース、文は値を生成しない
type Statement interface {
	Node
	statementNode()
}

// 式インターフェース、式は値を生成する
type Expression interface {
	Node
	expressionNode()
}

// -----------------------------------------------------

// 文: 値を生成しない
// let <identifier> = <expression>;
//
// 式: 値を生成する
// <expression>
//
// Monkey言語の「文」はletとreturnのみ、それ以外はすべて式

// Programノード
// プログラムは文の集合
// 構文解析器が生成するASTのルートノード
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// -----------------------------------------------------

// 識別子

type Identifier struct {
	Token token.Token // token.IDENT トークン
	Value string
}

func (i *Identifier) expressionNode() {
}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

// -----------------------------------------------------

// 整数リテラル

type IntegerLiteral struct {
	Token token.Token // token.IDENT トークン
	Value int64
}

func (il *IntegerLiteral) expressionNode() {
}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// -----------------------------------------------------

// 真偽値リテラル

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {
}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

// -----------------------------------------------------

// let文
// 構造: let <identifier> = <expression>;

type LetStatement struct {
	Token token.Token // token.LET トークン
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {
}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// -----------------------------------------------------

// return文
// 構造: return <expression>;

type ReturnStatement struct {
	Token       token.Token // 'return' トークン
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {
}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// -----------------------------------------------------

// block文

type BlockStatement struct {
	Token      token.Token // トークン
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {
}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// -----------------------------------------------------

// 式文 ステートメント
//
// e.g.
// let x = y + 1; // 式
// y + 1;		  // 式文
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {
}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// -----------------------------------------------------

// 前置演算子 Expression

type PrefixExpression struct {
	Token    token.Token // 前置トークン、例えば"!", "-"
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// -----------------------------------------------------

// 中置演算子 Expression

type InfixExpression struct {
	Token    token.Token // 演算子トークン、例えば"+"
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {
}

func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.Literal
}
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

// -----------------------------------------------------

// if式

type IfExpression struct {
	Token       token.Token // 'if' トークン
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {
}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// -----------------------------------------------------

// 関数リテラル

type FunctionLiteral struct {
	Token      token.Token // 'fn' トークン
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {
}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// -----------------------------------------------------

// 関数呼び出し

type CallExpression struct {
	Token     token.Token // '(' トークン
	Function  Expression  // Identifier または FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {
}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// -----------------------------------------------------

// 文字列リテラル

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// -----------------------------------------------------

// 配列リテラル

type ArrayLiteral struct {
	Token    token.Token // '['トークン
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {
}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// -----------------------------------------------------

// 添字演算子式の構文解析

type IndexExpression struct {
	Token token.Token // '['トークン
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
