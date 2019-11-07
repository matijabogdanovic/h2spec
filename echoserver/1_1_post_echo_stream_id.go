package echoserver

import (
	"errors"
	"fmt"
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
	"strconv"
)

func PostEchoStreamId() *spec.TestGroup {
	tg := NewTestGroup("1.1", "POST to an echo server")

	// Peer must reply to non sequential delivery of HEADER and DATA frames
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Send maxStream HEADER frames followed by maxStream DATA frames",
		Requirement: "The endpoint must reply with the same number of HEADER and DATA frames",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Skip this test case when SETTINGS_MAX_CONCURRENT_STREAMS is unlimited.
			maxStreams, ok := conn.Settings[http2.SettingMaxConcurrentStreams]
			if !ok {
				return spec.ErrSkipped
			}

			// POST maxStream headers
			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"
			var streamId uint32 = 1
			for i := uint32(0); i < maxStreams; i++ {
				hp1 := http2.HeadersFrameParam{
					StreamID:      streamId,
					EndStream:     false,
					EndHeaders:    true,
					BlockFragment: conn.EncodeHeaders(headers),
				}
				conn.WriteHeaders(hp1)
				streamId += 2
			}
			// POST maxStream DATA frames where each frame body is set to streamId corresponding to the stream being sent
			streamId = 1
			for i := uint32(0); i < maxStreams; i++ {
				s := fmt.Sprintf("%d", streamId)
				conn.WriteData(streamId, true, []byte(s))
				streamId += 2
			}

			// Count only responded HEADER & DATA frames. Ignore other frame types.
			var headersCount = uint32(0)
			var dataCount = uint32(0)
			for headersCount < maxStreams || dataCount < maxStreams {
				actual := conn.WaitEvent()

				switch event := actual.(type) {
				case spec.HeadersFrameEvent:
					headersCount += 1
				case spec.DataFrameEvent:

					dataCount += 1
					data := string(event.Data())
					// check if content of received data is equal to streamId
					if data != strconv.Itoa(int(event.StreamID)) {
						var message = fmt.Sprintf("received data body \"%s\" does not match streamId: %d", data, event.StreamID)
						return errors.New(message)
					}
				default:

				}
			}

			if headersCount != maxStreams {
				var message = fmt.Sprintf("number of received HEADER frames: %d does not match number of maxStreams: %d", headersCount, maxStreams)
				return errors.New(message)
			}

			if dataCount != maxStreams {
				var message = fmt.Sprintf("number of received DATA frames: %d does not match number of maxStreams: %d", headersCount, maxStreams)
				return errors.New(message)
			}

			return nil
		},
	})

	return tg
}
