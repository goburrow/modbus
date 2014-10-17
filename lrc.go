package modbus

type LRC struct {
	sum byte
}

func (lrc *LRC) pushByte(b byte) *LRC {
	lrc.sum += b
	return lrc
}

func (lrc *LRC) pushBytes(data []byte) *LRC {
	for _, b := range data {
		lrc.pushByte(b)
	}
	return lrc
}

func (lrc *LRC) value() byte {
	return -lrc.sum
}
