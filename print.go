package main

import (
  "bytes"
  "encoding/gob"
  "log"
)

func SerializeTable() {
  table := LoadTable(Url)
  var b bytes.Buffer
  err := gob.NewEncoder(&b).Encode(table)
  if err != nil {
    log.Fatal("GOB encode error, ", err)
  }
  gob := b.Bytes()
  println(len(gob))
}

func main() {
  PrintTable(LoadTable(Url))
}
