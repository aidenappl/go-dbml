// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/types/tokens.ts, src/core/lexer/utils.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

// Package lexer turns DBML source text into a token stream.
//
// The lexer is a faithful port of @dbml/parse's lexer: whitespace and comments
// become trivia attached to neighbouring tokens, invalid tokens are attached
// to the previous valid token, and all positions are tracked in UTF-16 code
// units so they agree with JavaScript (and Monaco) offsets.
package lexer

// Kind classifies a token.
type Kind uint8

// Token kinds.
const (
	Space Kind = iota
	Tab
	Newline
	Comma
	LParen
	RParen
	LBrace
	RBrace
	LBracket
	RBracket
	Op
	EOF
	NumericLiteral
	StringLiteral
	ColorLiteral
	FunctionExpression
	QuotedString
	Identifier
	Semicolon
	Colon
	SingleLineComment
	MultilineComment
	Wildcard
)

var kindNames = [...]string{
	Space: "<space>", Tab: "<tab>", Newline: "<newline>", Comma: "<comma>", LParen: "<lparen>", RParen: "<rparen>",
	LBrace: "<lbrace>", RBrace: "<rbrace>", LBracket: "<lbracket>", RBracket: "<rbracket>", Op: "<op>", EOF: "<eof>",
	NumericLiteral: "<number>", StringLiteral: "<string>", ColorLiteral: "<color>",
	FunctionExpression: "<function-expression>", QuotedString: "<variable>", Identifier: "<identifier>",
	Semicolon: "<semicolon>", Colon: "<colon>", SingleLineComment: "<single-line-comment>",
	MultilineComment: "<multiline-comment>", Wildcard: "<wildcard>",
}

// String returns the @dbml/parse name of the kind (e.g. "<identifier>").
func (k Kind) String() string {
	if int(k) >= 0 && int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "<unknown>"
}

// Position is a 0-based source position. Offset and Column are in UTF-16 code
// units (JavaScript semantics); Byte is the byte offset into the Go string.
// Fields are int32 to keep tokens and syntax nodes small (inputs are
// capped far below 2 GiB).
type Position struct {
	Offset int32
	Line   int32
	Column int32
	Byte   int32
}

// Token is a lexed token with its attached trivia.
//
// Trivia (whitespace and comments) is not stored as tokens: a token keeps
// the source range of its leading and trailing trivia, and LeadingTrivia /
// TrailingTrivia re-lex that range on demand. This keeps a token at about
// 100 bytes, with no allocation per space or newline.
type Token struct {
	Kind    Kind
	Invalid bool
	flags   uint8
	Value   string
	Start   Position
	End     Position

	leadStart   Position // start of the leading trivia (leadLen > 0)
	leadLen     int32    // bytes of leading trivia, ending at or before Start
	trailLen    int32    // bytes from End to the end of the trailing trivia
	trailEndOff int32    // UTF-16 offset after the trailing trivia
	src         *string
	inv         *invalidTokens
}

type invalidTokens struct{ leading, trailing []*Token }

const (
	flagTrailNewline uint8 = 1 << iota
	flagTrailSpace
	flagLeadNewline
)

// LeadingTrivia returns the whitespace/comment tokens before the token
// (only when it starts a line), re-lexed from the source.
func (t *Token) LeadingTrivia() []*Token {
	if t == nil || t.leadLen == 0 {
		return nil
	}
	return lexTrivia(*t.src, t.leadStart, t.leadStart.Byte+t.leadLen)
}

// TrailingTrivia returns the whitespace/comment tokens after the token up
// to and including the end of its line, re-lexed from the source.
func (t *Token) TrailingTrivia() []*Token {
	if t == nil || t.trailLen == 0 {
		return nil
	}
	return lexTrivia(*t.src, t.End, t.End.Byte+t.trailLen)
}

// LeadingInvalid returns the invalid tokens attached before the token.
func (t *Token) LeadingInvalid() []*Token {
	if t == nil || t.inv == nil {
		return nil
	}
	return t.inv.leading
}

// TrailingInvalid returns the invalid tokens attached after the token.
func (t *Token) TrailingInvalid() []*Token {
	if t == nil || t.inv == nil {
		return nil
	}
	return t.inv.trailing
}

// PrependLeadingInvalid attaches invalid tokens before the token.
func (t *Token) PrependLeadingInvalid(ts []*Token) {
	if len(ts) == 0 {
		return
	}
	if t.inv == nil {
		t.inv = &invalidTokens{}
	}
	t.inv.leading = append(append([]*Token(nil), ts...), t.inv.leading...)
}

// AppendTrailingInvalid attaches an invalid token after the token.
func (t *Token) AppendTrailingInvalid(x *Token) {
	if t.inv == nil {
		t.inv = &invalidTokens{}
	}
	t.inv.trailing = append(t.inv.trailing, x)
}

// IsTrivia reports whether the token is whitespace or a comment.
func (t *Token) IsTrivia() bool {
	switch t.Kind {
	case Newline, Space, Tab, SingleLineComment, MultilineComment:
		return true
	}
	return false
}

// HasTrailingNewline reports whether a newline follows the token on its line.
func (t *Token) HasTrailingNewline() bool { return t.flags&flagTrailNewline != 0 }

// HasTrailingSpaces reports whether a space or tab directly follows the token.
func (t *Token) HasTrailingSpaces() bool { return t.flags&flagTrailSpace != 0 }

// HasLeadingNewline reports whether a newline precedes the token.
func (t *Token) HasLeadingNewline() bool { return t.flags&flagLeadNewline != 0 }

// IsAtStartOfLine reports whether token begins a new line given the previous token.
func IsAtStartOfLine(previous, token *Token) bool {
	return token.HasLeadingNewline() || (previous != nil && previous.HasTrailingNewline())
}

// FullStart is the start offset including leading trivia (UTF-16 units).
func (t *Token) FullStart() int {
	if t.leadLen == 0 {
		return int(t.Start.Offset)
	}
	return int(t.leadStart.Offset)
}

// FullEnd is the end offset including trailing trivia (UTF-16 units).
func (t *Token) FullEnd() int {
	if t.trailLen == 0 {
		return int(t.End.Offset)
	}
	return int(t.trailEndOff)
}

// IsOpChar reports whether c can start an operator token.
func IsOpChar(c rune) bool {
	switch c {
	case '+', '-', '/', '%', '<', '>', '=', '!', '?', '.', '&', '|', '~':
		return true
	}
	return false
}

// Span reports the token's range (nil-safe: a nil token has no span). It
// lets tokens and syntax nodes be handled uniformly by the ast package.
func (t *Token) Span() (start, end Position, fullStart, fullEnd int, ok bool) {
	if t == nil {
		return
	}
	return t.Start, t.End, t.FullStart(), t.FullEnd(), true
}
