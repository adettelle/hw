package sqlstorage

import "context"

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	// TODO
	return nil
}
