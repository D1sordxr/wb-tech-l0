package reader

import "context"

type Reader struct {
	// TODO
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) Start(ctx context.Context) error {
	return nil
}

func (r *Reader) Stop(ctx context.Context) error {
	return nil
}
