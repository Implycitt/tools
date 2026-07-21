package main 

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
  "errors"
)

func download(config *configuration) error {
  client := new(http.Client)
  cookie := new(http.Cookie)

  req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", config.Year, config.Day), nil)
  if err != nil { return err }

  cookie.Name, cookie.Value = "session", config.SessionCookie 
  req.AddCookie(cookie)

  resp, err := client.Do(req)
  if err != nil { return err }

  defer resp.Body.Close()

  _, err = os.Stat(config.Path)
  if errors.Is(err, os.ErrNotExist) {
    os.MkdirAll(config.Path, os.ModePerm)
  } 

  file, err := os.OpenFile(fmt.Sprintf("%s/%s", config.Path, config.Output), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

  defer file.Close()

  _, err = io.Copy(file, resp.Body)
  if err != nil { return err }

  return nil
}

func addFlags(config *configuration) {
  var year, day int
  flags := flag.NewFlagSet("", flag.ContinueOnError)

  ignored := new(bytes.Buffer)
  flags.SetOutput(ignored)

  yearFlag := flags.String("year", "", "")
  dayFlag := flags.String("day", "", "")
  waitFlag := flags.Bool("force", false, "")
  flagErr := flags.Parse(os.Args[1:])

  if flagErr == nil {
    year, flagErr = parseIntFlag(*yearFlag)
  }
  if flagErr == nil {
    day, flagErr = parseIntFlag(*dayFlag)
  }
  if flagErr == flag.ErrHelp {
    fmt.Println("incorrect usage")
    os.Exit(0)
  }
  if flagErr != nil {
    fmt.Fprintln(os.Stderr, flagErr)
    os.Exit(1)
  }

  flagConfig := new(configuration)
  flagConfig.Year = year
  flagConfig.Day = day

  config.merge(flagConfig)
  
  if *waitFlag {
    config.Wait = true
  }
}

func parseIntFlag(flag string) (int, error) {
  if flag == "" { return 0, nil }
  value, err := strconv.ParseInt(flag, 10, 0)
  return int(value), err
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintln(os.Stderr, err) 
    os.Exit(1)
  }
}
