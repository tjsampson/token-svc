package datahelpers

// BytesToMb converts the input of bytes to MegaBytes
func BytesToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
