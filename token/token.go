package token

import (
  "errors"
  "fmt"
)

type TokenType int

// Tokens
const (
  ERROR TokenType = iota
  LPAREN // (
  RPAREN // )
  RBRACE // {
  LBRACE // }
  IDENT // abc
  FUNC // abc() Ident before call operator
  STRING // "hello"
  CHAR // 'c'
  INTEGER // 123
  FLOAT // 123.123e23
  OPERATOR // +
  KEYWORD // var
  TYPE // string
)

type GenericToken struct {
  Typ TokenType
  Val string
}

func (t *GenericToken) Value() string {
  return t.Val
}

func (t *GenericToken) Type() TokenType {
  return t.Typ
}

// Operators
type Operator int

const (
  _OPERATOR_BINARY_BEGIN Operator = iota

  // Mathematical operators
  OPERATOR_ADD // +
  OPERATOR_SUB // -
  OPERATOR_MUL // *
  OPERATOR_DIV // /
  OPERATOR_MOD // %

  // Tuple creation appender
  OPERATOR_COMMA // ,

  // Method or member resolution
  OPERATOR_DOT // .

  // Type cast operator
  OPERATOR_TYPE // :

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
  OPERATOR_TYPE: ":",
  OPERATOR_EQ: "==",
  OPERATOR_NEQ: "!=",
  OPERATOR_GT: ">",
  OPERATOR_GTE: ">=",
  OPERATOR_LT: "<",
  OPERATOR_LTE: "<=",
  OPERATOR_NOT: "!",
}

type OperatorToken struct {
  Operator Operator
}

func NewBinaryOperatorToken(l string) (*OperatorToken, error) {
  for i := _OPERATOR_BINARY_BEGIN; i < _OPERATOR_BINARY_END; i++ {
    if l == operators[i] {
      return &OperatorToken{ Operator: i }, nil
    }
  }
  return nil, errors.New(fmt.Sprintf("'%s' is not a binary operator", l))
}

func NewUnaryOperatorToken(l string) (*OperatorToken, error) {
  for i := _OPERATOR_UNARY_BEGIN; i < _OPERATOR_UNARY_END; i++ {
    if l == operators[i] {
      return &OperatorToken{ Operator: i }, nil
    }
  }
  return nil, errors.New(fmt.Sprintf("'%s' is not a unary operator", l))
}

func (t *OperatorToken) Value() string {
  return operators[t.Operator]
}

func (t *OperatorToken) Type() TokenType {
  return OPERATOR
}

type OperatorAssoc int

const (
  OPERATOR_ASSOC_LEFT OperatorAssoc = iota
  OPERATOR_ASSOC_RIGHT
)

func (t OperatorToken) Assoc() OperatorAssoc {
  if t.Operator == OPERATOR_NOT {
    return OPERATOR_ASSOC_RIGHT
  }
  return OPERATOR_ASSOC_LEFT
}

func (t OperatorToken) Prec() int {
  switch t.Operator {
  case OPERATOR_DOT:
    return 1
  case OPERATOR_TYPE:
    return 2
  case OPERATOR_ADD:
    return 3
  case OPERATOR_SUB:
    return 3
  case OPERATOR_MUL:
    return 4
  case OPERATOR_DIV:
    return 4
  case OPERATOR_MOD:
    return 4
  }
  return 15
}

// Type
type Type int

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

type TypeToken struct {
  Typ Type
}

func NewTypeToken(l string) (*TypeToken, error) {
  for i, s := range types {
    if l == s {
      return &TypeToken{ Typ: Type(i) }, nil
    }
  }
  return nil, errors.New(fmt.Sprintf("'%s' is not a type", l))
}

func (t *TypeToken) Value() string {
  return types[t.Typ]
}

func (t *TypeToken) Type() TokenType {
  return TYPE
}

// Keyword
type Keyword int

const (
  KEYWORD_TYPE Keyword = iota // :
  KEYWORD_ASSIGN // =
  KEYWORD_COMMA // ,
  KEYWORD_SEMICOLON // ;
  KEYWORD_LET // let
)

var keywords = []string {
  KEYWORD_TYPE: ":",
  KEYWORD_ASSIGN: "=",
  KEYWORD_COMMA: ",",
  KEYWORD_SEMICOLON: ";",
  KEYWORD_LET: "let",
}

type KeywordToken struct {
  Keyword Keyword
}

func NewKeywordToken(l string) (*KeywordToken, error) {
  for i, s := range keywords {
    if l == s {
      return &KeywordToken{ Keyword: Keyword(i) }, nil
    }
  }
  return nil, errors.New(fmt.Sprintf("'%s' is not a keyword", l))
}

func (t *KeywordToken) Value() string {
  return keywords[t.Keyword]
}

func (t *KeywordToken) Type() TokenType {
  return KEYWORD
}

type Token interface {
  Value() string
  Type() TokenType
}
