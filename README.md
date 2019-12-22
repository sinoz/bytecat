# ByteCat

An implementation of a byte string which is an immutable rope-like sequence designed for efficient concatenation of byte arrays.

## Installation

```
go get -v github.com/sinoz/bytecat
```

## Defining custom operators

You may would want to define your own operators for reading or writing bytes. To do this, it is advised to define a type alias for `bytecat.Iterator` / `bytecat.Builder`:

```
type MyIterator bytecat.Iterator

func (i MyIterator) ReadByteThenAdd16() (int, error) {
    value, err := i.ReadByte()
    if err != nil {
        return 0, err
    }

    return value + 16, nil
}

func readMyData(bs bytecat.String) {
    itr := MyIterator(bs.Iterator())
    fmt.Println(itr.ReadByteThenAdd16())
}
```
