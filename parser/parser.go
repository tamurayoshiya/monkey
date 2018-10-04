package parser

import (
	"fmt"
	"strconv"

	"github.com/tamurayoshiya/monkey/ast"
	"github.com/tamurayoshiya/monkey/lexer"
	"github.com/tamurayoshiya/monkey/token"
)

// -------------------------------------------------------

// 演算子の優先順位

const (
	// 次に来る定数にインクリメントしながら数を与える
	// _ = 0, LOWEST = 1, EQUALS = 2... と割り当てられる
	_ int = iota
	LOWEST
	EQUALS     // ==
	LESSGRETER // > または <
	SUM        // +
	PRODUCT    // *
	PREFIX     // -X または !X
	CALL       // myFunction(X)
)

// -------------------------------------------------------

// Pratt構文解析器

type (
	prefixParseFn func() ast.Expression               // 前置構文解析関数 (prefix parsing function)
	infixParseFn  func(ast.Expression) ast.Expression // 中置構文解析関数 (infix parsing function)
)

// -------------------------------------------------------

type Parser struct {
	l      *lexer.Lexer // 字句解析器インスタンスへのポインタ
	errors []string

	curToken  token.Token // 現在のトークン(cur -> current)
	peekToken token.Token // 次のトークン(peek 覗く)

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// 2つのトークンを読み込む。curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

// -------------------------------------------------------

// トークンを進める処理

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// -------------------------------------------------------

// パースの主な処理

func (p *Parser) ParseProgram() *ast.Program {
	// ASTのルートノードを生成
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// トークンをウォーク
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	// ルートノードを返却
	return program
}

// -------------------------------------------------------

// 文、式文のパース

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// -------------------------------------------------------

// 式のパース

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

// 識別子のパース
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// 整数リテラルのパース
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// -------------------------------------------------------

// ヘルパー関数

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// -------------------------------------------------------

// アサーション関数

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// -------------------------------------------------------

// 外部にエラーをエクスポート

func (p *Parser) Errors() []string {
	return p.errors
}
