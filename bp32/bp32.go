/*
 * Copyright (c) 2013 Zhen, LLC. http://zhen.io. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license.
 *
 */

package bp32

import (
	"errors"
	"github.com/reducedb/encoding"
	"github.com/reducedb/encoding/bitpacking"
)

const (
	DefaultBlockSize uint32 = 128
	DefaultPageSize uint32  = 65536
)

type BP32 struct {

}

var _ encoding.Integer = (*BP32)(nil)

func NewBP32() encoding.Integer {
	return &BP32{}
}

func (this *BP32) Compress(in []uint32, inpos *encoding.Cursor, inlength int, out []uint32, outpos *encoding.Cursor) error {
	//log.Printf("bp32/Compress: before inlength = %d\n", inlength)

	inlength = int(encoding.FloorBy(uint32(inlength), DefaultBlockSize))

	if inlength == 0 {
		return errors.New("BP32/Compress: block size less than 128. No work done.")
	}

	//log.Printf("bp32/Compress: after inlength = %d, len(in) = %d\n", inlength, len(in))

	out[outpos.Get()] = uint32(inlength)
	outpos.Increment()

	tmpoutpos := outpos.Get()

	for s := inpos.Get(); s < inpos.Get() + inlength; s += 32*4 {
		mbits1 := encoding.MaxBits(in, s, 32)
		mbits2 := encoding.MaxBits(in, s + 32, 32)
		mbits3 := encoding.MaxBits(in, s + 2*32, 32)
		mbits4 := encoding.MaxBits(in, s + 3*32, 32)

		//log.Printf("bp32/Compress: tmpoutpos = %d, s = %d\n", tmpoutpos, s)

		out[tmpoutpos] = (mbits1<<24) | (mbits2<<16) | (mbits3<<8) | mbits4
		tmpoutpos += 1

		//log.Printf("bp32/Compress: mbits1 = %d, mbits2 = %d, mbits3 = %d, mbits4 = %d, s = %d\n", mbits1, mbits2, mbits3, mbits4, out[tmpoutpos-1])

		bitpacking.FastPackWithoutMask(in, s, out, tmpoutpos, int(mbits1))
		//encoding.PrintUint32sInBits(in, s, 32)
		//encoding.PrintUint32sInBits(out, tmpoutpos, int(mbits1))
		tmpoutpos += int(mbits1)

		bitpacking.FastPackWithoutMask(in, s + 32, out, tmpoutpos, int(mbits2))
		//encoding.PrintUint32sInBits(in, s+32, 32)
		//encoding.PrintUint32sInBits(out, tmpoutpos, int(mbits2))
		tmpoutpos += int(mbits2)

		bitpacking.FastPackWithoutMask(in, s + 2*32, out, tmpoutpos, int(mbits3))
		//encoding.PrintUint32sInBits(in, s+2*32, 32)
		//encoding.PrintUint32sInBits(out, tmpoutpos, int(mbits3))
		tmpoutpos += int(mbits3)

		bitpacking.FastPackWithoutMask(in, s + 3*32, out, tmpoutpos, int(mbits4))
		//encoding.PrintUint32sInBits(in, s+3*32, 32)
		//encoding.PrintUint32sInBits(out, tmpoutpos, int(mbits4))
		tmpoutpos += int(mbits4)
	}

	inpos.Add(inlength)
	outpos.Set(tmpoutpos)

	return nil
}

func (this *BP32) Uncompress(in []uint32, inpos *encoding.Cursor, inlength int, out []uint32, outpos *encoding.Cursor) error {
	if inlength == 0 {
		return errors.New("BP32/Uncompress: Length is 0. No work done.")
	}

	outlength := in[inpos.Get()]
	inpos.Increment()

	tmpinpos := inpos.Get()

	//log.Printf("bp32/Uncompress: outlength = %d, inpos = %d, outpos = %d\n", outlength, inpos.Get(), outpos.Get())
	for s := outpos.Get(); s < outpos.Get() + int(outlength); s += 32*4 {
		mbits1 := in[tmpinpos]>>24
		mbits2 := (in[tmpinpos]>>16) & 0xFF
		mbits3 := (in[tmpinpos]>>8) & 0xFF
		mbits4 := (in[tmpinpos]) & 0xFF

		//log.Printf("bp32/Uncopmress: mbits1 = %d, mbits2 = %d, mbits3 = %d, mbits4 = %d, s = %d\n", mbits1, mbits2, mbits3, mbits4, s)
		tmpinpos += 1

		bitpacking.FastUnpack(in, tmpinpos, out, s, int(mbits1))
		tmpinpos += int(mbits1)
		//log.Printf("bp32/Uncompress: out = %v\n", out)

		bitpacking.FastUnpack(in, tmpinpos, out, s + 32, int(mbits2))
		tmpinpos += int(mbits2)
		//log.Printf("bp32/Uncompress: out = %v\n", out)

		bitpacking.FastUnpack(in, tmpinpos, out, s + 2*32, int(mbits3))
		tmpinpos += int(mbits3)
		//log.Printf("bp32/Uncompress: out = %v\n", out)

		bitpacking.FastUnpack(in, tmpinpos, out, s + 3*32, int(mbits4))
		tmpinpos += int(mbits4)

		//log.Printf("bp32/Uncompress: out = %v\n", out)
	}

	outpos.Add(int(outlength))
	inpos.Set(tmpinpos)

	return nil
}
