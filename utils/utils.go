package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseIP(p string) ([4]byte, error) {
	split := strings.Split(p, ".")
	if len(split) != 4 {
		return [4]byte{}, fmt.Errorf("not ep")
	}
	cp := [4]byte{}
	for i, v := range split {
		atoi, err := strconv.Atoi(v)
		if err != nil {
			return [4]byte{}, err
		}
		cp[i] = uint8(atoi)
	}

	return cp, nil
}

