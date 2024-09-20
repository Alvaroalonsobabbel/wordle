# Terminal Wordle
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Alvaroalonsobabbel/wordle) ![Test](https://github.com/Alvaroalonsobabbel/wordle/actions/workflows/test.yml/badge.svg) ![Latest Release](https://img.shields.io/github/v/release/Alvaroalonsobabbel/wordle?color=blue&label=Latest%20Release)

Play the NYT daily Wordle from the comfort of your terminal!

<img src="doc/example.gif" alt="drawing" width="282"/>

⚠️ this assumes you know how to use the terminal! If you don't you can find out how [here](https://www.google.com/search?q=how+to+use+the+terminal).

## Install

For Apple computers with ARM chips you can use the provided installer. For any other OS you'll have to compile the binary yourself.

### ARM (Apple Sillicon)

Open the terminal and run:

```bash
curl -sSL https://raw.githubusercontent.com/Alvaroalonsobabbel/wordle/main/bin/install.sh | bash
```

- You'll be required to enter your admin password.
- You might be required to allow the program to run in the *System Settings - Privavacy & Security* tab.

### Compiling the binary yourself

1. [Install Go](https://go.dev/doc/install)
2. Clone the repo `git clone git@github.com:Alvaroalonsobabbel/wordle.git`
3. CD into the repo `cd wordle`
4. Build the binary `make build`

## How to Play

You can check the official Wordle rules [here](https://www.nytimes.com/2023/08/01/crosswords/how-to-talk-about-wordle.html).

1. Start Wordle by running `wordle` in your Terminal.
2. Have fun!

You can quit the game at any time by pressing `Esc` or `Ctrl C`

## Options

### Hard Mode

Enables Worlde's hard mode.
`wordle -hard`
