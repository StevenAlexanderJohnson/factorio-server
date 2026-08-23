package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	SERVERDATA_AUTH           int32 = 3
	SERVERDATA_AUTH_RESPONSE  int32 = 2
	SERVERDATA_EXECCOMMAND    int32 = 2
	SERVERDATA_RESPONSE_VALUE int32 = 0
)

type Client struct {
	address  string
	password string
	conn     net.Conn
}

func NewClient(address string, password string) *Client {
	return &Client{
		address:  address,
		password: password,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to rcon at %s: %w", c.address, err)
	}
	c.conn = conn

	if err := c.authenticate(); err != nil {
		_ = c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *Client) authenticate() error {
	reqID := int32(1)
	if err := c.writePacket(reqID, SERVERDATA_AUTH, c.password); err != nil {
		return fmt.Errorf("failed to send auth packet: %w", err)
	}

	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer c.conn.SetReadDeadline(time.Time{})

	for {
		respID, respType, _, err := c.readPacket()
		if err != nil {
			return fmt.Errorf("failed to read auth response: %w", err)
		}

		if respType == SERVERDATA_AUTH_RESPONSE {
			if respID == -1 {
				return errors.New("rcon authentication failed: invalid password")
			}
			return nil
		}
	}
}

func (c *Client) Execute(cmd string) (string, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return "", err
		}
	}

	reqID := int32(2)
	if err := c.writePacket(reqID, SERVERDATA_EXECCOMMAND, cmd); err != nil {
		_ = c.Close()
		return "", fmt.Errorf("failed to send command packet: %w", err)
	}

	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer c.conn.SetReadDeadline(time.Time{})

	respID, _, body, err := c.readPacket()
	if err != nil {
		_ = c.Close()
		return "", fmt.Errorf("failed to read command response: %w", err)
	}

	if respID != reqID {
		return body, fmt.Errorf("mismatched request id: expected %d, got %d", reqID, respID)
	}

	return body, nil
}

func (c *Client) writePacket(id int32, packetType int32, body string) error {
	bodyBytes := []byte(body)
	length := int32(4 + 4 + len(bodyBytes) + 2)

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, length)
	_ = binary.Write(buf, binary.LittleEndian, id)
	_ = binary.Write(buf, binary.LittleEndian, packetType)
	buf.Write(bodyBytes)
	buf.Write([]byte{0x00, 0x00})

	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *Client) readPacket() (int32, int32, string, error) {
	var length int32
	if err := binary.Read(c.conn, binary.LittleEndian, &length); err != nil {
		return 0, 0, "", err
	}

	if length < 8 || length > 100000 {
		return 0, 0, "", fmt.Errorf("invalid packet length: %d", length)
	}

	packetData := make([]byte, length)
	if _, err := io.ReadFull(c.conn, packetData); err != nil {
		return 0, 0, "", err
	}

	id := int32(binary.LittleEndian.Uint32(packetData[0:4]))
	packetType := int32(binary.LittleEndian.Uint32(packetData[4:8]))

	bodyBytes := packetData[8:]
	if nullIdx := bytes.IndexByte(bodyBytes, 0x00); nullIdx != -1 {
		bodyBytes = bodyBytes[:nullIdx]
	}

	return id, packetType, string(bodyBytes), nil
}

func ExecuteCommand(address string, password string, cmd string) (string, error) {
	client := NewClient(address, password)
	if err := client.Connect(); err != nil {
		return "", err
	}
	defer client.Close()
	return client.Execute(cmd)
}
