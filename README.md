# envcmd

Command line tool for running per-environment commands.

## Installation

Build from source, install via Homebrew or see [Releases](https://github.com/brooknullsh/envcmd/releases).

### Source

1. Clone the repository:

```sh
git clone https://github.com/brooknullsh/envcmd.git
```

2. Build the binary (requires [Go](https://go.dev/doc/install) to be installed):

```sh
cd envcmd
go build -o ./bin/envcmd .
```

3. Run the binary in-place or add to your path manually:

```sh
./bin/envcmd
```

### Homebrew

```sh
brew install brooknullsh/tap/envcmd
```

## Usage

1. Run `envcmd create` or `envcmd c` to create your config file at
   `$HOME/.envcmd/config.json`
2. Run `envcmd list` or `envcmd l` to list all configurations
3. Change the config file with your directory/branch name and the commands to
   run
4. Run `envcmd` to execute the commands matching the directory/branch
