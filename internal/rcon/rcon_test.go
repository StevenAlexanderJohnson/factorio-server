package rcon

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"testing"
)

func startMockRCONServer(t *testing.T, password string) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			var length int32
			if err := binary.Read(conn, binary.LittleEndian, &length); err != nil {
				return
			}
			packetData := make([]byte, length)
			if _, err := io.ReadFull(conn, packetData); err != nil {
				return
			}

			id := int32(binary.LittleEndian.Uint32(packetData[0:4]))
			packetType := int32(binary.LittleEndian.Uint32(packetData[4:8]))
			bodyBytes := packetData[8:]
			if nullIdx := bytes.IndexByte(bodyBytes, 0x00); nullIdx != -1 {
				bodyBytes = bodyBytes[:nullIdx]
			}
			body := string(bodyBytes)

			if packetType == SERVERDATA_AUTH {
				respID := id
				if body != password {
					respID = -1
				}
				// Send SERVERDATA_AUTH_RESPONSE
				buf := new(bytes.Buffer)
				respLen := int32(10)
				_ = binary.Write(buf, binary.LittleEndian, respLen)
				_ = binary.Write(buf, binary.LittleEndian, respID)
				_ = binary.Write(buf, binary.LittleEndian, SERVERDATA_AUTH_RESPONSE)
				buf.Write([]byte{0x00, 0x00})
				_, _ = conn.Write(buf.Bytes())
			} else if packetType == SERVERDATA_EXECCOMMAND {
				respText := "response to: " + body
				respBody := []byte(respText)
				respLen := int32(4 + 4 + len(respBody) + 2)

				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, respLen)
				_ = binary.Write(buf, binary.LittleEndian, id)
				_ = binary.Write(buf, binary.LittleEndian, SERVERDATA_RESPONSE_VALUE)
				buf.Write(respBody)
				buf.Write([]byte{0x00, 0x00})
				_, _ = conn.Write(buf.Bytes())
			}
		}
	}()

	return listener.Addr().String()
}

func TestRCONClient(t *testing.T) {
	addr := startMockRCONServer(t, "secret123")

	resp, err := ExecuteCommand(addr, "secret123", "/version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "response to: /version"
	if resp != expected {
		t.Errorf("expected %q, got %q", expected, resp)
	}

	_, err = ExecuteCommand(addr, "wrongpassword", "/version")
	if err == nil {
		t.Errorf("expected auth failure error, got nil")
	}
}
