package cipher_client

import (
	"bufio"
	"encoding/json"
	"os"
)

// Producer the struct for read/write data in file
type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

// NewProducer create producer
func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// WriteFile write data
func (p *Producer) WriteFile(userID string, key []byte) error { // key int value string договорились на о2о

	type tmp struct {
		UID string `json:"UID"`
		Key []byte `json:"Key"`
	}

	t := tmp{
		UID: userID,
		Key: key,
	}

	data, err := json.Marshal(&t)
	if err != nil {
		return err
	}
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}
	return p.writer.Flush()
}

// Close func for close file
func (p *Producer) Close() error {
	return p.file.Close()
}
