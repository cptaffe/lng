package parse

import (
  "fmt"
  "../token"
)

type SyntaxTree struct {
  Value token.Token
  Children []*SyntaxTree
}

type Parser struct {
  Input chan token.Token
  stack []*SyntaxTree
  State StateFunc
}

func (s *SyntaxTree) Append(t *SyntaxTree) {
  s.Children = append(s.Children, t)
}

// Pushes a SyntaxTree onto the stack,
// setting it as root if root is nil
func (p *Parser) Push(s *SyntaxTree) {
  p.stack = append(p.stack, st)
}

// Pops a SyntaxTree from the stack
func (p *Parser) Pop() *SyntaxTree {
  if len(p.stack) > 0 {
    s := p.stack[len(p.stack)-1]
    p.stack = p.stack[:len(p.stack)-1]
    return s
  }
  return nil
}

type StateFunc func(*Parser)StateFunc

func letTypeState(p *Parser) StateFunc {
  // Expecting type
  t := <-p.Input
  if t.Type() == token.TYPE {
    k := t.(token.TypeToken)
    if k != nil {
      select k.Typ {
      case token.TYPE_TUPLE:
        return letTupleState
      default:
        // Expect '=' keyword

      }
    }
  }
  fmt.Printf("Expected type, found %#v\n", t)
}

func letStatementState(p *Parser) StateFunc {
  // Expecting ':' or '=' keyword
  t := <-p.Input
  if t.Type() == token.KEYWORD {
    k := t.(token.KeywordToken)
    if k != nil {
      select k.Keyword {
      case token.KEYWORD_ASSIGN:
        s := p.Pop()

        p.PushLast(&SyntaxTree{ Value: t })
        return letAssignState
      case token.KEYWORD_TYPE:
        p.Push(&SyntaxTree{ Value: t })
        return letTypeState
      }
    }
  }
  fmt.Printf("Expected keyword ':' or '=', found %#v\n", t)
  return nil
}

func startState(p *Parser) StateFunc {
  // Expecting 'let' keyword
  t := <-p.Input
  if t.Type() == token.KEYWORD {
    k := t.(token.KeywordToken)
    if k != nil && k.Keyword == token.KEYWORD_LET {
      p.Push(&SyntaxTree{ Value: t })
      return letStatementState
    }
  }
  fmt.Printf("Expected keyword 'let', found %#v\n", t)
  return nil
}

func (p *Parser) SyntaxTreeGenerator() chan *SyntaxTree {
  c := make(chan *SyntaxTree)
  go func(chan *SyntaxTree) {
    l.State = startState
    for l.State != nil {
      l.State = l.State(p)
    }
    close(c)
  }(c)
  return c
}

func (p *Parser) Parse() {
  for c := range p.SyntaxTreeGenerator() {
    fmt.Println(c)
  }
}
