package base

import (
	"bufio"
	"bytes"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
)

var DefaultNonUtfEncoding = simplifiedchinese.GBK

func TransformReaderEncoding(r io.Reader, encoding encoding.Encoding) (err error, newReader io.Reader) {
	var bomBytes []byte
	var ri interface{} = r
	if readerSeeker, ok := ri.(io.ReadSeeker); ok {
		bomBytes, err = bufio.NewReaderSize(readerSeeker, 1024).Peek(1024)
		if err != nil {
			return err, nil
		}
		readerSeeker.Seek(0, io.SeekStart)
		r = readerSeeker
	} else {
		contentBytes, e := ioutil.ReadAll(r)
		if e != nil {
			return e, nil
		}
		if len(contentBytes) > 1024 {
			bomBytes = contentBytes[0:1024]
		} else {
			bomBytes = contentBytes
		}
		r = bytes.NewReader(contentBytes)
	}
	_, name, _ := charset.DetermineEncoding(bomBytes, "")
	if name != "utf-8" {
		if encoding != nil {
			return nil, transform.NewReader(r, encoding.NewDecoder())
		} else {
			return nil, transform.NewReader(r, DefaultNonUtfEncoding.NewDecoder())
		}
	} else {
		return nil, r
		//system default encoding is utf-8
		//return nil, transform.NewReader(r,  unicode.UTF8.NewDecoder())
	}
}
