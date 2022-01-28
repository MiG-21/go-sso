package shared

import "bytes"

type (
	DiscardWriter struct {
		buff [][]byte
	}
)

func (w *DiscardWriter) Write(b []byte) (int, error) {
	w.buff = append(w.buff, b)
	return len(b), nil
}

func (w *DiscardWriter) EventOccurred(n []byte) bool {
	for _, b := range w.buff {
		if bytes.Contains(b, n) {
			return true
		}
	}
	return false
}

func (w *DiscardWriter) Reset() {
	w.buff = nil
}
