package bot

import (
	"bufio"
	"bytes"
	"github.com/jxo-me/netx/api/handler"
)

func ConvertJsonMsg(d any) (string, error) {
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)
	err := handler.Write(bio, d, "json")
	if err != nil {
		return "", err
	}
	err = bio.Flush()
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
