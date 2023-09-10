# steel

An experimental local development environment CLI, powered by Homebrew.

## Runtime Requirements

1. Homebrew
1. ZSH

## Quick Example

1. `go run . --workdir testdata/project_ruby setup`
1. `go run . --workdir testdata/project_ruby shell`
1. `ruby version.rb`

## Supported Languages

### Ruby

1. `steel --workdir testdata/project_ruby setup`
1. `steel --workdir testdata/project_ruby shell`
1. `ruby version.rb`

### Go

1. `steel --workdir testdata/project_go setup`
1. `steel --workdir testdata/project_go shell`
1. `go run .`

## Testing

`go test -v`
