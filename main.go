package main

import (
  "fmt"
  "os"

  "github.com/brooknullsh/envcmd/internal/command"
  "github.com/brooknullsh/envcmd/internal/config"
  "github.com/brooknullsh/envcmd/internal/log"
)

func printUsage() {
  fmt.Println("Command line tool for running per-environment commands")
  fmt.Println("\nUsage: envcmd [COMMAND]")
  fmt.Println("\nCommands:")
  fmt.Println("  create  c  Create configuration file.")
  fmt.Println("  delete  d  Delete configuration file.")
  fmt.Println("  list    l  Show configuration file contents.")
  fmt.Println("  help    h  Show this message.")
}

func main() {
  var cfg config.Config
  cfg.Init()

  args := os.Args[1:]
  if len(args) == 0 {
    command.Run(cfg)
    return
  }

  if len(args) > 1 {
    printUsage()
    return
  }

  switch args[0] {
  case "create", "c":
    cfg.Create()
  case "delete", "d":
    cfg.Delete()
  case "list", "l":
    cfg.List()
  case "help", "h":
    printUsage()
  default:
    log.Abort("unknown command: %s", args[0])
  }
}
