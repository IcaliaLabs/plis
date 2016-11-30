# About setting up Go

This is just a friendly reminder for myself, whenever I need to install Go on macOS, I need to set
the following variables on the `~/.bash_profile` file:

```bash
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

Most of the available online documentation and resources I found missed to set the `GOBIN` variable,
which causes an error whenever running the `go get` command in the project's directory.
