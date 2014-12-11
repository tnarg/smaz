// Package smaz is a pure Go implementation of the smaz algorithm
// (https://github.com/antirez/smaz) for compressing small strings.
package smaz

import (
	"errors"

	"github.com/kjk/smaz/trie"
)

var codeStrings = []string{" ",
	"the", "e", "t", "a", "of", "o", "and", "i", "n", "s", "e ", "r", " th",
	" t", "in", "he", "th", "h", "he ", "to", "\r\n", "l", "s ", "d", " a", "an",
	"er", "c", " o", "d ", "on", " of", "re", "of ", "t ", ", ", "is", "u", "at",
	"   ", "n ", "or", "which", "f", "m", "as", "it", "that", "\n", "was", "en",
	"  ", " w", "es", " an", " i", "\r", "f ", "g", "p", "nd", " s", "nd ", "ed ",
	"w", "ed", "http://", "for", "te", "ing", "y ", "The", " c", "ti", "r ", "his",
	"st", " in", "ar", "nt", ",", " to", "y", "ng", " h", "with", "le", "al", "to ",
	"b", "ou", "be", "were", " b", "se", "o ", "ent", "ha", "ng ", "their", "\"",
	"hi", "from", " f", "in ", "de", "ion", "me", "v", ".", "ve", "all", "re ",
	"ri", "ro", "is ", "co", "f t", "are", "ea", ". ", "her", " m", "er ", " p",
	"es ", "by", "they", "di", "ra", "ic", "not", "s, ", "d t", "at ", "ce", "la",
	"h ", "ne", "as ", "tio", "on ", "n t", "io", "we", " a ", "om", ", a", "s o",
	"ur", "li", "ll", "ch", "had", "this", "e t", "g ", "e\r\n", " wh", "ere",
	" co", "e o", "a ", "us", " d", "ss", "\n\r\n", "\r\n\r", "=\"", " be", " e",
	"s a", "ma", "one", "t t", "or ", "but", "el", "so", "l ", "e s", "s,", "no",
	"ter", " wa", "iv", "ho", "e a", " r", "hat", "s t", "ns", "ch ", "wh", "tr",
	"ut", "/", "have", "ly ", "ta", " ha", " on", "tha", "-", " l", "ati", "en ",
	"pe", " re", "there", "ass", "si", " fo", "wa", "ec", "our", "who", "its", "z",
	"fo", "rs", ">", "ot", "un", "<", "im", "th ", "nc", "ate", "><", "ver", "ad",
	" we", "ly", "ee", " n", "id", " cl", "ac", "il", "</", "rt", " wi", "div",
	"e, ", " it", "whi", " ma", "ge", "x", "e c", "men", ".com",
}

var codes = make([][]byte, len(codeStrings))
var codeTrie = trie.New()

func init() {
	for i, code := range codeStrings {
		codes[i] = []byte(code)
		codeTrie.Put([]byte(code), i)
	}
}

func next(d []byte, n int) ([]byte, []byte) {
	if len(d) < n {
		return d, nil
	}
	return d[:n], d[n:]
}

func appendSrc(dst, src []byte) []byte {
	// We can write a max of 255 continuous verbatim characters, because the
	// length of the continous verbatim section is represented by a single byte.
	for len(src) > 0 {
		chunk, src = next(src, 255)
		if len(chunk) == 1 {
			// 254 is code for a single verbatim byte
			dst = append(dst, byte(254))
			dst = append(dst, chunk[0])
		} else {
			// 255 is code for a verbatim string. It is followed by a byte
			// containing the length of the string.
			dst = append(dst, byte(255))
			dst = append(dst, byte(len(chunk)))
			dst = append(dst, chunk...)
		}
	}
	return dst
}

// Encode returns the encoded form of src. The returned slice may be a sub-slice
// of dst if dst was large enough to hold the entire encoded block. Otherwise,
// a newly allocated slice will be returned. It is valid to pass a nil dst.
func Encode(dst, src []byte) []byte {
	dst = dst[:0]
	root := codeTrie.Root()
	currPos := 0
	nLeft := len(src)
	tmp := src
	code := 0
	prefixLen := 0
	for currPos < nLeft {
		node := root
		for i, c := range tmp {
			next := node.Walk(c)
			if next == nil {
				break
			}
			node = next
			if node.Terminal() {
				prefixLen = i + 1
				code = node.Val()
			}
		}
		if prefixLen == 0 {
			currPos++
			tmp = tmp[1:]
			continue
		}
		dst = appendSrc(dst, src[:currPos])
		dst = append(dst, byte(code))
		src = src[currPos+prefixLen:]
		currPos = 0
		nLeft = len(src)
		tmp = src
		prefixLen = 0
	}
	dst = appendSrc(dst, src)
	return dst
}

// ErrCorrupt reports that the input is invalid.
var ErrCorrupt = errors.New("smaz: corrupt input")

// Decode returns the decoded form of src. The returned slice may be a sub-slice
// of dst if dst was large enough to hold the entire decoded block. Otherwise,
// a newly allocated slice will be returned. It is valid to pass a nil dst.
func Decode(dst, src []byte) ([]byte, error) {
	if cap(dst) < len(src) {
		dst = make([]byte, 0, len(src))
	}
	for len(src) > 0 {
		n := int(src[0])
		switch n {
		case 254: // Verbatim byte
			if len(src) < 2 {
				return nil, ErrCorrupt
			}
			dst = append(dst, src[1])
			src = src[2:]
		case 255: // Verbatim string
			if len(src) < 2 {
				return nil, ErrCorrupt
			}
			n = int(src[1])
			if len(src) < n+2 {
				return nil, ErrCorrupt
			}
			dst = append(dst, src[2:n+2]...)
			src = src[n+2:]
		default: // Look up encoded value
			d := codes[n]
			dst = append(dst, d...)
			src = src[1:]
		}
	}

	return dst, nil
}
