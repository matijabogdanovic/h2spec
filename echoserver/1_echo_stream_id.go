package echoserver

import "github.com/summerwind/h2spec/spec"

func EchoStreamId() *spec.TestGroup {
	tg := NewTestGroup("1", "Echo back a StreamId")

	tg.AddTestGroup(PostEchoStreamId())

	return tg
}
