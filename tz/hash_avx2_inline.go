package tz

import (
	"hash"

	"github.com/nspcc-dev/tzhash/gf127"
)

type digest3 struct {
	x [2]gf127.GF127x2
}

var _ hash.Hash = (*digest3)(nil)

func NewAVX2Inline() hash.Hash {
	d := new(digest3)
	d.Reset()
	return d
}

func (d *digest3) Write(data []byte) (n int, err error) {
	n = len(data)
	for _, b := range data {
		mulByteRightx2(&d.x[0], &d.x[1], b)
	}
	return
}

func (d *digest3) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	h := d0.checkSum()
	return append(in, h[:]...)
}
func (d *digest3) Reset() {
	d.x[0] = gf127.GF127x2{1, 0, 0, 0}
	d.x[1] = gf127.GF127x2{0, 0, 1, 0}
}
func (d *digest3) Size() int      { return hashSize }
func (d *digest3) BlockSize() int { return hashBlockSize }
func (d *digest3) checkSum() (b [hashSize]byte) {
	// Matrix is stored transposed,
	// but we need to use order consistent with digest.
	h := d.x[0].ByteArray()
	copy(b[:], h[:16])
	copy(b[32:], h[16:])

	h = d.x[1].ByteArray()
	copy(b[16:], h[:16])
	copy(b[48:], h[16:])
	return
}

func mulByteRightx2(c00c10 *gf127.GF127x2, c01c11 *gf127.GF127x2, b byte)
