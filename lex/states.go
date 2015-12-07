package lex

import (
  "fmt"
  "unicode"
  "../token"
)

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
  return c == '+' || c == '-' || c == '*' || c == '/' || c == '%' || c == '=' || c == '!' || c == '>' || c == '<' || c == ':' || c == '.' || c == ',';
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

// Operator charset is generated
// from characters used to begin operators currently.
func binaryOperatorCheck(c rune) bool {
  return c == '+' || c == '-' || c == '*' || c == '%' || c == '/' || c == '=' || c == '>' || c == '<' || c == ':' || c == '.' || c == ',';
}

// Operator charset is generated
// from characters used to begin operators currently.
func unaryOperatorCheck(c rune) bool {
  return c == '!'
}

func binaryOperatorState(l *Lexer) StateFunc {
  l.NextUpTo(operatorPredicate)
  l.Back()
  l.Emit(&token.Token{ Type: token.BinaryOperatorType(l.Lexed()), Value: l.Lexed() })
  return expressionState
}

func unaryOperatorState(l *Lexer) StateFunc {
  l.NextUpTo(operatorPredicate)
  l.Back()
  l.Emit(&token.Token{ Type: token.UnaryOperatorType(l.Lexed()), Value: l.Lexed() })
  return expressionState
}

func floatState(l *Lexer) StateFunc {
  c := l.NextUpTo(numberPredicate)
  if c == 'e' && len(l.lexed) > 1 {
    l.NextUpTo(numberPredicate)
    l.Back()
    if len(l.lexed) > 0 {
      l.Emit(&token.Token{ Type: token.FLOAT, Value: l.Lexed() })
      return expressionPostValueState
    }
  }
  l.Emit(&token.Token{ Type: token.ERROR, Value: fmt.Sprintf("Not a proper float '%s'\n", l.Lexed()) })
  return nil
}

func numberState(l *Lexer) StateFunc {
  switch c := l.NextUpTo(numberPredicate); {
  case c == '.':
    return floatState
  default:
    l.Back()
    l.Emit(&token.Token{ Type: token.INTEGER, Value: l.Lexed() })
    return expressionPostValueState
  }
}

func identifierState(l *Lexer) StateFunc {
  l.NextUpTo(identifierPredicate)
  l.Back()
  s := l.Lexed()
  c := l.IgnoreUpTo(whitespacePredicate)
  l.Back()
  if c == '(' {
    l.Emit(&token.Token{ Type: token.FUNC, Value: s })
    l.Next()
    return expressionLeftParenState
  } else {
    l.Emit(&token.Token{ Type: token.IDENT, Value: s })
    return expressionPostValueState
  }
}

// After lexing an identifier or literal
func expressionPostValueState(l *Lexer) StateFunc {
  switch c := l.IgnoreUpTo(whitespacePredicate); {
  case c == 0:
    return nil // EOF
  case c == ')':
    return expressionRightParenState
  case binaryOperatorCheck(c):
    return binaryOperatorState
  default:
    l.Emit(&token.Token{ Type: token.ERROR, Value: fmt.Sprintf("Unexpected rune '%c' after expression literal\n", c) })
    return nil
  }
}

func expressionRightParenState(l *Lexer) StateFunc {
  switch {
  case l.parenDepth > 0:
    l.Emit(&token.Token{ Type: token.RPAREN, Value: l.Lexed() })
    l.parenDepth -= 1
    return expressionPostValueState
  default:
    l.Emit(&token.Token{ Type: token.ERROR, Value: fmt.Sprintf("Extraneous right parenthesis '%s'\n", l.Lexed()) })
    return nil
  }
}

func expressionLeftParenState(l *Lexer) StateFunc {
  l.Emit(&token.Token{ Type: token.LPAREN, Value: l.Lexed() })
  l.parenDepth += 1
  return expressionState
}

func expressionState(l *Lexer) StateFunc {
  // Check for identifiers or type literals
  switch c := l.IgnoreUpTo(whitespacePredicate); {
  case c == 0:
    return nil // EOF
  case c == '(':
    return expressionLeftParenState
  case c == ')':
    return expressionRightParenState
  case unaryOperatorCheck(c):
    return unaryOperatorState
  case identifierCheck(c):
    return identifierState
  case numberCheck(c):
    return numberState
  case stringCheck(c):
    l.NextUpTo(stringPredicate)
    l.Emit(&token.Token{ Type: token.STRING, Value: l.Lexed() })
    return expressionPostValueState
  case charCheck(c):
    l.NextUpTo(charPredicate)
    l.Emit(&token.Token{ Type: token.CHAR, Value: l.Lexed() })
    return expressionPostValueState
  default:
    l.Emit(&token.Token{ Type: token.ERROR, Value: fmt.Sprintf("Unexpected rune '%c' in expression\n", c) })
    return nil
  }
}

func startState(l *Lexer) StateFunc {
  return expressionState
}
