package hubclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayChunkIterator_hasNext(t *testing.T) {
	type fields struct {
		chunks   []string
		position int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ArrayChunkIterator{
				chunks:   tt.fields.chunks,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, i.hasNext(), "hasNext()")
		})
	}
}

func TestArrayChunkIterator_next(t *testing.T) {
	type fields struct {
		chunks   []string
		position int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ArrayChunkIterator{
				chunks:   tt.fields.chunks,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, i.next(), "next()")
		})
	}
}

func TestNewArrayChunkIterator(t *testing.T) {
	type args struct {
		chunks []string
	}
	tests := []struct {
		name string
		args args
		want ArrayChunkIterator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewArrayChunkIterator(tt.args.chunks), "NewArrayChunkIterator(%v)", tt.args.chunks)
		})
	}
}
