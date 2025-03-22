package common

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CursorToTime converts a UUID cursor to a timestamp
func CursorToTime(cursor string) (*time.Time, error) {
	if cursor == "" {
		return nil, nil
	}

	uid, err := uuid.Parse(cursor)
	if err != nil {
		return nil, err
	}

	// Extract timestamp from UUID
	// This is a simple example using the high bits of the UUID to store a timestamp
	bytes := uid.String()[:8] // Use first 8 chars of the UUID string as hex representation
	var unixTime int64
	fmt.Sscanf(bytes, "%x", &unixTime) // Convert hex to int64
	t := time.Unix(unixTime, 0)
	return &t, nil
}

// TimeToCursor converts a timestamp to a UUID cursor
func TimeToCursor(t *time.Time) string {
	if t == nil {
		return ""
	}

	// Create a UUID where the first part encodes the timestamp
	// In production, you'd use a better algorithm that preserves ordering
	unixTime := t.Unix()
	namespaceBuf := make([]byte, 16)
	// Put the unix timestamp in the first 8 bytes
	for i := 0; i < 8; i++ {
		namespaceBuf[i] = byte(unixTime >> uint(8*(7-i)))
	}
	// Fill the rest with random bytes
	for i := 8; i < 16; i++ {
		namespaceBuf[i] = byte(time.Now().UnixNano() >> uint(i))
	}

	ns, err := uuid.FromBytes(namespaceBuf)
	if err != nil {
		// If error, fall back to a simple v4 UUID
		return uuid.New().String()
	}
	return ns.String()
}
