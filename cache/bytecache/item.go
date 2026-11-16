package bytecache

import "github.com/c2h5oh/datasize"

// item is an element of the byte cache and also a node of the least-recently
// used (LRU) list.
//
// The LRU list is a circular and doubly-linked list.  One item acts as the
// list head: it holds no key or value, it is never stored in the cache map, and
// it marks both ends of the list.  The list head's next field points to the
// least recently used item, and its prev field points to the most recently used
// item.
type item struct {
	// next is the next item in the least-recently used linked list.  It must
	// not be nil.
	next *item

	// prev is the previous item in the least-recently used linked list.  It
	// must not be nil.
	prev *item

	// key is the cache key data.  It is always nil for the list head.
	key []byte

	// value is the cached value.  It is always nil for the list head.
	value []byte
}

// newListHead returns a new empty list head.
func newListHead() (head *item) {
	head = &item{}
	head.next = head
	head.prev = head

	return head
}

// back returns the most recently used item of the list headed by it.  If the
// list is empty, back returns it itself.  res is never nil.
func (it *item) back() (res *item) {
	return it.prev
}

// front returns the least recently used item of the list headed by it.  If the
// list is empty, front returns it itself.  res is never nil.
func (it *item) front() (res *item) {
	return it.next
}

// insertAfter inserts it into the list immediately after at, which must not be
// nil and must already be in a list.  Passing the list head as at inserts it at
// the front; passing the most recently used item inserts it at the back.
func (it *item) insertAfter(at *item) {
	it.next, at.next.prev = at.next, it
	at.next, it.prev = it, at
}

// size returns the total size of the key and the value in the item.
func (it *item) size() (sz datasize.ByteSize) {
	return datasize.ByteSize(len(it.key) + len(it.value))
}

// unlink removes it from its list.  It should currently be in a list, and
// should not be the list head.  unlink must not be called twice on the same
// item.
func (it *item) unlink() {
	it.prev.next, it.next.prev = it.next, it.prev

	// Self-link the removed item to make an accidental double unlink actually
	// harmless and drop stray references to neighbors.
	it.next, it.prev = it, it
}
