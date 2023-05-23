module github.com/blackducksoftware/hub-client-go

go 1.17

require (
	github.com/h2non/gock v1.0.14
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.3.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/h2non/parth v0.0.0-20190131123155-b4df798d6542 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14
