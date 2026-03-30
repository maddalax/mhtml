module github.com/maddalax/htmgo/cli/htmgo

go 1.23.0

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/google/uuid v1.6.0
	github.com/maddalax/htmgo/framework v1.0.7-0.20250703190716-06f01b3d7c1b
	github.com/maddalax/htmgo/tools/html-to-htmgo v0.0.0-20250703190716-06f01b3d7c1b
	github.com/stretchr/testify v1.9.0
	golang.org/x/mod v0.21.0
	golang.org/x/sys v0.26.0
	golang.org/x/tools v0.25.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

require (
	github.com/bmatcuk/doublestar/v4 v4.7.1
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/maddalax/htmgo/framework => ../../framework
