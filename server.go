package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"
  "sync"
  "time"
)

var (
  storage = make(BatchTable)
  lock    sync.RWMutex
)

func TableUpdater() {
  for {
    time.Sleep(5 * time.Second)
    log.Printf("Reload")
    table := LoadTable(Url)
    lock.Lock()
    storage = table
    lock.Unlock()
  }
}

func GetBatch(id string) []BatchUpdate {
  lock.RLock()
  defer lock.RUnlock()
  return storage[id]
}

func Batch(w http.ResponseWriter, r *http.Request) {
  parts := strings.Split(r.URL.Path, "/")
  id := parts[2]
  updates := GetBatch(id)
  if updates == nil {
    fmt.Fprintf(w, "Batch not found")
  } else {
    for _, update := range updates {
      fmt.Fprintf(w, "%s\n%s\n\n", update.Status, update.Date)
    }
  }
}

func Server() {
  http.Handle("/batch/", http.HandlerFunc(Batch))

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
  storage = LoadTable(Url)
  go Server()
  TableUpdater()
}
