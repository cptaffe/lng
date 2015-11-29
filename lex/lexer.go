package lex

import (
  "../token"
)

type StateFunc func(*Lexer)StateFunc

type Lexer struct {
  Input chan rune
  State StateFunc
  lexed []rune
  back []rune
  toks chan token.Token
  parenDepth int
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

func (l *Lexer) Emit(t token.Token) {
  l.toks <- t
  l.lexed = []rune{}
}

func (l *Lexer) Lexed() string {
  s := ""
  for _, r := range l.lexed {
    s += string(r)
  }
  return s
}

func (l *Lexer) TokenGenerator() chan token.Token {
  l.toks = make(chan token.Token)
  go func(l *Lexer) {
    // state machine
    l.State = startState
    for l.State != nil {
      l.State = l.State(l)
    }
    close(l.toks)
  }(l)
  return l.toks
}
