package main 

import(
  "fmt"
  "os"
  "time"
  "bufio"
  "io"
  "errors"
  "strings"
)

type configuration struct {
  SessionCookie string
	Output        string 
	Dir           string 
	Path          string
	Year          int 
	Day           int 
	Wait          bool   
}

func createConfigFile() {
  secret := ""
  path := ""
  homedir, err := os.UserHomeDir()
  if err != nil { return }

  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter your session cookie: ")
  secret, _ = reader.ReadString('\n')
  fmt.Print("Enter your input path: ")
  path, _ = reader.ReadString('\n')

  os.Chdir(homedir)
  file, err := os.Create("aoc.config")
  if err != nil { return }
  defer file.Close()
  fmt.Fprint(file, (fmt.Sprintf("%s\n%s", secret, path)))
}

func defaultConfig() (*configuration, error) {

  config := new(configuration)

  homedir, err := os.UserHomeDir()
  if err != nil { return config, err }
  os.Chdir(homedir)

  _, err = os.Stat("aoc.config")
  if errors.Is(err, os.ErrNotExist) {
    createConfigFile()
  } 

  file, err := os.Open("aoc.config")
  if err != nil { return config, err }
  defer file.Close()

  var configs []string
  reader := bufio.NewReader(file)
  for {
    line, err := reader.ReadString('\n')
    if err == io.EOF || err != nil { break }
    configs = append(configs, line) 
  }

  config.SessionCookie = string(strings.TrimSuffix(configs[0], "\r\n"))
  config.Dir = strings.TrimSuffix(configs[2], "\r\n")
  config.SessionCookie = strings.TrimSuffix(config.SessionCookie, "\n") // for some reason linux complains without this line 
  config.Dir = strings.TrimSuffix(config.Dir, "\n")

  est, err := time.LoadLocation("EST")
  if err != nil { os.Exit(1) }
  now := time.Now().In(est)

  if now.Month() != time.December {
    config.Day = 25
    config.Year = int(now.Year())-1
  } else {
    config.Day = now.Day()
    config.Year = now.Year()
  }

  config.Output = fmt.Sprintf("%d.txt", config.Day)
  config.Path = fmt.Sprintf("%s/%d", config.Dir, config.Year)

  return config, nil
}

func (config *configuration) merge(other *configuration) {
  if other.Year != 0 { config.Year = other.Year }
  if other.Day != 0 { config.Day = other.Day }
  config.Output = fmt.Sprintf("%d.txt", config.Day)
  config.Path = fmt.Sprintf("%s/%d", config.Dir, config.Year)
}
