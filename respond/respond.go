package respond

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"net"

	"github.com/FreifunkBremen/yanic/data"
)

const (
	// default multicast group used by announced
	multicastAddressDefault = "ff05:0:0:0:0:0:2:1001"

	// default udp port used by announced
	port = 1001

	// maximum receivable size
	MaxDataGramSize = 8192
)

// Response of the respond request
type Response struct {
	Address *net.UDPAddr
	Raw     []byte
}

func NewRespone(res *data.ResponseData, addr *net.UDPAddr) (*Response, error) {
	buf := new(bytes.Buffer)
	flater, err := flate.NewWriter(buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	defer flater.Close()

	if err = json.NewEncoder(flater).Encode(res); err != nil {
		return nil, err
	}

	err = flater.Flush()

	return &Response{
		Raw:     buf.Bytes(),
		Address: addr,
	}, err
}

func (res *Response) parse() (*data.ResponseData, error) {
	// Deflate
	deflater := flate.NewReader(bytes.NewReader(res.Raw))
	defer deflater.Close()

	// Unmarshal
	rdata := &data.ResponseData{}
	err := json.NewDecoder(deflater).Decode(rdata)

	return rdata, err
}
