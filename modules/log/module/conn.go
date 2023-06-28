// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package module

import (
	"bufio"
	"net"
)

// The Conn type represents a WebSocket connection.
type Conn struct {
	net.Conn
	readBuf *bufio.Reader
}

func (c *Conn) Read(p []byte) (int, error) {
	return c.readBuf.Read(p)
}
