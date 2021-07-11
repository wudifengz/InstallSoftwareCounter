package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"net"
)

func sendData(host string, r iList) error {
	t, err := net.Dial("tcp", host+":8765")
	if err != nil {
		return err
	}
	defer t.Close()
	var bb bytes.Buffer
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	w := zlib.NewWriter(&bb)
	w.Write(b)
	w.Close()
	lb := bb.Len()
	n, err := t.Write(bb.Bytes())
	if err != nil {
		return err
	}
	if n != lb {
		return fmt.Errorf("[Error] - Must send %dB data, only send %dB", lb, n)
	}
	return nil
}
