package lex

import (
  "fmt"
  "unicode"
  "log"
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
  return c == '+' || c == '-' || c == '*' || c == '/' || c == '=' || c == '!' || c == '>' || c == '<' || c == ':' || c == '.' || c == ',';
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
  return c == '+' || c == '-' || c == '*' || c == '/' || c == '=' || c == '>' || c == '<' || c == ':' || c == '.' || c == ',';
}

// Operator charset is generated
// from characters used to begin operators currently.
func unaryOperatorCheck(c rune) bool {
  return c == '!'
}

func semicolonState(l *Lexer) StateFunc {
  // End of expression
  // Requires all parenthesis are matched and not inside
  // a function call.
  if len(l.funcDepth) > 0 || l.parenDepth > 0 {
    l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Expected ')', found ';'\n") }
    return nil
  }
  l.Emit(&token.GenericToken{ Typ: token.KEYWORD, Val: l.Lexed() })
  return startState
}

func binaryOperatorState(l *Lexer) StateFunc {
  l.NextUpTo(operatorPredicate)
  l.Back()
  // Check if an operator
  t, err := token.NewBinaryOperatorToken(l.Lexed())
  if err != nil {
    fmt.Println(err)
    return nil
  }
  l.Emit(t)
  return expressionState
}

func unaryOperatorState(l *Lexer) StateFunc {
  l.NextUpTo(operatorPredicate)
  l.Back()
  // Check if an operator
  t, err := token.NewUnaryOperatorToken(l.Lexed())
  if err != nil {
    fmt.Println(err)
    return nil
  }
  l.Emit(t)
  return expressionState
}

func floatState(l *Lexer) StateFunc {
  c := l.NextUpTo(numberPredicate)
  if c == 'e' && len(l.lexed) > 1 {
    l.NextUpTo(numberPredicate)
    l.Back()
    if len(l.lexed) > 0 {
      l.Emit(&token.GenericToken{ Typ: token.FLOAT, Val: l.Lexed() })
      return expressionPostValueState
    }
  }
  l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Not a proper float '%s'\n", l.Lexed()) }
  return nil
}

func numberState(l *Lexer) StateFunc {
  switch c := l.NextUpTo(numberPredicate); {
  case c == '.':
    return floatState
  default:
    l.Back()
    l.Emit(&token.GenericToken{ Typ: token.INTEGER, Val: l.Lexed() })
    return expressionPostValueState
  }
}

func identifierState(l *Lexer) StateFunc {
  l.NextUpTo(identifierPredicate)
  l.Back()
  l.Emit(&token.GenericToken{ Typ: token.IDENT, Val: l.Lexed() })
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == '(' {
    // Parenthesized expression
    l.Emit(&token.GenericToken{ Typ: token.LPAREN, Val: l.Lexed() })
    l.parenDepth += 1
    return expressionState
  }
  l.Back()
  return expressionPostValueState
}

// After lexing an identifier or literal
func expressionPostValueState(l *Lexer) StateFunc {
  switch c := l.IgnoreUpTo(whitespacePredicate); {
  case c == ')':
    return expressionRightParenState
  case c == ',':
    return expressionCommaState
  case c == ';':
    return semicolonState
  case binaryOperatorCheck(c):
    return binaryOperatorState
  default:
    l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Unexpected rune '%c' after expression literal\n", c) }
    return nil
  }
}

func expressionRightParenState(l *Lexer) StateFunc {
  switch {
  case l.parenDepth > 0:
    l.Emit(&token.GenericToken{ Typ: token.RPAREN, Val: l.Lexed() })
    l.parenDepth -= 1
    return expressionPostValueState
  default:
    l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Extraneous right parenthesis '%s'\n", l.Lexed()) }
    return nil
  }
}

func expressionState(l *Lexer) StateFunc {
  // Check for identifiers or type literals
  switch c := l.IgnoreUpTo(whitespacePredicate); {
  case c == '(':
    // Parenthesized expression
    l.Emit(&token.GenericToken{ Typ: token.LPAREN, Val: l.Lexed() })
    l.parenDepth += 1
    return expressionState
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
    l.Emit(&token.GenericToken{ Typ: token.STRING, Val: l.Lexed() })
    return expressionPostValueState
  case charCheck(c):
    l.NextUpTo(charPredicate)
    l.Emit(&token.GenericToken{ Typ: token.CHAR, Val: l.Lexed() })
    return expressionPostValueState
  default:
    l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Unexpected rune '%c' in expression\n", c) }
    return nil
  }
}

func letAssignmentState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == '=' {
    t, err := token.NewKeywordToken(l.Lexed())
    if err != nil {
      log.Fatal(err) // should not happen
    }
    l.Emit(t)
    return expressionState
  }
  l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Expected assignment, found unexpected rune '%c'\n", c) }
  return nil
}

// First time inside a tuple, lexing type subexpression
func letTypeTupleState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == '(' {
    l.parenDepth += 1
    l.Emit(&token.GenericToken{ Typ: token.LPAREN, Val: l.Lexed() })
    return letTypeState
  }
  l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Expected parenthesized subtype expression after tuple, found unexpected run '%v'\n", c) }
  return nil
}

// letTypeState has discovered that it is lexing a type nested in a tuple.
func letInsideTupleState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate)
  if c == ')' {
    l.parenDepth -= 1
    l.Emit(&token.GenericToken{ Typ: token.RPAREN, Val: l.Lexed() })
    if l.parenDepth > 0 {
      // still in tuple, expect comma or closing paren
      return letInsideTupleState
    } else {
      return letAssignmentState
    }
  } else if c == ',' {
    t, err := token.NewKeywordToken(l.Lexed())
    if err != nil {
      log.Fatal(err) // should not happen
    }
    l.Emit(t)
    // need to lex another type in tuple definition
    return letTypeState
  }
  l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Unexpected rune inside tuple '%v'\n", c) }
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
    t, err := token.NewTypeToken(l.Lexed())
    if err != nil {
      fmt.Println(err)
      return nil
    }
    l.Emit(t)
    if t.Typ == token.TYPE_TUPLE {
      // tuples are nested types tuple(i32, bool)
      return letTypeTupleState
    }
    if l.parenDepth > 0 {
      return letInsideTupleState
    }
    return letAssignmentState
  }
  l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Expected identifier in type expression, found unexpected rune '%c'", c) }
  return nil
}

// keywordState has lexed the 'let' keyword,
// we expect to find an identifier
func letIdentState(l *Lexer) StateFunc {
  c := l.IgnoreUpTo(whitespacePredicate) // kill whitespace
  if identifierCheck(c) {
    l.NextUpTo(identifierPredicate)
    l.Back()
    l.Emit(&token.GenericToken{ Typ: token.IDENT, Val: l.Lexed() })
    c := l.IgnoreUpTo(whitespacePredicate)
    if c == ':' {
      t, err := token.NewKeywordToken(l.Lexed())
      if err != nil {
        log.Fatal(err) // should not happen
      }
      l.Emit(t)
      return letTypeState
    } else {
      l.Back()
      return letAssignmentState
    }
    // fallthrough and error
  }
  l.Emit(&token.GenericToken{ Typ: token.ERROR, Val: fmt.Sprintf("Unexpected rune '%v'\n", c) }
  return nil
}

// startState has lexed the first character
// of what we expect to be a keyword.
func keywordState(l *Lexer) StateFunc {
  l.NextUpTo(identifierPredicate)
  l.Back()
  if l.Lexed() == "let" {
    t, err := token.NewKeywordToken(l.Lexed())
    if err != nil {
      log.Fatal(err) // should not happen
    }
    l.Emit(t)
    return letIdentState
  }
  return identifierState
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
  l.Back()
  return expressionState
}
