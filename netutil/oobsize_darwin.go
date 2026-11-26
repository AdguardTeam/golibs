//go:build darwin

package netutil

// IPDstOOBSize is the required size of the control-message buffer for
// [net.UDPConn.ReadMsgUDP] to read the original destination on darwin.
const IPDstOOBSize = 32
