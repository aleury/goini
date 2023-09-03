package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int

const (
	itemError itemType = iota // error occurred
	itemEOF                   // end of file has been reached

	itemLeftBracket   // left bracket [
	itemRrightBracket // right bracket ]
	itemNewLine       // new line \n
	itemEqualSign     // equals sign

	itemSection
	itemKey
	itemValue
)

const (
	eof          rune   = -1
	leftBracket  string = "["
	rightBracket string = "]"
	newLine      string = "\n"
	equalSign    string = "="
)

// item represents a token or text string returned from the scanner.
type item struct {
	Type  itemType
	Value string
}

func (i item) String() string {
	switch i.Type {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.Value
	}
	if len(i.Value) > 10 {
		return fmt.Sprintf("%.10q...", i.Value)
	}
	return fmt.Sprintf("%q", i.Value)
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read.
	state stateFn   // current state function
	items chan item // channel of scanned items
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		state: lexStart,
		items: make(chan item, 2),
	}
	return l
}

// run lexes the input by executing the state functions until
// the state is nil.
func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	l.width = w
	return r
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return nil
}

// state functions

func lexStart(l *lexer) stateFn {
loop:
	for {
		if strings.HasPrefix(l.input[l.pos:], leftBracket) {
			return lexLeftBracket
		}
		switch r := l.next(); {
		case r == eof:
			break loop
		case unicode.IsSpace(r):
			l.ignore()
		case unicode.IsLetter(r):
			return lexKey
		}
	}
	l.emit(itemEOF)
	return nil
}

func lexLeftBracket(l *lexer) stateFn {
	l.pos += len(leftBracket)
	l.emit(itemLeftBracket)
	return lexSection
}

func lexSection(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightBracket) {
			l.emit(itemSection)
			return lexRightBracket
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed section")
		case !unicode.IsLetter(r):
			return l.errorf("invalid section")
		}
	}
}

func lexRightBracket(l *lexer) stateFn {
	l.pos += len(rightBracket)
	l.emit(itemRrightBracket)
	return lexStart
}

func lexKey(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], equalSign) {
			l.emit(itemKey)
			return lexEqualSign
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unassigned key")
		case !unicode.IsLetter(r):
			return l.errorf("invalid key")
		}
	}
}

func lexEqualSign(l *lexer) stateFn {
	l.pos += len(equalSign)
	l.emit(itemEqualSign)
	return lexValue
}

func lexValue(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], newLine) {
			l.emit(itemValue)
			return lexStart
		}
		if l.next() == eof {
			return l.errorf("unexpected end of file")
		}
	}
}
