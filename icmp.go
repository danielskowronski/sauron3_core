// stolen from https://github.com/porjo/pingo2/blob/master/icmp.go

/*
The MIT License (MIT)

Copyright (c) 2014 Ian Bishop

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package sauron3

import (
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// non-privileged ping on Linux requires special sysctl setting:
//     sysctl -w net.ipv4.ping_group_range="0 65535"
//
// above code allows everyone to use raw_sockets (still limited by caps)
// See: http://stackoverflow.com/questions/8290046/icmp-sockets-linux/20105379#20105379
func Ping(hostname string) (reply bool, err error) {
	ipAddr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		return false, err
	}

	readDeadline := time.Now().Add(time.Duration(time.Millisecond * ICMPReadTimeout))
	writeDeadline := time.Now().Add(time.Duration(time.Millisecond * ICMPWriteTimeout))

	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return false, err
	}
	defer c.Close()

	if err = c.SetReadDeadline(readDeadline); err != nil {
		return false, err
	}
	if err = c.SetWriteDeadline(writeDeadline); err != nil {
		return false, err
	}

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return false, err
	}
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: ipAddr.IP}); err != nil {
		return false, err
	}

	rb := make([]byte, 1500)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		return false, err
	}
	rm, err := icmp.ParseMessage(1/*iana.ProtocolICMP*/, rb[:n])
	if err != nil {
		return false, err
	}

	if rm.Type == ipv4.ICMPTypeEchoReply {
		return true, nil
	}

	return false, nil
}
