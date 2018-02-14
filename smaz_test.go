package smaz_test

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/tnarg/smaz"
)

var antirezTestStrings = []string{"",
	"This is a small string",
	"foobar",
	"the end",
	"not-a-g00d-Exampl333",
	"Smaz is a simple compression library",
	"Nothing is more difficult, and therefore more precious, than to be able to decide",
	"this is an example of what works very well with smaz",
	"1000 numbers 2000 will 10 20 30 compress very little",
	"and now a few italian sentences:",
	"Nel mezzo del cammin di nostra vita, mi ritrovai in una selva oscura",
	"Mi illumino di immenso",
	"L'autore di questa libreria vive in Sicilia",
	"try it against urls",
	"http://google.com",
	"http://programming.reddit.com",
	"http://github.com/antirez/smaz/tree/master",
	"https://github.com/TheCoolKids",
	"/media/hdb1/music/Alben/The Bla",
}

func TestCorrectness(t *testing.T) {
	// Set up our slice of test strings.
	inputs := make([][]byte, 0)
	for _, s := range antirezTestStrings {
		inputs = append(inputs, []byte(s))
	}
	// An array with every possible byte value in it.
	allBytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		allBytes[i] = byte(i)
	}
	inputs = append(inputs, allBytes)
	// A long array of all 0s (the longest continuous string that can be represented is 256; any longer than
	// this and the compressor will need to split it into chunks)
	allZeroes := make([]byte, 300)
	for i := 0; i < 300; i++ {
		allZeroes[i] = byte(0)
	}
	inputs = append(inputs, allZeroes)

	for _, input := range inputs {
		compressed := smaz.DefaultCodec.Encode(nil, input)
		decompressed, err := smaz.DefaultCodec.Decode(nil, compressed)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(input, decompressed) {
			t.Fatalf("want %q after decompression; got %q\n", input, decompressed)
		}

		if len(input) > 1 && len(input) < 50 {
			compressionLevel := 100 - ((100.0 * len(compressed)) / len(input))
			if compressionLevel < 0 {
				t.Logf("%q enlarged by %d%%\n", input, -compressionLevel)
			} else {
				t.Logf("%q compressed by %d%%\n", input, compressionLevel)
			}
		}
	}
}

func TestCustomTable(t *testing.T) {
	table := []string{
		"This is a small string",
		"foobar",
		"the end",
		"not-a-g00d-Exampl333",
		"http://google.com",
		"http://programming.reddit.com",
	}
	cdc := smaz.NewCodec(table)

	inputs := make([][]byte, 0)
	for _, s := range table {
		inputs = append(inputs, []byte(s))
	}

	for _, input := range inputs {
		compressed := cdc.Encode(nil, input)
		decompressed, err := cdc.Decode(nil, compressed)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(input, decompressed) {
			t.Fatalf("want %q after decompression; got %q\n", input, decompressed)
		}
		if len(compressed) > 1 {
			t.Fatalf("want len(encode(%q)) == 1 after compression; got %d\n", input, len(compressed))
		}

		if len(input) > 1 && len(input) < 50 {
			compressionLevel := 100 - ((100.0 * len(compressed)) / len(input))
			if compressionLevel < 0 {
				t.Logf("%q enlarged by %d%%\n", input, -compressionLevel)
			} else {
				t.Logf("%q compressed by %d%%\n", input, compressionLevel)
			}
		}
	}
}

func TestCorrectnessWithCustomTable(t *testing.T) {
	// Set up our slice of test strings.
	inputs := make([][]byte, 0)
	for _, s := range antirezTestStrings {
		inputs = append(inputs, []byte(s))
	}
	// An array with every possible byte value in it.
	allBytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		allBytes[i] = byte(i)
	}
	inputs = append(inputs, allBytes)
	// A long array of all 0s (the longest continuous string that can be represented is 256; any longer than
	// this and the compressor will need to split it into chunks)
	allZeroes := make([]byte, 300)
	for i := 0; i < 300; i++ {
		allZeroes[i] = byte(0)
	}
	inputs = append(inputs, allZeroes)

	cdc := smaz.NewCodec([]string{"http://", "is", ".com"})
	for _, input := range inputs {
		compressed := cdc.Encode(nil, input)
		decompressed, err := cdc.Decode(nil, compressed)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(input, decompressed) {
			t.Fatalf("want %q after decompression; got %q\n", input, decompressed)
		}

		if len(input) > 1 && len(input) < 50 {
			compressionLevel := 100 - ((100.0 * len(compressed)) / len(input))
			if compressionLevel < 0 {
				t.Logf("%q enlarged by %d%%\n", input, -compressionLevel)
			} else {
				t.Logf("%q compressed by %d%%\n", input, compressionLevel)
			}
		}
	}
}

func loadTestData(t testing.TB) ([][]byte, int64) {
	f, err := os.Open("./testdata/pg5200.txt")
	if err != nil {
		t.Fatal(err)
	}

	var lines [][]byte
	var n int64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := []byte(scanner.Text()) // Note that .Bytes() would require us to manually copy
		lines = append(lines, input)
		n += int64(len(input))
	}
	return lines, n
}

func BenchmarkCompression(b *testing.B) {
	b.StopTimer()
	inputs, n := loadTestData(b)
	b.SetBytes(n)
	b.StartTimer()
	var dst []byte
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			dst = smaz.DefaultCodec.Encode(dst, input)
		}
	}
}

func BenchmarkDecompression(b *testing.B) {
	b.StopTimer()
	inputs, _ := loadTestData(b)
	compressedStrings := make([][]byte, len(inputs))
	var n int64
	for i, input := range inputs {
		compressed := smaz.DefaultCodec.Encode(nil, input)
		compressedStrings[i] = compressed
		n += int64(len(compressed))
	}
	b.SetBytes(n)
	b.StartTimer()
	var dst []byte
	var err error
	for i := 0; i < b.N; i++ {
		for _, compressed := range compressedStrings {
			dst, err = smaz.DefaultCodec.Decode(dst, compressed)
			if err != nil {
				b.Fatalf("Decompress failed with %s", err)
			}
		}
	}
}
