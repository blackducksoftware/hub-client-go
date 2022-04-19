package hubclient

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

func (i *ArrayChunkIterator) hasNext() bool {
	if len(i.chunks) < (i.position + 1) {
		return true
	} else {
		return false
	}
}

func (i *ArrayChunkIterator) next() string {
	i.position++
	return i.chunks[i.position]
}
