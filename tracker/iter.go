package tracker

import (
	"context"
	"fmt"
	"iter"
)

// pageOptions is a pointer to ListOptions or to an options struct that embeds
// it.
type pageOptions[O any] interface {
	*O
	listOptions() *ListOptions
}

func (o *ListOptions) listOptions() *ListOptions { return o }

// cursorOptions is a pointer to the options of a list that pages by the ID of
// the last item it returned.
type cursorOptions[O any] interface {
	*O
	cursor() *string
}

func (o *ChangelogOptions) cursor() *string { return &o.ID }

func (o *CommentListOptions) cursor() *string { return &o.ID }

// pageIter yields every item of a page-numbered list, starting with the page
// list returns for opts. It copies opts when called and sets only the page of
// the copy. Page 0 sends no page parameter, which the server reads as page 1.
// The iteration ends after page X-Total-Pages or, when a response lacks that
// header, on the first empty page.
func pageIter[O any, P pageOptions[O], T any](ctx context.Context, opts P, list func(context.Context, P) ([]T, *Response, error)) iter.Seq2[T, error] {
	var base O
	if opts != nil {
		base = *opts
	}

	return func(yield func(T, error) bool) {
		page := P(&base).listOptions().Page
		for {
			o := base
			P(&o).listOptions().Page = page
			items, resp, err := list(ctx, &o)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			page = max(page, 1)
			if len(items) == 0 || (resp != nil && resp.TotalPages > 0 && page >= resp.TotalPages) {
				return
			}
			page++
		}
	}
}

// cursorIter yields every item of a cursor list, starting with the page list
// returns for opts and then from the ID of each page's last item. It copies
// opts when called and sets only the cursor of the copy.
//
// Tracker sends rel="next" from the last full page too, and X-Total-Pages only
// on some cursor lists, so only an empty page ends the iteration. A page whose
// last ID is empty or repeats the cursor would be fetched forever, so it ends
// the iteration with an error.
func cursorIter[O any, P cursorOptions[O], T any](ctx context.Context, opts P, id func(T) string, list func(context.Context, P) ([]T, *Response, error)) iter.Seq2[T, error] {
	var base O
	if opts != nil {
		base = *opts
	}

	return func(yield func(T, error) bool) {
		var zero T
		cursor := *P(&base).cursor()
		for {
			o := base
			*P(&o).cursor() = cursor
			items, _, err := list(ctx, &o)
			if err != nil {
				yield(zero, err)
				return
			}
			if len(items) == 0 {
				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			next := id(items[len(items)-1])
			if next == "" || next == cursor {
				yield(zero, fmt.Errorf("tracker: cannot page past cursor %q: the last item of the page has ID %q", cursor, next))
				return
			}
			cursor = next
		}
	}
}
