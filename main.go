package main

import (
  "fmt"
  "unicode"
  "os"
  "bufio"
  "io"
  "log"
)

type Type int

// Tokens
const (
  TOKEN_ERROR Type = iota
  TOKEN_RPAREN // (
  TOKEN_LPAREN // )
  TOKEN_RBRACE // {
  TOKEN_LBRACE // }
  TOKEN_IDENT // abc
  TOKEN_STRING // "hello"
  TOKEN_CHAR // 'c'
  TOKEN_INTEGER // 123
  TOKEN_FLOAT // 123.123e23
  TOKEN_OPERATOR // +
  TOKEN_KEYWORD // var
  TOKEN_TYPE // string
)

// Operators
const (
  _OPERATOR_BINARY_BEGIN = iota
  // Mathematical operators
  OPERATOR_ADD // +
  OPERATOR_SUB // -
  OPERATOR_MUL // *
  OPERATOR_DIV // /
  OPERATOR_MOD // %
  // Method or member resolution
  OPERATOR_DOT // .
  // Boolean logical operators
  OPERATOR_EQ // ==
  OPERATOR_NEQ // !=
  OPERATOR_GT // >
  OPERATOR_GTE // >=
  OPERATOR_LT // <
  OPERATOR_LTE // <=
  _OPERATOR_BINARY_END
  _OPERATOR_UNARY_BEGIN
  // Logical Not Operator
  OPERATOR_NOT // !
  _OPERATOR_UNARY_END
)

// maps operators to strings
var operators = []string {
  OPERATOR_ADD: "+",
  OPERATOR_SUB: "-",
  OPERATOR_MUL: "*",
  OPERATOR_DIV: "/",
  OPERATOR_MOD: "%",
  OPERATOR_DOT: ".",
  OPERATOR_EQ: "==",
  OPERATOR_NEQ: "!=",
  OPERATOR_GT: ">",
  OPERATOR_GTE: ">=",
  OPERATOR_LT: "<",
  OPERATOR_LTE: "<=",
  OPERATOR_NOT: "!",
}

// Types
const (
  TYPE_I32 = iota // i32
  TYPE_F32 // f32
  TYPE_BOOL // bool
  TYPE_STRING // string
  TYPE_RUNE // rune
  TYPE_TUPLE // tuple(i32, bool)
)

var types = []string {
  TYPE_I32: "i32",
  TYPE_F32: "f32",
  TYPE_BOOL: "bool",
  TYPE_STRING: "string",
  TYPE_RUNE: "rune",
  TYPE_TUPLE: "tuple",
}

// Keyword
const (
  KEYWORD_TYPE = iota // :
  KEYWORD_ASSIGN // =
  KEYWORD_COMMA // ,
  KEYWORD_SEMICOLON // ;
)

var keywords = []string {
  KEYWORD_TYPE: ":",
  KEYWORD_ASSIGN: "=",
  KEYWORD_COMMA: ",",
  KEYWORD_SEMICOLON: ";",
}

type Token struct {
  Type Type
  Value string
}

type Lexer struct {
  Input chan rune
  lexed []rune
  back []rune
  toks chan *Token
  parenDepth int
  funcDepth []int
}

type RunePredicate func(rune)bool

func (l *Lexer) next() rune {
  if len(l.back) > 0 {
    c := l.back[len(l.back)-1]
    l.back = l.back[:len(l.back)-1]
    return c
  }
  return <-l.Input
}

func (l *Lexer) Next() rune {
  c := l.next()
  l.lexed = append(l.lexed, c)
  return c
}

func (l *Lexer) Ignore() rune {
  return l.next()
}

func (l *Lexer) Back() {
  if len(l.lexed) > 0 {
    l.back = append(l.back, l.lexed[len(l.lexed)-1])
    l.lexed = l.lexed[:len(l.lexed)-1]
  }
}

func (l *Lexer) NextUpTo(pred RunePredicate) rune {
  for {
    c := l.Next()
    if !pred(c) {
      return c
    }
  }
}

func (l *Lexer) IgnoreUpTo(pred RunePredicate) rune {
  for {
    c := l.Ignore()
    if c == 0 {
      return c
    }
    if !pred(c) {
      l.lexed = append(l.lexed, c) // remember
      return c
    }
  }
}

func (l *Lexer) Emit(t Type) {
  l.toks <- &Token{ Type: t, Value: l.Lexed() }
  l.lexed = []rune{}
}

func (l *Lexer) Lexed() string {
  s := ""
  for _, r := range l.lexed {
    s += string(r)
  }
  return s
}

type StateFunc func(*Lexer)StateFunc

func (l *Lexer) TokenGenerator() chan *Token {
  l.toks = make(chan *Token)
  go func(l *Lexer) {
    // state machine
    f := startState
    for f != nil {
      f = f(l)
    }
    close(l.toks)
    fmt.Printf("Lexed: '%s'\n", l.Lexed())
  }(l)
  return l.toks
}

func whitespacePredicate(c rune) bool {
  return unicode.IsSpace(c)
}

// Identifier after first character (first cannot be number)
func identifierPredicate(c rune) bool {
  return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_'
}

// The operator charset is defined by
// the characters currently used in operators.
func operatorPredicate(c rune) bool {
  for _, op := range operators {
    for _, d := range op {
      if c == d {
        return true
      }
    }
  }
  return false
}

func stringPredicate(c rune) bool {
  return unicode.IsPrint(c) && c != '"'
}

func numberPredicate(c rune) bool {
  return unicode.IsNumber(c) || c == '_'
}

func charPredicate(c rune) bool {
  return unicode.IsPrint(c) && c != '\''
}

func identifierCheck(c rune) bool {
  return unicode.IsLetter(c) || c == '_'
}

func whitespaceCheck(c rune) bool {
  return unicode.IsSpace(c)
}

func stringCheck(c rune) bool {
  return c == '"'
}

func numberCheck(c rune) bool {
  return unicode.IsNumber(c)
}

func charCheck(c rune) bool {
  return c == '\''
}

func semicolonState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == ';' {
    l.Emit(TOKEN_KEYWORD)
    // end of expression
    return startState
  }
  fmt.Printf("Expected semicolon, found '%c'\n", c)
  return nil
}

func binaryOperatorState(l *Lexer) StateFunc {
  l.IgnoreUpTo(whitespacePredicate)
  l.Back()
  l.NextUpTo(operatorPredicate)
  l.Back()
  // Check if an operator
  for j := _OPERATOR_BINARY_BEGIN+1; j < _OPERATOR_BINARY_END; j++ {
    if l.Lexed() == operators[j] {
      l.Emit(TOKEN_OPERATOR)
      return expressionState
    }
  }
  // Not an operator
  return semicolonState
}

func unaryOperatorState(l *Lexer) StateFunc {
  l.IgnoreUpTo(whitespacePredicate)
  l.Back()
  l.NextUpTo(operatorPredicate)
  l.Back()
  // Check if an operator
  for j := _OPERATOR_UNARY_BEGIN+1; j < _OPERATOR_UNARY_END; j++ {
    if l.Lexed() == operators[j] {
      l.Emit(TOKEN_OPERATOR)
      return expressionState
    }
  }
  // Not an operator
  fmt.Printf("Unexpected runes '%s' in expression\n", l.Lexed())
  return nil
}

func floatState(l *Lexer) StateFunc {
  c := l.NextUpTo(numberPredicate)
  if c == 'e' && len(l.lexed) > 1 {
    l.NextUpTo(numberPredicate)
    l.Back()
    if len(l.lexed) > 0 {
      l.Emit(TOKEN_FLOAT)
      return expressionPostLiteralState
    }
  }
  fmt.Printf("Not a proper float '%s'\n", l.Lexed())
  return nil
}

func numberState(l *Lexer) StateFunc {
  c := l.NextUpTo(numberPredicate)
  if c == '.' {
    return floatState
  }
  l.Back()
  l.Emit(TOKEN_INTEGER)
  return expressionPostLiteralState
}

func funcCallState(l *Lexer) StateFunc {
  l.Emit(TOKEN_LPAREN)
  // function call stack,
  // needed because grammars for parenthesized
  // expressions and function calls differ
  l.funcDepth = append(l.funcDepth, l.parenDepth)
  l.parenDepth = 0
  return expressionState
}

func identifierState(l *Lexer) StateFunc {
  l.NextUpTo(identifierPredicate)
  l.Emit(TOKEN_IDENT)
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == '(' {
    return funcCallState
  }
  l.Back()
  return expressionPostLiteralState
}

// After lexing an identifier or literal
func expressionPostLiteralState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == ')' {
    return expressionRightParenState
  } else if c == ',' {
    return expressionCommaState
  } else {
    l.Back()
    return binaryOperatorState
  }
}

func expressionRightParenState(l *Lexer) StateFunc {
  if l.parenDepth > 0 {
    l.Emit(TOKEN_RPAREN)
    l.parenDepth -= 1
    return binaryOperatorState
  } else if len(l.funcDepth) > 0 {
    l.Emit(TOKEN_RPAREN)
    l.parenDepth = l.funcDepth[len(l.funcDepth)-1]
    l.funcDepth = l.funcDepth[:len(l.funcDepth)-1]
    return binaryOperatorState
  }
  fmt.Printf("Extraneous right parenthesis '%s'\n", l.Lexed())
  return nil
}

func expressionCommaState(l *Lexer) StateFunc {
  // Commas are found in function calls only
  if l.parenDepth == 0 && len(l.funcDepth) > 0 {
    // Inside a function, not nested in parenthesized expression.
    l.Emit(TOKEN_KEYWORD)
    return expressionState
  }
  fmt.Printf("Commas are only used as separators in function calls '%s'\n", l.Lexed())
  return nil
}

func expressionState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  // Check for identifiers or type literals
  if identifierCheck(c) {
    return identifierState
  } else if numberCheck(c) {
    return numberState
  } else if stringCheck(c) {
    l.NextUpTo(stringPredicate)
    l.Emit(TOKEN_STRING)
    return binaryOperatorState
  } else if charCheck(c) {
    l.NextUpTo(charPredicate)
    l.Emit(TOKEN_CHAR)
    return binaryOperatorState
  } else if c == '(' {
    // Parenthesized expression
    l.Emit(TOKEN_LPAREN)
    l.parenDepth += 1
    return expressionState
  } else if c == ')' {
    return expressionRightParenState
  } else if c == ',' {
    return expressionCommaState
  }
  // maybe unary operator
  return unaryOperatorState
}

func letAssignmentState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == '=' {
    l.Emit(TOKEN_KEYWORD)
    return expressionState
  }
  fmt.Printf("Expected assignment, found unexpected rune '%c'\n", c)
  return nil
}

// First time inside a tuple, lexing type subexpression
func letTypeTupleState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == '(' {
    l.parenDepth += 1
    l.Emit(TOKEN_LPAREN)
    return letTypeState
  }
  fmt.Printf("Expected parenthesized subtype expression after tuple, found unexpected run '%v'\n", c)
  return nil
}

// letTypeState has discovered that it is lexing a type nested in a tuple.
func letInsideTupleState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == ')' {
    l.parenDepth -= 1
    l.Emit(TOKEN_RPAREN)
    if l.parenDepth > 0 {
      // still in tuple, expect comma or closing paren
      return letInsideTupleState
    } else {
      return letAssignmentState
    }
  } else if c == ',' {
    l.Emit(TOKEN_KEYWORD)
    // need to lex another type in tuple definition
    return letTypeState
  }
  fmt.Printf("Unexpected rune inside tuple '%v'\n", c)
  return nil
}

// letIdentState has lexed an identifier and found ':'.
func letTypeState(l *Lexer) StateFunc {
  // For now we only allow simple type annotation
  // like 'i32', 'bool', 'string', etc.
  c := l.IgnoreUpTo(whitespacePredicate)
  if identifierCheck(c) {
    l.NextUpTo(identifierPredicate)
    l.Back()
    for _, t := range types {
      if t == l.Lexed() {
        l.Emit(TOKEN_TYPE)
        if t == types[TYPE_TUPLE] {
          // tuples are nested types tuple(i32, bool)
          return letTypeTupleState
        }
        if l.parenDepth > 0 {
          return letInsideTupleState
        }
        return letAssignmentState
      }
    }
  }
  fmt.Printf("Expected identifier in type expression, found unexpected rune '%c'", c)
  return nil
}

// keywordState has lexed the 'let' keyword,
// we expect to find an identifier
func letIdentState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate) // kill whitespace
  if identifierCheck(c) {
    l.NextUpTo(identifierPredicate)
    l.Back()
    l.Emit(TOKEN_IDENT)
    c := l.IgnoreUpTo(whitespacePredicate)
    if c == ':' {
      l.Emit(TOKEN_KEYWORD)
      return letTypeState
    } else {
      l.Back()
      return letAssignmentState
    }
    // fallthrough and error
  }
  fmt.Printf("Unexpected rune '%v'\n", c)
  return nil
}

// startState has lexed the first character
// of what we expect to be a keyword.
func keywordState(l *Lexer) StateFunc {
  l.NextUpTo(identifierPredicate)
  l.Back()
  if l.Lexed() == "let" {
    l.Emit(TOKEN_KEYWORD)
    return letIdentState
  }
  fmt.Printf("Unexpected keyword '%s'\n", l.Lexed())
  return nil
}

func startState(l *Lexer) StateFunc {
  // get rid of whitespace
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == 0 {
    return nil // end of input (is ok here)
  }
  if identifierCheck(c) {
    return keywordState
  }
  return expressionState
}

func main() {
  r := make(chan rune)
  l := Lexer { Input: r }
  c := l.TokenGenerator()
  go func() {
    reader := bufio.NewReader(os.Stdin)
    for {
      fmt.Print(">> ")
      text, err := reader.ReadString('\n')
      if err != nil {
        if err == io.EOF {
          close(r)
          return
        }
        log.Fatal(err)
      }
      for _, c := range text {
        r <- c
      }
    }
  }()
  for ch := range c {
    fmt.Println(ch)
  }
}
