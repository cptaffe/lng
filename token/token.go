package token

type Type int

// Tokens
const (
  ERROR Type = iota
  INVALID
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

  _OPERATOR_BEGIN
  _OPERATOR_BINARY_BEGIN
  OPERATOR_ADD // +
  OPERATOR_SUB // -
  OPERATOR_MUL // *
  OPERATOR_DIV // /
  OPERATOR_MOD // %
  OPERATOR_COMMA // ,
  OPERATOR_DOT // .
  OPERATOR_TYPE // :
  OPERATOR_EQ // ==
  OPERATOR_NEQ // !=
  OPERATOR_GT // >
  OPERATOR_GTE // >=
  OPERATOR_LT // <
  OPERATOR_LTE // <=
  _OPERATOR_BINARY_END
  _OPERATOR_UNARY_BEGIN
  OPERATOR_NOT // !
  _OPERATOR_UNARY_END
  _OPERATOR_END

  // HACK: Special use in Parser
  EMPTY
)

type Token struct {
  Type Type
  Value string
}

func (t *Token) String() string {
  return t.Value
}

// maps operators to strings
var operators = []string {
  OPERATOR_ADD: "+",
  OPERATOR_SUB: "-",
  OPERATOR_MUL: "*",
  OPERATOR_DIV: "/",
  OPERATOR_MOD: "%",
  OPERATOR_COMMA: ",",
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

func Operator(t Type) string {
  if t > _OPERATOR_BEGIN && t < _OPERATOR_END {
    return operators[t]
  }
  return ""
}

func BinaryOperatorType(s string) Type {
  for i := _OPERATOR_BINARY_BEGIN+1; i < _OPERATOR_BINARY_END; i++ {
    if s == operators[i] {
      return i
    }
  }
  return INVALID
}

func UnaryOperatorType(s string) Type {
  for i := _OPERATOR_UNARY_BEGIN+1; i < _OPERATOR_UNARY_END; i++ {
    if s == operators[i] {
      return i
    }
  }
  return INVALID
}

func (t Type) IsOperator() bool {
  return t > _OPERATOR_BEGIN && t < _OPERATOR_END
}

func (t Type) IsBinaryOperator() bool {
  return t > _OPERATOR_BINARY_BEGIN && t < _OPERATOR_BINARY_END
}

func (t Type) IsUnaryOperator() bool {
  return t > _OPERATOR_UNARY_BEGIN && t < _OPERATOR_UNARY_END
}

func (t Type) Prec() int {
  switch t {
  case OPERATOR_COMMA:
    return 1
  case OPERATOR_ADD:
    return 2
  case OPERATOR_SUB:
    return 2
  case OPERATOR_MUL:
    return 3
  case OPERATOR_DIV:
    return 3
  case OPERATOR_MOD:
    return 3
  case OPERATOR_DOT:
    return 4
  }
  return 15
}

const (
  LEFT_ASSOC = iota
  RIGHT_ASSOC
)

func (t Type) Assoc() int {
  if t == OPERATOR_NOT {
    return RIGHT_ASSOC
  }
  return LEFT_ASSOC
}
