package aoc

func main() {
  config, err := defaultConfig()
  checkError(err)
  addFlags(config)
  download(config)
}

