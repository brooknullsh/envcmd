package main

import (
  "flag"
  "fmt"
  "os"

  "github.com/brooknullsh/envcmd/internal/command"
  "github.com/brooknullsh/envcmd/internal/config"
)

// build flag provided by goreleaser, using latest git tag
// https://goreleaser.com/cookbooks/using-main.version/
var version string

func printUsage() {
  fmt.Println("Command line tool for running per-environment commands")
  fmt.Println("\nUsage: envcmd [COMMAND]")
  fmt.Println("\nCommands:")
  fmt.Println("  create   c  Create configuration file.")
  fmt.Println("  delete   d  Delete configuration file.")
  fmt.Println("  list     l  Show configuration file contents.")
  fmt.Println("  add      a  Add new commands to run.")
  fmt.Println("              -a async     bool")
  fmt.Println("              -k kind      string")
  fmt.Println("              -t target    string")
  fmt.Println("              .. commands  string[]")
  fmt.Println("  version  v  Show the version you're running.")
  fmt.Println("  help     h  Show this message.")
}

func printVersion() {
  fmt.Printf("envcmd (v%s)\n", version)
}

func setupFlags() (*flag.FlagSet, *bool, *string, *string) {
  flags := flag.NewFlagSet("add", flag.ExitOnError)

  async := flags.Bool("a", false, "async")
  kind := flags.String("k", "", "kind")
  target := flags.String("t", "", "target")
  return flags, async, kind, target
}

func main() {
  var cfg config.Config
  cfg.Init()

  args := os.Args[1:]
  if len(args) == 0 {
    command.Run(cfg)
    return
  }

  switch args[0] {
  case "create", "c":
    cfg.Create()
  case "delete", "d":
    cfg.Delete()
  case "list", "l":
    cfg.List()
  case "add", "a":
    flags, async, kind, target := setupFlags()
    if err := flags.Parse(args[1:]); err != nil {
      printUsage()
      return
    }

    commands := flags.Args()
    if *kind == "" || *target == "" || len(commands) == 0 {
      printUsage()
      return
    }

    newCfg := config.Config{
      Async:    *async,
      Kind:     *kind,
      Target:   *target,
      Commands: commands,
    }

    cfg.Add(newCfg)
  case "version", "v":
    printVersion()
  case "help", "h":
    printUsage()
  default:
    printUsage()
  }
}
