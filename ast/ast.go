package ast

import (
	"github.com/tamurayoshiya/monkey/token"
)

// ASTノード共通インターフェース
type Node interface {
	TokenLiteral() string // デバッグ用
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

// -----------------------------------------------------

// Identifier

type Identifier struct {
	Token token.Token // token.IDENT トークン
	Value string
}

func (i *Identifier) expressionNode() {
}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// -----------------------------------------------------

// let ステートメント（文）
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

// -----------------------------------------------------

// return ステートメント（文）
// 構造: return <expression>;

type ReturnStatement struct {
	Token       token.Token // 'return' トークン
	ReturnValue Expression
}

func (ls *ReturnStatement) statementNode() {
}
func (ls *ReturnStatement) TokenLiteral() string {
	return ls.Token.Literal
}
