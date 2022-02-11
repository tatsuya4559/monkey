package parser

import (
	"fmt"
	"strconv"

	"github.com/tatsuya4559/monkey/ast"
	"github.com/tatsuya4559/monkey/lexer"
	"github.com/tatsuya4559/monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // <, >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X, !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.MOD:      PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.MACRO, p.parseMacroLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// read token for setup
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) newPeekError(t token.TokenType) error {
	return fmt.Errorf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken

	tok := p.l.NextToken()
	// consume comments
	for tok.Type == token.COMMENT {
		tok = p.l.NextToken()
	}
	p.peekToken = tok
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt, err := p.parseStatement()
		// FIXME: return error
		if err == nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.curToken}

	if err := p.expectPeek(token.IDENT); err != nil {
		return nil, err
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if err := p.expectPeek(token.ASSIGN); err != nil {
		return nil, err
	}

	p.nextToken()

	var err error
	stmt.Value, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.SEMICOLON); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) error {
	if p.peekTokenIs(t) {
		p.nextToken()
		return nil
	}
	p.peekError(t)
	return p.newPeekError(t)
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	var err error
	stmt.ReturnValue, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.SEMICOLON); err != nil {
		return nil, err
	}

	return stmt, nil
}

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	var err error
	stmt.Expression, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) newNoPrefixParseFnError(t token.TokenType) error {
	return fmt.Errorf("no prefix parse function for %s found", t)
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		// FIXME: remove p.errors
		p.noPrefixParseFnError(p.curToken.Type)
		return nil, p.newNoPrefixParseFnError(p.curToken.Type)
	}
	leftExpr, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpr, nil
		}

		p.nextToken()

		leftExpr, err = infix(leftExpr)
		if err != nil {
			return nil, err
		}
	}

	return leftExpr, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil, fmt.Errorf("could not parse %q as integer", p.curToken.Literal)
	}

	lit.Value = value
	return lit, nil
}

func (p *Parser) parseStringLiteral() (ast.Expression, error) {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	var err error
	expr.Right, err = p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	var err error
	expr.Right, err = p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseBoolean() (ast.Expression, error) {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}, nil
}

func (p *Parser) parseGroupExpression() (ast.Expression, error) {
	p.nextToken()

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.RPAREN); err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	expr := &ast.IfExpression{Token: p.curToken}

	// skip IF
	p.nextToken()
	var err error
	expr.Condition, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.LBRACE); err != nil {
		return nil, err
	}

	expr.Consequence, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if err := p.expectPeek(token.LBRACE); err != nil {
			return nil, err
		}

		expr.Alternative, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken() // skip LBRACE

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expression, error) {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if err := p.expectPeek(token.LPAREN); err != nil {
		return nil, err
	}

	var err error
	lit.Parameters, err = p.parseFuntionParameters()
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.LBRACE); err != nil {
		return nil, err
	}

	lit.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return lit, nil
}

func (p *Parser) parseFuntionParameters() ([]*ast.Identifier, error) {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers, nil
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if err := p.expectPeek(token.RPAREN); err != nil {
		return nil, err
	}

	return identifiers, nil
}

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	expr := &ast.CallExpression{Token: p.curToken, Function: function}
	var err error
	expr.Arguments, err = p.parseExpressionList(token.RPAREN)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) parseArrayLiteral() (ast.Expression, error) {
	array := &ast.ArrayLiteral{Token: p.curToken}
	var err error
	array.Elements, err = p.parseExpressionList(token.RBRACKET)
	if err != nil {
		return nil, err
	}
	return array, nil
}

func (p *Parser) parseExpressionList(end token.TokenType) ([]ast.Expression, error) {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list, nil
	}

	p.nextToken()
	item, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	list = append(list, item)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		item, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	if err := p.expectPeek(end); err != nil {
		return nil, err
	}

	return list, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	expr := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	var err error
	expr.Index, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.RBRACKET); err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseHashLiteral() (ast.Expression, error) {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		if err := p.expectPeek(token.COLON); err != nil {
			return nil, err
		}

		p.nextToken()
		value, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) {
			if err := p.expectPeek(token.COMMA); err != nil {
				return nil, err
			}
		}
	}

	if err := p.expectPeek(token.RBRACE); err != nil {
		return nil, err
	}

	return hash, nil
}

func (p *Parser) parseMacroLiteral() (ast.Expression, error) {
	lit := &ast.MacroLiteral{Token: p.curToken}

	if err := p.expectPeek(token.LPAREN); err != nil {
		return nil, err
	}

	var err error
	lit.Parameters, err = p.parseFuntionParameters()
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.LBRACE); err != nil {
		return nil, err
	}

	lit.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return lit, nil
}

func (p *Parser) parseWhileStatement() (*ast.WhileStatement, error) {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if err := p.expectPeek(token.LPAREN); err != nil {
		return nil, err
	}

	p.nextToken() // skip LPAREN
	var err error
	stmt.Condition, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.expectPeek(token.RPAREN); err != nil {
		return nil, err
	}
	if err := p.expectPeek(token.LBRACE); err != nil {
		return nil, err
	}
	stmt.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return stmt, nil
}
