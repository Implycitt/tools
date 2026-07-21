package main 

func main() {
  config, err := defaultConfig()
  checkError(err)
  addFlags(config)
  download(config)
}

