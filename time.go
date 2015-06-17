package gocast

import (
  "time"
)

var timeFormats = []string{
  time.RFC1123Z,
  time.RFC3339Nano,
  time.UnixDate,
  time.RubyDate,
  time.RFC1123,
  time.RFC3339,
  time.RFC822,
  time.RFC850,
  time.RFC822Z,
}

func parseTime(tm string) (t time.Time, err error) {
  for _, f := range timeFormats {
    if t, err = time.Parse(f, tm); nil == err {
      break
    }
  }
  return
}
