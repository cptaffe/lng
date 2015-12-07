package parse

import (
  "fmt"
  "../token"
)

type SyntaxTree struct {
  Value *token.Token
  Children []*SyntaxTree
}

func (s *SyntaxTree) String() string {
  str := ""
  str += fmt.Sprintf("%s %v", s.Value.Value, s.Children)
  return str
}

const (
  MAX_ERROR = 20
)

type Parser struct {
  Input <-chan *token.Token
}

type ShuntingYard struct {
  Input <-chan *token.Token
  stack []*token.Token
}

func (y *ShuntingYard) shunt() <-chan *token.Token {
  c := make(chan *token.Token)
  go func(){
    for t := range y.Input {
      fmt.Printf("'%s': %s\n", t.Value, y.stack)
      switch {
      case t.Type == token.ERROR:
        fmt.Printf("Error: %s\n", t.Value)
      case t.Type == token.INVALID:
        fmt.Printf("Encountered Invalid Token: '%s'\n", t.Value)
      case t.Type == token.INTEGER || t.Type == token.FLOAT || t.Type == token.IDENT:
        c <- t
      case t.Type == token.FUNC:
        y.stack = append(y.stack, t)
      case t.Type.IsOperator():
        for len(y.stack) > 0 {
          st := y.stack[len(y.stack)-1]
          if st.Type.IsOperator() {
            if (t.Type.Assoc() == token.LEFT_ASSOC && t.Type.Prec() <= st.Type.Prec()) || (t.Type.Assoc() == token.RIGHT_ASSOC && t.Type.Prec() < st.Type.Prec()) {
              c <- st
              y.stack = y.stack[:len(y.stack)-1]
              continue
            }
          }
          break
        }
        y.stack = append(y.stack, t)
      case t.Type == token.LPAREN:
        y.stack = append(y.stack, t)
      case t.Type == token.RPAREN:
        hasFound := false
        hasOther := 0
        for len(y.stack) > 0 {
          st := y.stack[len(y.stack)-1]
          y.stack = y.stack[:len(y.stack)-1]
          if st.Type == token.LPAREN {
            hasFound = true
            if len(y.stack) > 0 && y.stack[len(y.stack)-1].Type == token.FUNC {
                c <- y.stack[len(y.stack)-1]
                y.stack = y.stack[:len(y.stack)-1]
            } else if hasOther == 0 {
              // HACK HACK HACK
              // empty tuple, null field or something?
              c <- &token.Token{ Type: token.EMPTY, Value: "()" }
            }
            break
          } else {
            hasOther += 1
            c <- st
          }
        }
        if !hasFound {
          panic("Mismatched parenthesis")
        }
      }
    }
    // EOF
    for len(y.stack) > 0 {
      st := y.stack[len(y.stack)-1]
      if st.Type == token.RPAREN {
        panic("Mismatched parenthesis")
      }
      c <- st
      y.stack = y.stack[:len(y.stack)-1]
    }
    close(c)
  }()
  return c
}

func (y *ShuntingYard) Shunt() []*SyntaxTree {
  stack := []*SyntaxTree{}
  for t := range y.shunt() {
    switch {
    case t.Type.IsOperator():
      s := make([]*SyntaxTree, 2)
      copy(s, stack)
      st := &SyntaxTree{ Value: t, Children: s }
      if len(stack) > 2 {
        stack = stack[2:]
      } else {
        stack = []*SyntaxTree{}
      }
      stack = append(stack, st)
    case t.Type == token.FUNC:
      // Parent tuple
      st := &SyntaxTree{ Value: t }
      if len(stack) > 0 {
        st.Children = stack[:1]
        stack = stack[1:]
      }
      stack = append(stack, st)
    default:
      stack = append(stack, &SyntaxTree{ Value: t })
    }
  }
  return stack
}

func (p *Parser) Parse() []*SyntaxTree {
  sy := &ShuntingYard{ Input: p.Input }
  return sy.Shunt()
}
