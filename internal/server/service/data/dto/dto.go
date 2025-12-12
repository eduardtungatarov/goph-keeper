package dto

type Create struct {
	Type  string
	Data  []byte
	Title string
}

type Read struct {
	ID int
}

type ReadResult struct {
	Type string
	Data []byte
}

type Delete struct {
	ID int
}
