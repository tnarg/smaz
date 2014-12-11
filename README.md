This is a pure Go implementation of [antirez's](https://github.com/antirez)
[SMAZ](https://github.com/antirez/smaz), a library for compressing short,
English strings.

## Installation

    $ go get github.com/kjk/smaz

## Usage

``` go
import (
  "github.com/kjk/smaz"
)

func main() {
  s := "Now is the time for all good men to come to the aid of the party."
  compressed := smaz.Encode(nil, []byte(s))
  decompressed, err := smaz.Decode(nil, compressed)
  if err != nil {
    fmt.Printf("decompressed: %s\n", string(decompressed))
    ...
}
```

Full [API documentation](http://godoc.org/github.com/kjk/smaz).

## Notes

This is not a direct port of the C version. It is not guaranteed that the output
of `smaz.Encode` will be precisely the same as the C library. However, the
output should be decompressible by the C library, and the output of the C
library should be decompressible by `smaz.Decode`.

## Author

[Salvatore Sanfilippo](https://github.com/antirez) designed SMAZ and wrote
[C implementation]](https://github.com/antirez/smaz).

[Caleb Spare](https://github.com/cespare) wrote initial
[Go port](https://github.com/cespare/go-smaz).

[Krzysztof Kowalczyk](http://blog.kowalczyk.info) improved speed of
decompression (2.4x faster) and compression (1.3x faster).

## Contributors

[Antoine Grondin](https://github.com/aybabtme)

## License

MIT Licensed.

## Other implementations

* [The original C implementation](https://github.com/antirez/smaz)
* [Javascript](https://npmjs.org/package/smaz)
