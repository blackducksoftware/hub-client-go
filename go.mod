module github.com/blackducksoftware/hub-client-go

go 1.15

require (
	github.com/h2non/gock v0.0.0-00010101000000-000000000000
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
)

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14
