package parser

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/daneofmanythings/calcuroller/pkg/interpreter/ast"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/lexer"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/token"
)

// type alias for the testing suite
type dice string

const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	EXPONENT
	PREFIX
)

var precedences = map[token.TokenType]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.MODULUS:  PRODUCT,
	token.CARET:    EXPONENT,
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	dicemodParseFn func(*ast.DiceLiteral)
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	dicemodParseFns map[token.TokenType]dicemodParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.DICE, p.parseDiceExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.ILLEGAL, p.parseIllegalExpression)
	p.registerPrefix(token.ASTERISK, p.parseIllegalExpression)
	p.registerPrefix(token.SLASH, p.parseIllegalExpression)
	p.registerPrefix(token.MODULUS, p.parseIllegalExpression)
	p.registerPrefix(token.EOF, p.parseIllegalExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.CARET, p.parseInfixExpression)
	p.registerInfix(token.MODULUS, p.parseInfixExpression)

	p.dicemodParseFns = make(map[token.TokenType]dicemodParseFn)
	p.registerDicemod(token.METATAG, p.parseDiceTag)
	p.registerDicemod(token.DICEQUANT, p.parseDiceQuant)
	p.registerDicemod(token.DICEMIN, p.parseDiceMin)
	p.registerDicemod(token.DICEMAX, p.parseDiceMax)
	p.registerDicemod(token.DICEKEEPLOWEST, p.parseDiceLowest)
	p.registerDicemod(token.DICEKEEPHIGHEST, p.parseDiceHighest)

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerDicemod(tokenType token.TokenType, fn dicemodParseFn) {
	p.dicemodParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	return p.parseExpressionStatement()
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseIllegalExpression() ast.Expression {
	expression := &ast.IllegalLiteral{
		Token:   p.curToken,
		Literal: p.curToken.Literal,
	}

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken, Tags: []string{}}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	for p.peekToken.Type == token.METATAG {
		p.nextToken()
		lit.Tags = append(lit.Tags, p.curToken.Literal)
	}

	return lit
}

// NOTE: This does not happily recurse like the rest of the parser
func (p *Parser) parseDiceExpression() ast.Expression {
	dice := &ast.DiceLiteral{
		Token: p.curToken,
		Tags:  []string{},
		Size:  p.validateIntToUInt(p.curToken.Literal),
	}

	for slices.Contains(token.DiceMods, p.peekToken.Type) {
		p.nextToken()
		p.dicemodParseFns[p.curToken.Type](dice)
	}

	return dice
}

// TODO: collapse this logic
func (p *Parser) parseDiceTag(d *ast.DiceLiteral) {
	d.Tags = append(d.Tags, p.curToken.Literal)
}

func (p *Parser) parseDiceQuant(d *ast.DiceLiteral) {
	d.Quantity = p.validateIntToUInt(p.curToken.Literal)
}

func (p *Parser) parseDiceMin(d *ast.DiceLiteral) {
	d.MinValue = p.validateIntToUInt(p.curToken.Literal)
}

func (p *Parser) parseDiceMax(d *ast.DiceLiteral) {
	d.MaxValue = p.validateIntToUInt(p.curToken.Literal)
}

func (p *Parser) parseDiceLowest(d *ast.DiceLiteral) {
	d.KeepLowest = p.validateIntToUInt(p.curToken.Literal)
}

func (p *Parser) parseDiceHighest(d *ast.DiceLiteral) {
	d.KeepHighest = p.validateIntToUInt(p.curToken.Literal)
}

func (p *Parser) validateIntToUInt(lit string) uint32 {
	value, err := strconv.ParseInt(lit, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}

	return uint32(value)
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
