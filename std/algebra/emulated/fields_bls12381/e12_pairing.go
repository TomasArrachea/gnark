package fields_bls12381

// ExptHalf set z to x^(t/2) in E12 and return z
// const t/2 uint64 = 7566188111470821376 // negative
func (e Ext12) ExptHalf(x *E12) *E12 {
	// FixedExp computation is derived from the addition chain:
	//
	//	_10      = 2*1
	//	_11      = 1 + _10
	//	_1100    = _11 << 2
	//	_1101    = 1 + _1100
	//	_1101000 = _1101 << 3
	//	_1101001 = 1 + _1101000
	//	return     ((_1101001 << 9 + 1) << 32 + 1) << 15
	//
	// Operations: 62 squares 5 multiplies
	//
	// Generated by github.com/mmcloughlin/addchain v0.4.0.

	// Step 1: z = x^0x2
	z := e.Square(x)

	// Step 2: z = x^0x3
	z = e.Mul(x, z)

	z = e.CyclotomicSquare(z)
	z = e.CyclotomicSquare(z)

	// Step 5: z = x^0xd
	z = e.Mul(x, z)

	// Step 8: z = x^0x68
	z = e.NCycloSquareCompressed(z, 3)
	z = e.DecompressKarabina(z)

	// Step 9: z = x^0x69
	z = e.Mul(x, z)

	// Step 18: z = x^0xd200
	z = e.NCycloSquareCompressed(z, 9)
	z = e.DecompressKarabina(z)

	// Step 19: z = x^0xd201
	z = e.Mul(x, z)

	// Step 51: z = x^0xd20100000000
	z = e.NCycloSquareCompressed(z, 32)
	z = e.DecompressKarabina(z)

	// Step 52: z = x^0xd20100000001
	z = e.Mul(x, z)

	// Step 67: z = x^0x6900800000008000
	z = e.NCycloSquareCompressed(z, 15)
	z = e.DecompressKarabina(z)

	z = e.Conjugate(z) // because tAbsVal is negative

	return z
}

// Expt set z to xᵗ in E12 and return z
// const t uint64 = 15132376222941642752 // negative
func (e Ext12) Expt(x *E12) *E12 {
	z := e.ExptHalf(x)
	z = e.CyclotomicSquare(z)
	return z
}

// MulBy014 multiplies z by an E12 sparse element of the form
//
//	E12{
//		C0: E6{B0: 1, B1: c1, B2: 0},
//		C1: E6{B0: 0, B1: c4, B2: 0},
//	}
//
// TODO : correct MulByO14 and not 034
func (e *Ext12) MulBy014(z *E12, c1, c4 *E2) *E12 {

	a := z.C0
	b := z.C1

	b = *e.MulBy01(&b, c1, c4)
	one := e.Ext2.One()

	c1 = e.Ext2.Add(one, c1)
	d := e.Ext6.Add(&z.C0, &z.C1)
	d = e.MulBy01(d, c1, c4)

	zC1 := e.Ext6.Add(&a, &b)
	zC1 = e.Ext6.Neg(zC1)
	zC1 = e.Ext6.Add(zC1, d)
	zC0 := e.Ext6.MulByNonResidue(&b)
	zC0 = e.Ext6.Add(zC0, &a)

	return &E12{
		C0: *zC0,
		C1: *zC1,
	}
}
