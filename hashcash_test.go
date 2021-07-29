package hashcash

import (
	"reflect"
	"testing"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func fixedLeadingZeros(sha1sum []byte, n uint) bool {
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
	// TODO: turn on nth bit?
	return leadingBits(sha1sum[:], m)
}

func TestLeadingBits(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1234) // to generate reproducible results

	properties := gopter.NewProperties(nil)
	properties.Property("random shasum byte array with a known number of leading zeros", prop.ForAll(
		fixedLeadingZeros,
		gen.SliceOfN(20, gen.UInt8Range(0,255),
			reflect.TypeOf(uint8(0))).
			SuchThat(func(v interface{}) bool {
				return len(v.([]uint8)) > 0
			}),
		gen.UIntRange(0, 20*8),
	))
	properties.TestingRun(t)
}
