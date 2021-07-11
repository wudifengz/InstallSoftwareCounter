package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
)

func uncompress(b []byte) []byte {
	var bb bytes.Buffer
	r, _ := zlib.NewReader(bytes.NewReader(b))
	io.Copy(&bb, r)
	return bb.Bytes()
}
func recv(res chan<- iList) {
	l, err := net.Listen("tcp", ":8765")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer c.Close()
		go func(conn net.Conn) {
			raddr := conn.RemoteAddr()
			buf := make([]byte, 4096)
			lenb, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			var rs iList
			err = json.Unmarshal(uncompress(buf[:lenb]), &rs)
			rs["remote"] = strings.Split(raddr.String(), ":")[0]
			res <- rs
			conn.Close()
		}(c)
	}
}
