package processbar

import "testing"

func TestUpdateMsg(_ *testing.T) {
	// TODO
	// Test these
	// 1 - Sending messages without starting to listen
	//   - messages fail after limit (tests trySendMsgToChannel())
	//   - blocking messages stuck forever - timeout after 0.5sec (Just do 1)
	//     -  Tests sendMsgToChannelBlocking()
	// 2 - Use SendAddProcessMsg() and verify that new process is added soon
	// 3 - Use SendUpdateProcessNameMsg() and verify update
	// 4 - Verify that stopListeningMsg works. Use SendStopListeningMsgBlocking()
	//     - Test m.IsListeningForUpdates()

}
