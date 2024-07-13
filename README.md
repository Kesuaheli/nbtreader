# nbtreader

Parse and show Minecrafts NBT files in the command-line.

## Installing

Clone repo and build `cli` with Go

```sh
go build -o nbtreader ./cli
```

PS: For Windows its recommended to add a `.exe` to the filename:

```bash
go build -o nbtreader ./cli
```

This produces an ready-to-use executable file named `nbtreader` (or `nbtreader.exe`). See [next Section](#usage) for using the command and its syntax.

You may want to add a symlink in your PATH to this file, so you can use this executable systemwide. In the following examples I call the command `nbtreader` to call the file.

## Usage

After installation the binary can be used simply by passing a filename as argument:

```sh
nbtreader files/test.dat
```

produces

```
{
  name: "Bananrama"
}
```

and

```sh
nbtreader files/bigtest.nbt
```

produces

```
{
  longTest: 9223372036854775807l,
  shortTest: 32767s,
  stringTest: "HELLO WORLD THIS IS A TEST STRING ÅÄÖ!",
  floatTest: 0.49823147f,
  intTest: 2147483647,
  nested compound test: {
    ham: {
      name: "Hampus",
      value: 0.75f
    },
    egg: {
      name: "Eggbert",
      value: 0.5f
    }
  },
  listTest (long): [11l, 12l, 13l, 14l, 15l],
  listTest (compound): [{
    name: "Compound tag #0",
    created-on: 1264099775885l
  }, {
    name: "Compound tag #1",
    created-on: 1264099775885l
  }],
  byteTest: 127b,
  byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...)): [B; 0b, 62b, <trimmed 996 values>, 6b, 48b],
  doubleTest: 0.4931287132182315d
}
```

If no file is given, nbtreader reads from stdin:

```sh
# e.g. with pipe
cat files/test.nbt | nbtreader
# or redirect
nbreader < files/test.nbt
```

### Flags

Available flags are

- `-inType <string>`
- `-outType <string>`
- `-out <string>`

#### Flag `inType` and `outType`

With theese flags you can specify the in- and output type.

Current valid values for `inType`:
- `NBT` *(default if ommited)*

Current valid values for `outType`:
- `JSON`
- `NJSON` *([see spec](https://docs.google.com/document/d/1efDB9wyMLU4uWPTGY_nWNxBviS85iuicB8251kGiu2k/edit?usp=drivesdk))*
- `SNBT` *(default if ommited)*

Example:

```sh
nbtreader -outType JSON files/test.nbt
```

produces

```json
{
	"name": "Bananrama"
}
```

#### Flag `out`

To easily save the result in a file you could once again use redirecting:

```sh
nbtreader files/test.nbt > files output.txt
```

However, if you pass in the optional flag `-out <filename>` the output is instead written to the given filename:

```sh
nbtreader -out files/output.txt files/test.nbt
```
