package toc

import "github.com/yuin/goldmark/ast"

const _defaultMarker = '*'

// RenderList renders a table of contents as a nested list with a sane,
// default configuration for the ListRenderer.
func RenderList(toc *TOC) ast.Node {
	return new(ListRenderer).Render(toc)
}

// ListRenderer builds a nested list from a table of contents.
//
// For example,
//
//  # Foo
//  ## Bar
//  ## Baz
//  # Qux
//
// Becomes,
//
//  - Foo
//    - Bar
//    - Baz
//  - Qux
type ListRenderer struct {
	// Marker for elements of the list, e.g. '-', '*', etc.
	//
	// Defaults to '*'.
	Marker byte
}

// Render renders the table of contents into Markdown.
func (r *ListRenderer) Render(toc *TOC) ast.Node {
	return r.renderItems(toc.Items)
}

func (r *ListRenderer) renderItems(items Items) ast.Node {
	if len(items) == 0 {
		return nil
	}

	mkr := r.Marker
	if mkr == 0 {
		mkr = _defaultMarker
	}

	list := ast.NewList(mkr)
	for _, item := range items {
		list.AppendChild(list, r.renderItem(item))
	}
	return list
}

func (r *ListRenderer) renderItem(n *Item) ast.Node {
	item := ast.NewListItem(0)

	if t := n.Title; len(t) > 0 {
		title := ast.NewString(t)
		title.SetRaw(true)
		if len(n.ID) > 0 {
			link := ast.NewLink()
			link.Destination = append([]byte("#"), n.ID...)
			link.AppendChild(link, title)
			item.AppendChild(item, link)
		} else {
			item.AppendChild(item, title)
		}
	}

	if items := r.renderItems(n.Items); items != nil {
		item.AppendChild(item, items)
	}

	return item
}
