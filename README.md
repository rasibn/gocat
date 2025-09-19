# gocat

## Install

go install github.com/rasibn/gocat@v0.1.1

## Options

```bash
NAME:
   gocat - Concatenate files in a directory using ripgrep + optional fzf

USAGE:
   gocat [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --fzf       Use fzf to interactively pick files (press Tab to multi-select) (default: false)
   --hidden    Include hidden files (dotfiles) (default: false)
   --help, -h  show help
```

## Example

```bash
gocat --hidden --fzf ./src
```

```bash
gocat ./src
```
