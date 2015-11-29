package main

import (
  "fmt"
  "os"
  "bufio"
  "io"
  "log"
  "./lex"
  // "./parse"
)

func main() {
  r := make(chan rune)
  l := lex.Lexer { Input: r }
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
  // p := &parse.Parser{ Input: c }
  // p.Parse()
  for t := range l.TokenGenerator() {
    fmt.Printf("%#v\n", t)
  }
}
