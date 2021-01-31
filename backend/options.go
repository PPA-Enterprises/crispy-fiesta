package PPA

type BulkFetchOptions struct {
	All bool
	Sort bool
	Source uint64
	Next uint64
}

func NewBulkFetchOptions() BulkFetchOptions {
	return BulkFetchOptions {
		All: false,
		Sort: false,
		Source: 0,
		Next: 10,
	}
}
