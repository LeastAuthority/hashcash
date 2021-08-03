package hashcash

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"time"
	"strconv"
)

type Stamp struct {
	Version  int
	Bits     uint
	Date     string
	Resource string
	Rand     string
	Counter  string
}

func (stamp Stamp) String() string {
	return fmt.Sprintf("%d:%d:%s:%s::%s:%s", stamp.Version, stamp.Bits, stamp.Date, stamp.Resource, stamp.Rand, stamp.Counter)
}

func Mint(bits uint, resource string) (string, error) {
	randBits := make([]byte, 12) // 96-bits of random data
	counterBits := make([]byte, 8) // for counter

	if bits > (sha1.Size * 8) {
		return "", fmt.Errorf("number of bits should be ≤ %d", sha1.Size * 8)
	}

	_, err := rand.Read(randBits)
	if err != nil {
		return "", err
	}
	randString := base64.StdEncoding.EncodeToString(randBits)

	_, err = rand.Read(counterBits)
	if err != nil {
		return "", err
	}
	counter := binary.BigEndian.Uint64(counterBits)
	// had to look up the source code to understand the format
	// string to be given. https://golang.org/src/time/format.go
	timestamp := time.Now().Format("060102")
	for true {
		countString := strconv.Itoa(int(counter))
		attempt := Stamp{
			Version:  1,
			Bits:     bits,
			Date:     timestamp,
			Resource: resource,
			Rand:     randString,
			Counter:  countString,
		}
		if Valid(attempt.String(), bits) {
			return attempt.String(), nil
		}
		counter += 1
	}

	return "", fmt.Errorf("could not mint a stamp for %d bits and resource \"%s\"", bits, resource)
}

func Valid(stamp string, bits uint) bool {
	buffer := bytes.NewBufferString(stamp)
	hash := sha1.New()
	sha1sum := hash.Sum(buffer.Bytes())

	n := countLeadingZeros(sha1sum)
	return (n >= bits)
}

func countLeadingZeros(buf []byte) uint {
	var zCount uint
	for _, b := range buf {
		if b == 0 {
			zCount += 8
		} else {
			var mask byte
			mask = 1 << 7
			for i := 0; i < 8; i++ {
				if (byte(b) & mask) != 0 {
					return (zCount + uint(i))
				}
				mask = mask >> 1
			}
		}
	}

	return zCount
}
