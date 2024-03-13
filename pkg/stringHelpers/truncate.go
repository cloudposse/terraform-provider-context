package stringHelpers

import (
	"fmt"
	"strconv"

	"github.com/sigurn/crc16"
)

func TruncateWithHash(input string, maxLength int) string {
	// If the input is already within the maxLength, no need to truncate
	if len(input) <= maxLength {
		return input
	}

	data := []byte(input)

	// Choose a specific CRC16 polynomial from the predefined ones
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

	// Calculate the CRC16 checksum
	checksum := crc16.Checksum(data, table)
	hashLength := len(strconv.FormatUint(uint64(checksum), 10))

	truncated := fmt.Sprintf("%s%d", input[:maxLength-hashLength], checksum)

	return truncated
}
