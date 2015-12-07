package main

import (
  "fmt"
  "./lex"
  "./parse"
)

func main() {
  // go func() {
  //   reader := bufio.NewReader(os.Stdin)
  //   for {
  //     fmt.Print(">> ")
  //     text, err := reader.ReadString('\n')
  //     if err != nil {
  //       if err == io.EOF {
  //         close(r)
  //         return
  //       }
  //       log.Fatal(err)
  //     }
  //     for _, c := range text {
  //       r <- c
  //     }
  //   }
  // }()

  texts := []string {
    "hello",
    "hello(4,5)",
    "4*5%6+3",
    "_(f(5), 7*8)",
    "h.string",
    "().tuple.do",
    "((((6,7,8),8))).string",
    "(7*8, 6+3)",
  }

  for _, text := range texts {
    r := make(chan rune)
    l := lex.Lexer { Input: r }
    p := &parse.Parser{ Input: l.TokenGenerator() }
    go func() {
      for _, c := range text {
        r <- c
      }
      close(r)
    }()
    fmt.Printf("%s->%s\n", text, p.Parse())
  }
}
