package contentenc

// intraBlock identifies a part of a file block
type intraBlock struct {
	BlockNo uint64 // Block number in file
	Skip    uint64 // Offset into block plaintext
	Length  uint64 // Length of plaintext data in this block
	fs      *ContentEnc
}

// isPartial - is the block partial? This means we have to do read-modify-write.
func (ib *intraBlock) IsPartial() bool {
	if ib.Skip > 0 || ib.Length < ib.fs.plainBS {
		return true
	}
	return false
}

// BlockCipherOff returns the ciphertext offset corresponding to BlockNo
func (ib *intraBlock) BlockCipherOff() (offset uint64) {
	return ib.fs.BlockNoToCipherOff(ib.BlockNo)
}

// BlockPlainOff returns the plaintext offset corresponding to BlockNo
func (ib *intraBlock) BlockPlainOff() (offset uint64) {
	return ib.fs.BlockNoToPlainOff(ib.BlockNo)
}

// CropBlock - crop a potentially larger plaintext block down to the relevant part
func (ib *intraBlock) CropBlock(d []byte) []byte {
	lenHave := len(d)
	lenWant := int(ib.Skip + ib.Length)
	if lenHave < lenWant {
		return d[ib.Skip:lenHave]
	}
	return d[ib.Skip:lenWant]
}

// Ciphertext range corresponding to the sum of all "blocks" (complete blocks)
func (ib *intraBlock) JointCiphertextRange(blocks []intraBlock) (offset uint64, length uint64) {
	firstBlock := blocks[0]
	lastBlock := blocks[len(blocks)-1]

	offset = ib.fs.BlockNoToCipherOff(firstBlock.BlockNo)
	offsetLast := ib.fs.BlockNoToCipherOff(lastBlock.BlockNo)
	length = offsetLast + ib.fs.cipherBS - offset

	return offset, length
}
