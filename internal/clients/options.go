package clients

type BulkFetch struct {
	All bool
	Sorted bool
	Source uint
	Next uint
}

func BulkFetchOptions() *BulkFetch {
	return &BulkFetch{
		All: false,
		Sorted: false,
		Source: 0,
		Next: 10,
	}
}

func (opts *BulkFetch) FetchAll() *BulkFetch {
	opts.All = true
	return opts
}

func (opts *BulkFetch) FetchSorted() *BulkFetch {
	opts.Sorted = true
	return opts
}

func (opts *BulkFetch) SetSource(source uint) *BulkFetch {
	opts.Source = source
	return opts
}

func (opts *BulkFetch) SetNext(next uint) *BulkFetch {
	opts.Next = next
	return opts
}
