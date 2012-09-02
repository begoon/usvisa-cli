package main

import (
  "io/ioutil"
  "log"
  "net/http"
  "runtime"
  "strings"
)

func checker(jobs chan string, ack chan bool, table BatchTable) {
  for {
    id := <-jobs
    if id == "" {
      return
    }
    resp, err := http.Get("http://localhost:8080/batch/" + id)
    if err != nil {
      log.Printf("GET failed, %s", err)
      ack <- true
      continue
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Printf("GET read failed, %s", err)
      ack <- true
      continue
    }
    lines := strings.Split(string(body), "\n")
    updates := make([]BatchUpdate, 0)
    for i := 0; i < len(lines); i++ {
      if len(lines[i]) == 0 {
        continue
      }
      updates = append(updates, BatchUpdate{lines[i], lines[i+1]})
      i += 1
    }
    expectedUpdates := table[id]
    if len(expectedUpdates) != len(updates) {
      log.Fatalf("Updates length doesn't match for '%s'", id)
    }
    for i, update := range expectedUpdates {
      if update.Status != updates[i].Status {
        log.Fatalf("Status doesn't match for '%s', '%s' != '%s'", id, update.Status, updates[i].Status)
      }
      if update.Date != updates[i].Date {
        log.Fatalf("Date doesn't match for '%s', '%s' != '%s'", id, update.Date, updates[i].Date)
      }
    }
    ack <- true
  }
}

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())
  table := LoadTable(Url)
  for {
    jobs := make(chan string, len(table))
    ack := make(chan bool, len(table))
    go checker(jobs, ack, table)
    for id, _ := range table {
      jobs <- id
    }
    jobs <- ""
    for _ = range table {
      <-ack
    }
    log.Printf("Checked %d records", len(table))
  }
}
