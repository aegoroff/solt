package solution

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

type tokenType int

const (
    itemError tokenType = iota
    itemNIL    // used in the parser to indicate no type
    itemEOF
    itemDatetime
)

const (
    eof          = 0
    comma        = ','
    eq           = '='
    commentStart = '#'
    stringStart  = '"'
    stringEnd    = '"'
    parenOpen    = '('
    parenClose   = ')'
)

type stateFn func(lx *lexer) stateFn

type lexer struct {
    input string
    start int
    pos   int
    line  int
    state stateFn
    items chan yySymType

    // Allow for backing up up to three runes.
    // This is necessary because TOML contains 3-rune tokens (""" and ''').
    prevWidths [3]int
    nprev      int // how many of prevWidths are in use
    // If we emit an eof, we can still back up, but it is not OK to call
    // next again.
    atEOF bool

    // A stack of state functions used to maintain context.
    // The idea is to reuse parts of the state machine in various places.
    // For example, values can appear at the top level or within arbitrarily
    // nested arrays. The last state on the stack is used after a value has
    // been lexed. Similarly for comments.
    stack []stateFn
}

func (lx *lexer) nextItem() yySymType {
    for {
        select {
        case item := <-lx.items:
            return item
        default:
            lx.state = lx.state(lx)
        }
    }
}

func newLexer(input string) *lexer {
    lx := &lexer{
        input: input,
        state: lexTop,
        line:  1,
        items: make(chan yySymType, 10),
        stack: make([]stateFn, 0, 10),
    }
    return lx
}

func (lx *lexer) push(state stateFn) {
    lx.stack = append(lx.stack, state)
}

func (lx *lexer) pop() stateFn {
    if len(lx.stack) == 0 {
        return lx.errorf("BUG in lexer: no states to pop")
    }
    last := lx.stack[len(lx.stack)-1]
    lx.stack = lx.stack[0 : len(lx.stack)-1]
    return last
}

func (lx *lexer) current() string {
    return lx.input[lx.start:lx.pos]
}

func (lx *lexer) emit(typ tokenType) {
    lx.items <- yySymType{tok: typ, str: lx.current(), line: lx.line}
    lx.start = lx.pos
}

func (lx *lexer) emitTrim(typ tokenType) {
    lx.items <- yySymType{tok: typ, str: strings.TrimSpace(lx.current()), line: lx.line}
    lx.start = lx.pos
}

func (lx *lexer) next() (r rune) {
    if lx.atEOF {
        panic("next called after EOF")
    }
    if lx.pos >= len(lx.input) {
        lx.atEOF = true
        return eof
    }

    if lx.input[lx.pos] == '\n' {
        lx.line++
    }
    lx.prevWidths[2] = lx.prevWidths[1]
    lx.prevWidths[1] = lx.prevWidths[0]
    if lx.nprev < 3 {
        lx.nprev++
    }
    r, w := utf8.DecodeRuneInString(lx.input[lx.pos:])
    lx.prevWidths[0] = w
    lx.pos += w
    return r
}

// ignore skips over the pending input before this point.
func (lx *lexer) ignore() {
    lx.start = lx.pos
}

// backup steps back one rune. Can be called only twice between calls to next.
func (lx *lexer) backup() {
    if lx.atEOF {
        lx.atEOF = false
        return
    }
    if lx.nprev < 1 {
        panic("backed up too far")
    }
    w := lx.prevWidths[0]
    lx.prevWidths[0] = lx.prevWidths[1]
    lx.prevWidths[1] = lx.prevWidths[2]
    lx.nprev--
    lx.pos -= w
    if lx.pos < len(lx.input) && lx.input[lx.pos] == '\n' {
        lx.line--
    }
}

// accept consumes the next rune if it's equal to `valid`.
func (lx *lexer) accept(valid rune) bool {
    if lx.next() == valid {
        return true
    }
    lx.backup()
    return false
}

// peek returns but does not consume the next rune in the input.
func (lx *lexer) peek() rune {
    r := lx.next()
    lx.backup()
    return r
}

// skip ignores all input that matches the given predicate.
func (lx *lexer) skip(pred func(rune) bool) {
    for {
        r := lx.next()
        if pred(r) {
            continue
        }
        lx.backup()
        lx.ignore()
        return
    }
}

// errorf stops all lexing by emitting an error and returning `nil`.
// Note that any value that is a character is escaped if it's a special
// character (newlines, tabs, etc.).
func (lx *lexer) errorf(format string, values ...interface{}) stateFn {
    lx.items <- yySymType{
        tok:  itemError,
        str:  fmt.Sprintf(format, values...),
        line: lx.line,
    }
    return nil
}

// lexTop consumes elements at the top level of TOML data.
func lexTop(lx *lexer) stateFn {
    r := lx.next()
    if isWhitespace(r) || isNL(r) {
        return lexSkip(lx, lexTop)
    }

    switch {
    case r == eof:
        if lx.pos > lx.start {
            return lx.errorf("unexpected EOF")
        }
        lx.emit(eof)
        return nil
    }

    // At this point, the only valid item can be an identifier, so we back up
    // and let the key lexer do the rest.
    lx.backup()
    lx.push(lexTopEnd)
    return lexIdentifier
}

// lexTopEnd is entered whenever a top-level item has been consumed. (A value
// or a table.) It must see only whitespace, and will turn back to lexTop
// upon a newline. If it sees EOF, it will quit the lexer successfully.
func lexTopEnd(lx *lexer) stateFn {
    r := lx.next()
    switch {
    case r == commentStart:
        // a comment will read to a newline for us.
        lx.push(lexTop)
        return lexCommentStart
    case isWhitespace(r):
        return lexTopEnd
    case isNL(r):
        lx.ignore()
        return lexTop
    case r == eof:
        lx.emit(itemEOF)
        return nil
    }
    return lx.errorf("expected a top-level item to end with a newline, "+
        "comment, or EOF, but got %q instead", r)
}

func lexIdentifier(lx *lexer) stateFn {
    switch r := lx.next(); {
    case isIdentifierChar(r):
        return lexIdentifier
    default:
        lx.backup()
        lx.emit(IDENTIFIER)
        return lexOther
    }
}

func lexOther(lx *lexer) stateFn {
    switch r := lx.next(); {
    case isIdentifierChar(r):
        return lexIdentifier
    case r == comma:
        lx.ignore()
        lx.emit(COMMA)
        return lexSkip(lx, lexOther)
    case r == eq:
        lx.ignore()
        lx.emit(EQ)
        return lexOther
    case isDigit(r):
        lx.push(lexOther)
        lx.backup() // avoid an extra state and use the same as above
        return lexNumberOrDateStart
    case isWhitespace(r):
        return lexSkip(lx, lexOther)
    case r == parenOpen:
        lx.ignore()
        lx.emit(PAREN_OPEN)
        return lexIdentifierOrString
    case r == stringStart:
        lx.push(lexOther)
        lx.backup()
        return lexValue
    case r == parenClose:
        lx.ignore()
        lx.emit(PAREN_CLOSE)
        return lexOther
    case r == '\r':
        return lexSkip(lx, lexOther)
    case isCrlf(r):
        lx.ignore()
        lx.emit(CRLF)
        return lexNewLine
    default:
        return lx.errorf("identifiers cannot contain %q", r)
    }
}

func lexIdentifierOrString(lx *lexer) stateFn {
    switch r := lx.peek(); {
    case isIdentifierChar(r):
        return lexIdentifier
    case r == stringStart:
        lx.push(lexOther)
        return lexValue
    default:
        return lx.errorf("identifier or string cannot cannot start from %q", r)
    }
}

func lexNewLine(lx *lexer) stateFn {
    switch r := lx.next(); {
    case isIdentifierChar(r):
        lx.backup()
        return lexIdentifier
    case r == commentStart:
        lx.backup()
        // a comment will read to a newline for us.
        lx.push(lexOther)
        return lexCommentStart
    case r == '\t':
        lx.ignore()
        return lexBareStringStart
    case r == eof:
        lx.emit(itemEOF)
        return nil
    default:
        return lx.errorf("line cannot start from %q", r)
    }
}

func lexBareStringStart(lx *lexer) stateFn {
    switch r := lx.next(); {
    case r == '\t':
        lx.ignore()
        return lexBareString
    case r == ' ':
        lx.ignore()
        return lexBareString
    case isIdentifierChar(r):
        return lexIdentifier
    default:
        return lx.errorf("bare string cannot start from %q", r)
    }
}

func lexBareString(lx *lexer) stateFn {
    switch r := lx.next(); {
    case isCrlf(r):
        lx.backup()
        lx.emit(BARE_STRING)
        return lexBareStringEnd
    case r == eq:
        lx.backup()
        lx.emit(BARE_STRING)
        return lexBareStringEnd
    case r == '\r':
        lx.backup()
        lx.emit(BARE_STRING)
        return lexBareStringEnd
    default:
        return lexBareString
    }
}

func lexBareStringEnd(lx *lexer) stateFn {
    switch r := lx.next(); {
    case r == '\r':
        lx.ignore()
        return lexBareStringEnd
    case r == eq:
        lx.ignore()
        lx.emit(EQ)
        return lexBareStringStart
    case isCrlf(r):
        lx.ignore()
        lx.emit(CRLF)
        return lexNewLine
    default:
        return lexBareString
    }
}

// lexValue starts the consumption of a value anywhere a value is expected.
// lexValue will ignore whitespace.
// After a value is lexed, the last state on the next is popped and returned.
func lexValue(lx *lexer) stateFn {
    // We allow whitespace to precede a value, but NOT newlines.
    // In array syntax, the array states are responsible for ignoring newlines.
    r := lx.next()
    switch {
    case isWhitespace(r):
        return lexSkip(lx, lexValue)
    case isDigit(r):
        lx.backup() // avoid an extra state and use the same as above
        return lexNumberOrDateStart
    }
    switch r {
    case stringStart:
        lx.ignore() // ignore the '"'
        return lexString
    case '+', '-':
        return lexNumberStart
    case '.': // special error case, be kind to users
        return lx.errorf("floats must start with a digit, not '.'")
    }
    return lx.errorf("expected value but found %q instead", r)
}

// lexString consumes the inner contents of a string. It assumes that the
// beginning '"' has already been consumed and ignored.
func lexString(lx *lexer) stateFn {
    r := lx.next()
    switch {
    case r == eof:
        return lx.errorf("unexpected EOF")
    case isNL(r):
        return lx.errorf("strings cannot contain newlines")
        //case r == '\\':
        //    lx.push(lexString)
        //    return lexStringEscape
    case r == stringEnd:
        lx.backup()
        lx.emit(STRING)
        lx.next()
        lx.ignore()
        return lx.pop()
    }
    return lexString
}

// lexNumberOrDateStart consumes either an integer, a float, or datetime.
func lexNumberOrDateStart(lx *lexer) stateFn {
    r := lx.next()
    if isDigit(r) {
        return lexNumberOrDate
    }
    switch r {
    case '_':
        return lexNumber
    case 'e', 'E':
        return lexFloat
    case '.':
        return lx.errorf("floats must start with a digit, not '.'")
    }
    return lx.errorf("expected a digit but got %q", r)
}

// lexNumberOrDate consumes either an integer, float or datetime.
func lexNumberOrDate(lx *lexer) stateFn {
    r := lx.next()
    if isDigit(r) {
        return lexNumberOrDate
    }
    switch r {
    case '-':
        return lexDatetime
    case '_':
        return lexNumber
    case '.', 'e', 'E':
        return lexFloat
    }

    lx.backup()
    lx.emit(NUMBER)
    return lx.pop()
}

// lexDatetime consumes a Datetime, to a first approximation.
// The parser validates that it matches one of the accepted formats.
func lexDatetime(lx *lexer) stateFn {
    r := lx.next()
    if isDigit(r) {
        return lexDatetime
    }
    switch r {
    case '-', 'T', ':', '.', 'Z', '+':
        return lexDatetime
    }

    lx.backup()
    lx.emit(itemDatetime)
    return lx.pop()
}

// lexNumberStart consumes either an integer or a float. It assumes that a sign
// has already been read, but that *no* digits have been consumed.
// lexNumberStart will move to the appropriate integer or float states.
func lexNumberStart(lx *lexer) stateFn {
    // We MUST see a digit. Even floats have to start with a digit.
    r := lx.next()
    if !isDigit(r) {
        if r == '.' {
            return lx.errorf("floats must start with a digit, not '.'")
        }
        return lx.errorf("expected a digit but got %q", r)
    }
    return lexNumber
}

// lexNumber consumes an integer or a float after seeing the first digit.
func lexNumber(lx *lexer) stateFn {
    r := lx.next()
    if isDigit(r) {
        return lexNumber
    }
    switch r {
    case '_':
        return lexNumber
    case '.', 'e', 'E':
        return lexFloat
    }

    lx.backup()
    lx.emit(NUMBER)
    return lx.pop()
}

// lexFloat consumes the elements of a float. It allows any sequence of
// float-like characters, so floats emitted by the lexer are only a first
// approximation and must be validated by the parser.
func lexFloat(lx *lexer) stateFn {
    r := lx.next()
    if isDigit(r) {
        return lexFloat
    }
    switch r {
    case '_', '.', '-', '+', 'e', 'E':
        return lexFloat
    }

    lx.backup()
    lx.emit(NUMBER)
    return lx.pop()
}

// lexCommentStart begins the lexing of a comment. It will emit
// itemCommentStart and consume no characters, passing control to lexComment.
func lexCommentStart(lx *lexer) stateFn {
    lx.ignore()
    return lexComment
}

// lexComment lexes an entire comment. It assumes that '#' has been consumed.
// It will consume *up to* the first newline character, and pass control
// back to the last state on the stack.
func lexComment(lx *lexer) stateFn {
    r := lx.peek()
    if isNL(r) || r == eof {
        lx.emit(COMMENT)
        return lx.pop()
    }
    lx.next()
    return lexComment
}

// lexSkip ignores all slurped input and moves on to the next state.
func lexSkip(lx *lexer, nextState stateFn) stateFn {
    return func(lx *lexer) stateFn {
        lx.ignore()
        return nextState
    }
}

// isWhitespace returns true if `r` is a whitespace character according
// to the spec.
func isWhitespace(r rune) bool {
    return r == '\t' || r == ' '
}

func isNL(r rune) bool {
    return r == '\n' || r == '\r'
}

func isCrlf(r rune) bool {
    return r == '\n'
}

func isDigit(r rune) bool {
    return r >= '0' && r <= '9'
}

func isHexadecimal(r rune) bool {
    return (r >= '0' && r <= '9') ||
        (r >= 'a' && r <= 'f') ||
        (r >= 'A' && r <= 'F')
}

func isIdentifierChar(r rune) bool {
    return (r >= 'A' && r <= 'Z') ||
        (r >= 'a' && r <= 'z')
}

func (itype tokenType) String() string {
    switch itype {
    case CRLF:
        return "CRLF"
    case IDENTIFIER:
        return "IDENTIFIER"
    case NUMBER:
        return "NUMBER"
    case COMMA:
        return "COMMA"
    case COMMENT:
        return "COMMENT"
    case STRING:
        return "STRING"
    case BARE_STRING:
        return "BARE_STRING"
    case PAREN_OPEN:
        return "PAREN_OPEN"
    case PAREN_CLOSE:
        return "PAREN_CLOSE"
    case EQ:
        return "EQ"
    case itemError:
        return "Error"
    case itemEOF:
        return "EOF"
    case itemNIL:
        return "NIL"
    case itemDatetime:
        return "DateTime"
    }
    panic(fmt.Sprintf("BUG: Unknown type '%d'.", int(itype)))
}

func (item yySymType) String() string {
    return fmt.Sprintf("(%s, %s)", item.tok.String(), item.str)
}
