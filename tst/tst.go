// Copyright 2013 Codehack. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package tst implements a Tenary Search Tree (TST) inspired by
// the description on DrDobb's Journal 1998 by Jon Bentley and
// Bob Sedgewick[1].
//
// 1- http://www.drdobbs.com/database/ternary-search-trees/184410528?pgno=1
//
package tst

import (
	"strings"
)

// TST is the search tree
type TST struct {
	// Size is the total number of words in the tree
	Size int

	// root is the root node of this search tree.
	// All insertions and retrievals are done from the root node.
	root *tnode
}

// New returns a pointer to a new TST tree
func New() *TST {
	return &TST{Size: 0, root: &tnode{}}
}

// Get returns the value stored matching string s in tree t, or nil.
func (t *TST) Get(s string) interface{} {
	if s == "" {
		return nil
	}
	n := t.root.get(s)
	if n == nil {
		return nil
	}
	return n.val
}

// Find returns true if string s is found in tree t, false otherwise
func (t *TST) Find(s string) bool {
	return t.Get(s) != nil
}

// Insert inserts string s into tree t and assigns its value val to it.
func (t *TST) Insert(s string, val interface{}) bool {
	if !t.Find(s) {
		t.Size++
	}
	return t.root.rinsert(s, val) != nil
}

// Match returns all words that match pattern p, or nil if none.
// Patterns are strings with wildcards '*' (match 0 or more chars) or
// underscores '_' (match 1 char)
//
// 	foo := tst.Match("can*") // matches "can" "cannon" "candy" ...
//
//		bar := tst.Match("r_n_")  // matches "runs" "rank" "rent" ...
//
// BUG: wildcard pattern matching not implemented yet
func (t *TST) Match(p string) []string {
	if p == "" {
		return nil
	}
	// if p has no matching, just find the string
	if !strings.ContainsAny(p, "*_") {
		if t.Find(p) {
			return []string{p}
		}
		return nil
	}
	// when matching for "word*" use prefix instead since it's cheaper
	if idx := strings.LastIndex(p, "*"); strings.Count(p, "*") == 1 && idx == len(p)-1 {
		return t.Prefix(p[:idx])
	}
	matches := []string{}
	t.root.rmatch(p, "", &matches)
	if len(matches) > 0 {
		return matches
	}
	return nil
}

// Prefix returns all words that start with prefix p, or nil if none.
func (t *TST) Prefix(p string) []string {
	if p == "" {
		return nil
	}
	n := t.root.get(p)
	if n == nil {
		return nil // nothing has this prefix
	}
	matches := []string{}
	if n.val != nil {
		matches = append(matches, p)
	}
	n.eqkid.rprefix(p, &matches)
	if len(matches) > 0 {
		return matches
	}
	return nil
}

/*****************************************************************
 * Private functions on nodes
 *****************************************************************/

type tnode struct {
	c                   byte        // split char
	lokid, eqkid, hikid *tnode      // lo, eq and hi kids
	val                 interface{} // node value
}

func (n *tnode) get(s string) *tnode {
	if s == "" {
		return nil
	}
	if n == nil {
		return nil
	}
	c := s[0]
	switch {
	case c < n.c:
		return n.lokid.get(s)
	case c > n.c:
		return n.hikid.get(s)
	case len(s)-1 > 0:
		return n.eqkid.get(s[1:])
	default:
		return n
	}
}

func (n *tnode) rinsert(s string, val interface{}) *tnode {
	if s == "" {
		return nil
	}
	c := s[0]
	if n == nil {
		n = &tnode{c: c}
	}
	switch {
	case c < n.c:
		n.lokid = n.lokid.rinsert(s, val)
	case c > n.c:
		n.hikid = n.hikid.rinsert(s, val)
	case len(s)-1 > 0:
		n.eqkid = n.eqkid.rinsert(s[1:], val)
	default:
		n.val = val
	}
	return n
}

// TODO: wildcard '*' pattern matching
func (n *tnode) rmatch(pat, prefix string, m *[]string) {
	if n == nil {
		return
	}
	c := pat[0]
	if c == '_' || c < n.c {
		n.lokid.rmatch(pat, prefix, m)
	}
	if c == '_' || c == n.c {
		if n.val != nil && len(pat)-1 == 0 {
			*m = append(*m, prefix+string(n.c))
		}
		if len(pat)-1 > 0 {
			n.eqkid.rmatch(pat[1:], prefix+string(n.c), m)
		}
	}
	if c == '_' || c > n.c {
		n.hikid.rmatch(pat, prefix, m)
	}
}

func (n *tnode) rprefix(prefix string, m *[]string) {
	if n == nil {
		return
	}
	n.lokid.rprefix(prefix, m)
	if n.val != nil {
		*m = append(*m, prefix+string(n.c))
	}
	n.eqkid.rprefix(prefix+string(n.c), m)
	n.hikid.rprefix(prefix, m)
}
