package re0

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

const (
	TOKEN_MULTI  rune = '*'
	TOKEN_SINGLE rune = '?'

	MODE_MULTI  int = 3
	MODE_SINGLE int = 2
	MODE_STATIC int = 1
)

type Expression []*Token

type Tokens []*Token

type Token struct {
	buf   []byte
	skip  int
	multi bool
}

func (t *Token) String() string {
	return fmt.Sprintf("%q %d %t", string(t.buf), t.skip, t.multi)
}

func (t *Token) Shard() byte {
	if len(t.buf) > 0 {
		return t.buf[0]
	} else {
		return 0
	}
}

func (t *Token) Equal(t1 *Token) bool {
	if !bytes.Equal(t.buf, t1.buf) {
		return false
	}
	if t.skip != t1.skip {
		return false
	}
	if t.multi != t1.multi {
		return false
	}
	return true
}

func (t *Token) Fuzzy() bool {
	return t.multi || t.skip > 0
}

func (t *Token) Less(t1 *Token) bool {
	return bytes.Compare(t.buf, t1.buf) < 0
}

type parserState struct {
	lastToken *Token
	lastMode  int
	exp       Expression
}

func appendRune(b []byte, r rune) []byte {
	if r < utf8.RuneSelf {
		return append(b, byte(r))
	}

	rb := make([]byte, utf8.UTFMax)
	n := utf8.EncodeRune(rb, r)
	return append(b, rb[0:n]...)
}

func (self *parserState) process(r rune) {
	mode := self.modeByR(r)

	if self.lastToken == nil {
		buf := []byte{}
		self.lastToken = &Token{buf: buf}
	}

	modMode := false
	if mode == MODE_MULTI || mode == MODE_SINGLE {
		modMode = true
	}

	lastModMode := false
	if self.lastMode == MODE_MULTI || self.lastMode == MODE_SINGLE {
		lastModMode = true
	}

	// changed
	if self.lastMode > 0 && modMode && !lastModMode {
		self.exp = append(self.exp, self.lastToken)

		buf := []byte{}
		self.lastToken = &Token{buf: buf}
	}

	// update
	switch r {
	case TOKEN_SINGLE:
		self.lastToken.skip++
	case TOKEN_MULTI:
		self.lastToken.multi = true
	default:
		self.lastToken.buf = appendRune(self.lastToken.buf, r)
	}

	self.lastMode = mode
}

func (self *parserState) modeByR(r rune) int {
	if r == TOKEN_MULTI {
		return MODE_MULTI
	} else if r == TOKEN_SINGLE {
		return MODE_SINGLE
	}
	return MODE_STATIC
}

func (self *parserState) last() {
	// save prev token
	if self.lastToken != nil {
		self.exp = append(self.exp, self.lastToken)
	}
}

func (e *Token) MatchOne(r []byte) (bool, []byte) {
	if e.skip > 0 {
		if len(r) < e.skip {
			return false, nil
		}
		r = r[e.skip:]
	}

	if len(e.buf) == 0 {
		return e.multi, nil // TODO fix for ** case
	}

	if len(r) < len(e.buf) {
		return false, nil
	}

	if e.multi {
		ind := bytes.Index(r, e.buf)
		if ind == -1 {
			return false, nil
		}
		r = r[ind+len(e.buf):]
	} else {
		if !bytes.Equal(r[:len(e.buf)], e.buf) {
			return false, nil
		}
		r = r[len(e.buf):]
	}

	return true, r
}

func (self Expression) Match(r []byte) bool {
	match := false
	for _, e := range self {
		match, r = e.MatchOne(r)
		if len(r) == 0 {
			return match
		}
	}

	return false
}

func Compile(s []byte) Expression {
	state := &parserState{
		exp:      Expression{},
		lastMode: -1,
	}
	reader := bytes.NewReader(s)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		state.process(r)
	}
	state.last()

	return state.exp
}
