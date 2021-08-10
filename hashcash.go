package hashcash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type stamp struct {
	Version  int
	Bits     uint
	Date     string
	Resource string
	Rand     string
	Counter  string
}

func (stamp stamp) String() string {
	return fmt.Sprintf("%d:%d:%s:%s::%s:%s", stamp.Version, stamp.Bits, stamp.Date, stamp.Resource, stamp.Rand, stamp.Counter)
}

func Mint(bits uint, resource string) (string, error) {
	randBits := make([]byte, 12)   // 96-bits of random data
	counterBits := make([]byte, 8) // for counter

	if bits > (sha1.Size * 8) {
		return "", fmt.Errorf("number of bits should be ≤ %d", sha1.Size*8)
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
	for {
		countString := strconv.Itoa(int(counter))
		attempt := stamp{
			Version:  1,
			Bits:     bits,
			Date:     timestamp,
			Resource: resource,
			Rand:     randString,
			Counter:  countString,
		}
		if validatePartialHash(attempt.String(), bits) {
			return attempt.String(), nil
		}
		counter += 1
	}

	return "", fmt.Errorf("could not mint a stamp for %d bits and resource \"%s\"", bits, resource)
}

// XXX This is for the server side to evaluate the stamp. we are not
// looking at the timestamp for expiry. This function should take some
// kind of a duration in "days" as input (since the resolution of
// timestamp is in days (YYMMDD) and reject stamps that has expired.
func Evaluate(stamp string, requiredBits uint, resource string) bool {
	parts := strings.Split(stamp, ":")
	if len(parts) != 7 {
		return false
	}
	// stamp is of the form:
	// <ver>:<bits>:<timestamp>:<resource>::rand:counter
	ver := parts[0]
	bits, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return false
	}
	res := parts[3]

	return validatePartialHash(stamp, requiredBits) &&
		(res == resource) &&
		(ver == "1") &&
		(bits == uint64(requiredBits))
}

func validatePartialHash(stamp string, requiredBits uint) bool {
	buffer := bytes.NewBufferString(stamp)
	sha1sum := sha1.Sum(buffer.Bytes())

	actualBits := countLeadingZeros(sha1sum[:])
	return (actualBits >= requiredBits)
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
