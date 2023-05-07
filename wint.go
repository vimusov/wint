/*
   wint - An utility which waits for Internet connection being established.

   Copyright (C) 2023 Vadim Kuznetsov <vimusov@gmail.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	mrand "math/rand"
	"net"
	"os"
	"syscall"
	"time"
)

func doPing(seq int, ip string) bool {
	var sock = -1
	var err error

	sock, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM|syscall.O_CLOEXEC, syscall.IPPROTO_ICMP)
	if err != nil {
		fmt.Printf("%s: socket(): %s.\n", ip, err)
		return false
	}
	defer func() {
		if err = syscall.Close(sock); err != nil {
			fmt.Printf("%s: close(%d): %s.\n", ip, sock, err)
		}
	}()

	timeout := syscall.Timeval{Sec: 5}
	for _, opt := range []int{syscall.SO_SNDTIMEO, syscall.SO_RCVTIMEO} {
		err = syscall.SetsockoptTimeval(sock, syscall.SOL_SOCKET, opt, &timeout)
		if err != nil {
			fmt.Printf("%s: setsockopt(%d): %s.\n", ip, opt, err)
			return false
		}
	}

	outData := make([]byte, 56)
	if _, err = crand.Read(outData); err != nil {
		fmt.Printf("%s: rand(): %s.\n", ip, err)
		return false
	}

	outMsg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: outData,
		},
	}
	outPkt, err := outMsg.Marshal(nil)
	if err != nil {
		fmt.Printf("%s: pack(%v): %s.\n", ip, outMsg, err)
		return false
	}

	addr := syscall.SockaddrInet4{Addr: ([4]byte)(net.ParseIP(ip).To4())}
	if err = syscall.Sendto(sock, outPkt, 0, &addr); err != nil {
		fmt.Printf("%s: sendto(): %s.\n", ip, err)
		return false
	}

	inPkt := make([]byte, 1500)

	size, _, err := syscall.Recvfrom(sock, inPkt, 0)
	if err != nil {
		fmt.Printf("%s: recvfrom(): %s.\n", ip, err)
		return false
	}

	inMsg, err := icmp.ParseMessage(1, inPkt[:size])
	if err != nil {
		fmt.Printf("%s: unpack(%v): %s.\n", ip, inPkt[:size], err)
		return false
	}

	if inMsg.Type != ipv4.ICMPTypeEchoReply {
		fmt.Printf("%s: wrong type(%v): %s.\n", ip, inMsg.Type, err)
		return false
	}

	echoMsg, ok := inMsg.Body.(*icmp.Echo)
	if !ok {
		fmt.Printf("%s: wrong class(%v): %s.\n", ip, inMsg, err)
		return false
	}

	ok = bytes.Compare(echoMsg.Data, outData) == 0
	if ok {
		fmt.Printf("%s: alive.\n", ip)
	} else {
		fmt.Printf("%s: wrong payload.\n", ip)
	}

	return ok
}

func main() {
	time.AfterFunc(5*time.Minute, func() {
		fmt.Println("offline timeout is over.")
		os.Exit(1)
	})

	ips := []string{
		"1.0.0.1",
		"1.1.1.1",
		"208.67.220.220",
		"208.67.222.222",
		"4.2.2.2",
		"4.2.2.6",
		"8.8.4.4",
		"8.8.8.8",
	}
	count := len(ips)
	success := 0

	for seq := 0; success < 3; seq++ {
		if doPing(seq, ips[mrand.Intn(count)]) {
			success++
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Println("internet connection established.")
}
