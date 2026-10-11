package tui

// BookmarkManager tracks pinned/bookmarked log row indices.
type BookmarkManager struct {
	bookmarks map[int]bool
}

// NewBookmarkManager creates a BookmarkManager.
func NewBookmarkManager() *BookmarkManager {
	return &BookmarkManager{bookmarks: make(map[int]bool)}
}

// Toggle flips bookmark state on a row.
func (bm *BookmarkManager) Toggle(idx int) {
	if bm.bookmarks[idx] {
		delete(bm.bookmarks, idx)
	} else {
		bm.bookmarks[idx] = true
	}
}

// IsBookmarked returns true if row is marked.
func (bm *BookmarkManager) IsBookmarked(idx int) bool {
	return bm.bookmarks[idx]
}

// Next returns the next bookmarked index after cur.
func (bm *BookmarkManager) Next(cur, maxRows int) int {
	for i := cur + 1; i < maxRows; i++ {
		if bm.bookmarks[i] {
			return i
		}
	}
	for i := 0; i <= cur; i++ {
		if bm.bookmarks[i] {
			return i
		}
	}
	return cur
}
