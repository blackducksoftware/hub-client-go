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

import "github.com/pkg/errors"

type ArrayChunkIterator struct {
	chunks   []string
	position int
}

func NewArrayChunkIterator(chunks []string) ArrayChunkIterator {
	i := new(ArrayChunkIterator)
	i.chunks = chunks
	i.position = -1
	return *i
}

func (i *ArrayChunkIterator) HasNext() bool {
	if len(i.chunks) > (i.position + 1) {
		return true
	} else {
		return false
	}
}

// Next returns true and the next chunk if there is a next chunk. Returns false otherwise.
func (i *ArrayChunkIterator) Next() (string, error) {
	i.position++
	if i.position < len(i.chunks) {
		return i.chunks[i.position], nil
	}
	return "", errors.New("A new chunk was requested but no chunks exist.")
}
