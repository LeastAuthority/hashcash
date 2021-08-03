package hashcash

import (
	"testing"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// zero out n bits out of the given sha1sum slice and check for n
// leading zero bits.
func fixedLeadingZeros(sha1sum []byte, n uint) {
	// zero out leading n bits from randomly generated sha1sum
	count := uint(0)
	for i, _ := range sha1sum {
		if count == n {
			break
		}

		if n < 8 {
			mask := uint8((1 << uint(8 - n)) - 1)
			sha1sum[i] = sha1sum[i] & mask
			break
		} else {
			sha1sum[i] = 0
			count += 8
		}
	}
}

func TestLeadingBits(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1234) // to generate reproducible results

	properties := gopter.NewProperties(nil)
	properties.Property("random shasum byte array with a random number of leading zeros", prop.ForAll(
		func(input []byte, bits uint) bool {
			fixedLeadingZeros(input, bits)
			return leadingBits(input, bits)
		},
		gen.SliceOfN(20, gen.UInt8Range(1,255)),
		gen.UIntRange(0, 20*8),
	))
	properties.TestingRun(t)
}

func TestLeadingZeros(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1234) // to generate reproducible results

	properties := gopter.NewProperties(nil)
	properties.Property("random shasum byte array with a random number of leading zeros", prop.ForAll(
		func(input []byte, bits uint) bool {
			fixedLeadingZeros(input, bits)
			n := countLeadingZeros(input)
			return (n >= bits)
		},
		gen.SliceOfN(20, gen.UInt8Range(1,255)),
		gen.UIntRange(0, 20*8),
	))
	properties.TestingRun(t)
}

