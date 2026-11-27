//go:build linux

package netutil

// IPDstOOBSize is the required size of the control-message buffer for
// [net.UDPConn.ReadMsgUDP] to read the original destination on Linux.
const IPDstOOBSize = 40
