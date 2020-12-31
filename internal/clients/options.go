package clients

type BulkFetch struct {
	All bool
	Sort bool
	Source uint64
	Next uint64
}

type FuzzySearch struct {
	Term string
	Source uint64
	Next uint64
	Sort bool
}

func FuzzySearchOptions() *FuzzySearch {
	return &FuzzySearch{
		Term: "",
		Source: 0,
		Next: 10,
		Sort: false,
	}
}

func BulkFetchOptions() *BulkFetch {
	return &BulkFetch{
		All: false,
		Sort: false,
		Source: 0,
		Next: 10,
	}
}

func (opts *BulkFetch) FetchAll() *BulkFetch {
	opts.All = true
	return opts
}

func (opts *BulkFetch) FetchSorted() *BulkFetch {
	opts.Sort = true
	return opts
}

func (opts *BulkFetch) SetSource(source uint64) *BulkFetch {
	opts.Source = source
	return opts
}

func (opts *BulkFetch) SetNext(next uint64) *BulkFetch {
	opts.Next = next
	return opts
}

