# ByteCat

An implementation of a byte string which is an immutable rope-like sequence designed for efficient concatenation of byte arrays. It works similar to Go's slices that are built into the language. This implementation however, adds a higher-level API that will make your code slightly easier to read. Unlike regular slices, the subordinate Iterator-and Builder of the byte string also implement Go's Reader-and Writer interfaces and thus seamlessly integrate with other I/O packages.

## Installation

```
go get -v github.com/sinoz/bytecat
```

## Examples

There is a list of examples that you can consult [here](https://github.com/sinoz/bytecat/tree/master/examples).