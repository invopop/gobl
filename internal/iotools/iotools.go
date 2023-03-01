// Package iotools helps with reading documents.
package iotools

import (
	"context"
	"io"
)

type cancelReader struct {
	ctx context.Context
	r   io.Reader
}

// CancelableReader wraps r such that when ctx is cancelled, Read will return
// an error immediately.
func CancelableReader(ctx context.Context, r io.Reader) io.Reader {
	return &cancelReader{
		ctx: ctx,
		r:   r,
	}
}

func (r *cancelReader) Read(p []byte) (int, error) {
	var c int
	var err error
	wait := make(chan struct{}, 1)
	go func() {
		c, err = r.r.Read(p)
		close(wait)
	}()
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case <-wait:
		return c, err
	}
}
