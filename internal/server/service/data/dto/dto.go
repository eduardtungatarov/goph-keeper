package dto

type Create struct {
	Type  string
	Data  []byte
	Title string
}

type Read struct {
	Id int
}

type ReadResult struct {
	Type string
	Data []byte
}
