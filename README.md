# paswitcher
paswitcher is a small Go-based program to switch between available Pulseaudio outputs. Intended to be used for easy output switching via keybind while keeping devices plugged in.

paswitcher depends on `pactl`, `pacmd`, and `amixer`

To run, just do `go run paswitcher.go`

or build a binary with `go build paswitcher.go`