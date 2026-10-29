module github.com/infrahq/infra/test

go 1.27.3

replace github.com/infrahq/infra => ../

require (
	github.com/infrahq/infra v0.21.0
	go.yaml.in/yaml/v3 v3.0.5
	gotest.tools/v3 v3.5.2
)

require (
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/rs/zerolog v1.35.1 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/term v0.47.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
