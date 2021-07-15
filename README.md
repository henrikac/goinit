# goinit

goinit is a CLI that makes it easy to generate a new Go project.  

A new project contains:
+ `.git`
+ `go.mod`
+ `README.md`
+ `LICENSE` (MIT)
+ `.gitignore`

## Installation

Run `go get github.com/henrikac/goinit` to install goinit.  

After goinit has been installed it is recommended to set `GO_INIT_PATH`. This will tell goinit where it should generate new projects.

An example could be `export GO_INIT_PATH=$HOME/go/src/github.com/your-github-user`

## Usage

#### Generate a new project

New projects will be generated in `GO_INIT_PATH` if non-empty, otherwise in `GOPATH`.
If both `GO_INIT_PATH` and `GOPATH` are empty a default path is used `user-home-dir/go` (e.g. `$HOME/go` on Unix).

Run `goinit new <project name>` to create a new project.

## Contributing

1. Fork it (<https://github.com/henrikac/goinit/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## Contributors

- [Henrik Christensen](https://github.com/henrikac) - creator and maintainer
