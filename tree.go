package browscap_go

import (
	"bytes"
	"sort"
)

type ExpressionTree struct {
	root *node
}

func NewExpressionTree() *ExpressionTree {
	return &ExpressionTree{
		root: &node{},
	}
}

func (r *ExpressionTree) Find(userAgent string) string {
	res, _ := r.root.findBest([]byte(userAgent), 0)
	return res
}

func (r *ExpressionTree) Add(name string) {
	exp := CompileExpression(bytes.ToLower([]byte(name)))

	last := r.root
	for _, e := range exp {
		shard := e.Shard()

		var found *node
		for _, node := range last.nodesPure[shard] {
			if node.token.Equal(e) {
				found = node
				break
			}
		}
		if found == nil {
			for _, node := range last.nodesFuzzy {
				if node.token.Equal(e) {
					found = node
					break
				}
			}
		}
		if found == nil {
			found = &node{
				token:  e,
				parent: last,
			}
			if e.Fuzzy() {
				last.nodesFuzzy = append(last.nodesFuzzy, found)
				sort.Sort(sort.Reverse(last.nodesFuzzy))
			} else {
				if last.nodesPure == nil {
					last.nodesPure = map[byte]nodes{}
				}
				last.nodesPure[shard] = append(last.nodesPure[shard], found)
				sort.Sort(sort.Reverse(last.nodesPure[shard]))
			}
		}
		last = found
	}

	last.exp = exp
	last.name = name

	score := len(name)
	for last != nil {
		if score > last.topScore {
			last.topScore = score
		}
		last = last.parent
	}
}

type node struct {
	nodesPure  map[byte]nodes
	nodesFuzzy nodes

	token *Token

	name string
	exp  Expression

	topScore int

	parent *node
}

func (n *node) findBest(s []byte, minScore int) (res string, maxScore int) {
	if n.topScore < minScore {
		return "", minScore
	}

	match := false
	if n.token != nil {
		match, s = n.token.MatchOne(s)
		if !match {
			return "", minScore
		}

		if n.name != "" {
			res = n.name
			minScore = n.topScore
		}
	}

	if len(s) > 0 {
		for _, nd := range n.nodesPure[s[0]] {
			r, ms := nd.findBest(s, minScore)
			if r != "" {
				if ms > minScore {
					res = r
					minScore = ms
				} else {
					break
				}
			}
		}

		for _, nd := range n.nodesFuzzy {
			r, ms := nd.findBest(s, minScore)
			if r != "" {
				if ms > minScore {
					res = r
					minScore = ms
				} else {
					break
				}
			}
		}
	}

	return res, minScore
}

type nodes []*node

func (n nodes) Len() int {
	return len(n)
}

func (n nodes) Less(i, j int) bool {
	return n[i].topScore < n[j].topScore
}

func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
