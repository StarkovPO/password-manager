package cipher_client

import (
	"bufio"
	"encoding/json"
	"os"
)

// Consumer the struct for read/write data in file
type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

// NewConsumer create consumer
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return &Consumer{file: file, scanner: bufio.NewScanner(file)}, nil
}

// ReadFile read data
func (c *Consumer) ReadFile() ([]byte, error) {

	type tmp struct {
		UID string `json:"UID"`
		Key []byte `json:"Key"`
	}

	var buf []tmp
	for c.scanner.Scan() {
		data := c.scanner.Bytes()
		tmp := tmp{}

		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return nil, err
		}
		buf = append(buf, tmp)
	}

	b, err := json.Marshal(buf)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// Close func for close file
func (c *Consumer) Close() error {
	return c.file.Close()
}
