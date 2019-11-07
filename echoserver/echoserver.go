package echoserver

import "github.com/summerwind/h2spec/spec"

var key = "echoserver"

func NewTestGroup(section string, name string) *spec.TestGroup {
	return &spec.TestGroup{
		Key:     key,
		Section: section,
		Name:    name,
	}
}

func Spec() *spec.TestGroup {
	tg := &spec.TestGroup{
		Key:  key,
		Name: "ECHOSERVER: Echo server tests for HTTP/2",
	}

	tg.AddTestGroup(EchoStreamId())

	return tg
}
