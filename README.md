# cork

An efficient binary serialisation format for Go (Golang).

[![](https://img.shields.io/circleci/token/cd14c319e747f26df1403e8e089a4c48fea767a1/project/abcum/cork/master.svg?style=flat-square)](https://circleci.com/gh/abcum/cork) [![](https://img.shields.io/badge/status-alpha-ff00bb.svg?style=flat-square)](https://github.com/abcum/cork) [![](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abcum/cork) [![](https://goreportcard.com/badge/github.com/abcum/cork?style=flat-square)](https://goreportcard.com/report/github.com/abcum/cork) [![](https://img.shields.io/badge/license-Apache_License_2.0-00bfff.svg?style=flat-square)](https://github.com/abcum/cork) 

#### Features

- Simple and efficient encoding
- Based on MsgPack serialization algorithm
- Stores go type information inside encoded data
- Faster serialization than gob binary encoding
- More efficient output data size than gob binary encoding
- Serializes native go types, and arbritary structs or interfaces
- Enables predetermined encoding for structs without run-time reflection
- Allows serialization to and from maps, structs, slices, and nil interfaces

#### Installation

```bash
go get github.com/abcum/cork
```
