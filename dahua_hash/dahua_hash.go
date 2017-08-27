package dahua_hash

import (
	"crypto/md5"
)

// Original implementation: https://github.com/haicen/DahuaHashCreator

// DahuaHash implements the Dahua password hashing algorithm found in some Chinese DVRs, NVRs, and maybe other camera-related hardware
func DahuaHash(passw string) string {
	digest := md5.Sum([]byte(passw))
	return compressor(digest)
}

// compressor compresses a MD5 digest into the alphanumeric character space
// the result is always 8 characters long
func compressor(digest [md5.Size]byte) string {
	out := make([]byte, 0, md5.Size/2)

	for i := 0; i < len(digest); i += 2 {
		b1, b2 := int(digest[i]), int(digest[i+1]) // stretch to int to prevent wraparound in addition on next line
		compressed := byte((b1 + b2) % 62)
		if compressed < 10 {
			// 0..9 => 48..57 aka '0'..'9'
			compressed += 48
		} else if compressed < 36 {
			// 10..35 => 65..90 aka 'A'..'Z'
			compressed += 55
		} else {
			// 36..61 => 97..122 aka 'a'..'z'
			compressed += 61
		}
		out = append(out, compressed)
	}

	return string(out[:])
}
