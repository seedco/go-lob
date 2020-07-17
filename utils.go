package lob

//String is a helper function for creating a pointer string
func String(val string) *string {
	return &val
}

//Bool is a helper function for creating a pointer boolean
func Bool(val bool) *bool {
	return &val
}

//Uint64 is a helper function for creating a pointer uint64
func Uint64(val uint64) *uint64 {
	return &val
}
