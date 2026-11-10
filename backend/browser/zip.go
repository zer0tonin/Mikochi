package browser

import (
)

type ZipCompressor struct {}

func NewZipCompressor() *TarGzCompressor {
	return &ZipCompressor{}
}
