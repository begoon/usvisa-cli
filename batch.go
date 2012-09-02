package main

import (
  "bytes"
  "compress/zlib"
  "io/ioutil"
  "log"
  "net/http"
  "regexp"
  "strconv"
  "strings"
)

const (
  Url = "http://photos.state.gov/libraries/unitedkingdom/164203/cons-visa/admin_processing_dates.pdf"
)

type BatchUpdate struct {
  Status, Date string
}

type BatchTable map[string][]BatchUpdate

var (
  BTETRE = *regexp.MustCompile("(?ms)BT\\r\\n(.+?)ET\\r\\n")
  TextRE = *regexp.MustCompile("\\((.+?)\\)")
)

func loadFile(name string) []byte {
  bytes, err := ioutil.ReadFile(name)
  if err != nil {
    log.Printf("Unable to read file [%s]", name)
    return nil
  }
  log.Printf("Loaded: %d\n", len(bytes))
  return bytes
}

func loadFromUrl(url string) []byte {
  response, err := http.Get(url)
  log.Printf("Started downloading")
  if err != nil {
    log.Printf("GET failed: %s", err)
    return nil
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    log.Printf("GET read failed: %s", err)
    return nil
  }
  log.Printf("Loaded: %d\n", len(contents))
  return contents
}

const (
  StreamStartMarker = "stream\x0D\x0A"
  StreamEndMarker   = "endstream\x0D\x0A"
)

func LoadTable(url string) BatchTable {
  var pdf []byte
  if strings.HasPrefix(url, "http") {
    pdf = loadFromUrl(url)
  } else {
    pdf = loadFile(url)
  }
  if pdf == nil {
    log.Printf("PDF file wasn't loaded")
    return nil
  }

  table := make(BatchTable)

  for {
    begin := bytes.Index(pdf, []byte(StreamStartMarker))
    if begin == -1 {
      break
    }
    pdf = pdf[begin+len(StreamStartMarker):]
    end := bytes.Index(pdf, []byte(StreamEndMarker))
    if end == -1 {
      break
    }
    section := pdf[0:end]
    pdf = pdf[end+len(StreamEndMarker):]

    buf := bytes.NewBuffer(section)
    unzipReader, err := zlib.NewReader(buf)
    if err != nil {
      log.Printf("Unzip initialization failed, %v", err)
      continue
    }
    unzipped, err := ioutil.ReadAll(unzipReader)
    if err != nil {
      log.Printf("Unzip failed, %v", err)
      continue
    }
    records := make([]string, 0)
    for _, group := range BTETRE.FindAllSubmatch(unzipped, -1) {
      lines := make([][]byte, 0)
      for _, group := range TextRE.FindAllSubmatch(group[1], -1) {
        lines = append(lines, group[1])
      }
      records = append(records, string(bytes.Join(lines, []byte{})))
    }
    for i := 0; i < len(records)-2; i++ {
      v, err := strconv.ParseInt(records[i], 10, 64)
      if err == nil && v >= 20000000000 && v < 29000000000 {
        id := records[i]
        if _, exists := table[id]; !exists {
          table[id] = make([]BatchUpdate, 0)
        }
        table[id] = append(table[id], BatchUpdate{records[i+1], records[i+2]})
        i += 2
      }
    }
  }
  return table
}

func PrintTable(table BatchTable) {
  for id, updates := range table {
    log.Printf("%s (%d)", id, len(updates))
    for _, update := range updates {
      log.Printf("%s, %s", update.Status, update.Date)
    }
  }
}
