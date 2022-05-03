// Copyright 2022 Synopsys, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hubclient

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayChunkIterator_HasNext(t *testing.T) {
	type fields struct {
		chunks   []string
		position int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test has next",
			fields: fields{
				chunks:   []string{"chunk1", "chunk2"},
				position: 0,
			},
			want: true,
		},
		{
			name: "Test has no next",
			fields: fields{
				chunks:   []string{"chunk1", "chunk2"},
				position: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ArrayChunkIterator{
				chunks:   tt.fields.chunks,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, i.HasNext(), "HasNext()")
		})
	}
}

func TestArrayChunkIterator_Next(t *testing.T) {
	type fields struct {
		chunks   []string
		position int
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test get existing next",
			fields: fields{
				chunks:   []string{"chunk1", "chunk2"},
				position: 0,
			},
			want:    "chunk2",
			wantErr: assert.NoError,
		},
		{
			name: "Test fail when there is no next",
			fields: fields{
				chunks:   []string{"chunk1", "chunk2"},
				position: 1,
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ArrayChunkIterator{
				chunks:   tt.fields.chunks,
				position: tt.fields.position,
			}
			got, err := i.Next()
			if !tt.wantErr(t, err, fmt.Sprintf("Next()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Next()")
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
		{
			name: "Create new ArrayChunkIterator",
			args: args{
				chunks: []string{"chunk1", "chunk2"},
			},
			want: ArrayChunkIterator{chunks: []string{"chunk1", "chunk2"}, position: -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewArrayChunkIterator(tt.args.chunks), "NewArrayChunkIterator(%v)", tt.args.chunks)
		})
	}
}
