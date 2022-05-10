
# Minesweeper 🚩💣

Play Minesweeper in your terminal!

<p align="center">
    <img src="https://i.postimg.cc/L5GGYZYG/demo.png" alt="Demo Screenshot">
</p>

<p align="center">
    <a href="https://goreportcard.com/report/github.com/benhsm/minesweeper"><img src="https://goreportcard.com/badge/benhsm/minesweeper" alt="Go ReportCard"></a>
</p>

A simple TUI implementation of the immortal desktop video game,
[Minesweeper](https://en.wikipedia.org/wiki/Minesweeper_(video_game)), made
using the wonderful [Charm libraries](https://charm.sh/libs/),
[Bubbletea](https://github.com/charmbracelet/bubbletea) and
[Lipgloss](https://github.com/charmbracelet/lipgloss).

As of now, only the most essential features that make the game playable are
implemented; the first move is not guaranteed to be safe, and board
configurations in which the player is forced to guess may occur.

## Installation

1. Grab the latest release binary for your platform

2. Build it from source

```sh
$ git clone https://github.com/benhsm/minesweeper.git
$ cd minesweeper
$ go build
$ ./minesweeper
```
