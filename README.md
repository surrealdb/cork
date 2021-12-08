# cork

An efficient binary serialisation format for Go (Golang).

[![](https://img.shields.io/badge/status-1.0.0-ff00bb.svg?style=flat-square)](https://github.com/surrealdb/cork) [![](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/surrealdb/cork) [![](https://goreportcard.com/badge/github.com/surrealdb/cork?style=flat-square)](https://goreportcard.com/report/github.com/surrealdb/cork) [![](https://img.shields.io/badge/license-Apache_License_2.0-00bfff.svg?style=flat-square)](https://github.com/surrealdb/cork) 

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
go get github.com/surrealdb/cork
```
