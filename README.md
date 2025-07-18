# envcmd

Command line tool for running per-environment commands.

## Installation

There are a couple of ways to use `envcmd`:

### Source

1. Clone the repository:

```sh
git clone https://github.com/brooknullsh/envcmd.git
```

2. Build the binary:

> Requires [Cargo](https://github.com/rust-lang/cargo) to be installed

```sh
cd envcmd
cargo build -r
```

3. Run the binary in-place or add to your path:

```sh
./target/release/envcmd
```

### Homebrew

```sh
brew install brooknullsh/tap/envcmd
```

### Releases

See [the releases](https://github.com/brooknullsh/envcmd/releases) for the
latest version.

## Usage

> Run `envcmd help`, `envcmd -h` or `envcmd --help` to see usage

1. Run `envcmd create` or `envcmd c` to create your config file at
   `$HOME/.envcmd/config.json`

> Run `envcmd delete` or `envcmd d` to delete your config file

2. Run `envcmd list` or `envcmd l` to list all configurations
3. Change the config file with your directory/branch name and the commands to
   run
4. Run `envcmd` to execute the commands matching the directory/branch
