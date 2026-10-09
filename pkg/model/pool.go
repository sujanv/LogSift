package model

import "sync"

// RecordPool provides an allocation-free pool for Record instances.
type RecordPool struct {
	pool sync.Pool
}

// NewRecordPool creates a new pool of reusable Record pointers.
func NewRecordPool() *RecordPool {
	return &RecordPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Record{
					Fields: make(map[string]interface{}, 16),
				}
			},
		},
	}
}

// DefaultRecordPool is a package-level singleton pool for records.
var DefaultRecordPool = NewRecordPool()

// Get retrieves a clean Record from the pool.
func (p *RecordPool) Get() *Record {
	r := p.pool.Get().(*Record)
	r.Reset()
	return r
}

// Put returns a Record to the pool after resetting its fields.
func (p *RecordPool) Put(r *Record) {
	if r == nil {
		return
	}
	r.Reset()
	p.pool.Put(r)
}
