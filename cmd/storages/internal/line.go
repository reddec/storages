package internal

import (
	"encoding/base64"
	"encoding/json"
	"io"
)

func NewStringJSONLine(out io.Writer) (*jsonCodec, error) {
	cdc := &jsonCodec{
		first: true,
		out:   out,
	}
	return cdc, cdc.init()
}

type jsonCodec struct {
	first bool
	out   io.Writer
}

func (j *jsonCodec) init() error {
	_, err := j.out.Write([]byte("["))
	return err
}

func (j *jsonCodec) Close() error {
	if j.first {
		_, err := j.out.Write([]byte("]"))
		return err
	} else {
		_, err := j.out.Write([]byte("\n]"))
		return err
	}
}

func (j *jsonCodec) Write(key []byte) error {
	var err error
	if !j.first {
		_, err = j.out.Write([]byte(",\n  "))
	} else {
		_, err = j.out.Write([]byte("\n  "))
		j.first = false
	}
	if err != nil {
		return err
	}
	data, err := json.Marshal(string(key))
	if err != nil {
		return err
	}
	_, err = j.out.Write(data)
	return err
}
func NewPlainLine(out io.Writer, sep byte, tail bool) *plainCodec {
	return &plainCodec{
		first: true,
		sep:   sep,
		tail:  tail,
		out:   out,
	}
}

type plainCodec struct {
	first bool
	sep   byte
	tail  bool
	out   io.Writer
}

func (p *plainCodec) Close() error {
	if p.first || !p.tail {
		return nil
	}
	_, err := p.out.Write([]byte{p.sep})
	return err
}

func (p *plainCodec) Write(line []byte) error {
	var err error
	if !p.first {
		_, err = p.out.Write([]byte{p.sep})
	}
	p.first = false
	if err != nil {
		return err
	}
	_, err = p.out.Write(line)
	return err
}

func NewBase64Line(out io.Writer) *base64Codec {
	return &base64Codec{out: out}
}

type base64Codec struct {
	out io.Writer
}

func (b *base64Codec) Close() error {
	return nil
}

func (b *base64Codec) Write(k []byte) error {
	line := base64.StdEncoding.EncodeToString(k)
	_, err := b.out.Write([]byte(line))
	if err != nil {
		return err
	}
	_, err = b.out.Write([]byte("\n"))
	return err
}
