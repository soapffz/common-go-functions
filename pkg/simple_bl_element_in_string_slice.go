package pkg

import (
	"github.com/zeromicro/go-zero/core/lang"
)

func SimpleBlElementInStringSlice(stringL []string) []string {
	// 字符串列表取最小交集，用于黑名单关键词去重精简,by soapffz 2022-11-22

	// 主要实现方式，使用使用[zeromicro/go-zero](https://github.com/zeromicro/go-zero)里面的组件[stringx](https://go-zero.dev/cn/docs/blog/tool/keywords/)实现
	// 遍历整个数组，每次选择一个元素，剩余元素构建为筛选器，将元素与筛选器进行匹配，如果匹配到，则将该元素从数组中删除，以此得到最小字符串列表

	for _, i := range stringL {
		// 深拷贝一份原始数组
		tmpCopiedL := make([]string, len(stringL))
		copy(tmpCopiedL, stringL)
		// 删除当前元素
		tmpStringSlice := RemoveSpecifiedElementInStringSlice(tmpCopiedL, i)
		keyWordsFilter := NewTrie(tmpStringSlice)
		// 如果当前元素包含剩余列表的关键词
		matchedKeyWordsL := keyWordsFilter.FindKeywords(i)
		if len(matchedKeyWordsL) > 0 {
			// 从原始列表中删除当前元素
			stringL = RemoveSpecifiedElementInStringSlice(stringL, i)
			// log.Println("[+] Remove", i, "from list cause it contains", matchedKeyWordsL[0])
		}

	}
	// log.Println(stringL)
	return stringL

}

const defaultMask = '*'

type (
	// TrieOption defines the method to customize a Trie.
	TrieOption func(trie *trieNode)

	// A Trie is a tree implementation that used to find elements rapidly.
	Trie interface {
		Filter(text string) (string, []string, bool)
		FindKeywords(text string) []string
	}

	trieNode struct {
		node
		mask rune
	}

	scope struct {
		start int
		stop  int
	}
)

// NewTrie returns a Trie.
func NewTrie(words []string, opts ...TrieOption) Trie {
	n := new(trieNode)

	for _, opt := range opts {
		opt(n)
	}
	if n.mask == 0 {
		n.mask = defaultMask
	}
	for _, word := range words {
		n.add(word)
	}

	n.build()

	return n
}

func (n *trieNode) Filter(text string) (sentence string, keywords []string, found bool) {
	chars := []rune(text)
	if len(chars) == 0 {
		return text, nil, false
	}

	scopes := n.find(chars)
	keywords = n.collectKeywords(chars, scopes)

	for _, match := range scopes {
		// we don't care about overlaps, not bringing a performance improvement
		n.replaceWithAsterisk(chars, match.start, match.stop)
	}

	return string(chars), keywords, len(keywords) > 0
}

func (n *trieNode) FindKeywords(text string) []string {
	chars := []rune(text)
	if len(chars) == 0 {
		return nil
	}

	scopes := n.find(chars)
	return n.collectKeywords(chars, scopes)
}

func (n *trieNode) collectKeywords(chars []rune, scopes []scope) []string {
	set := make(map[string]lang.PlaceholderType)
	for _, v := range scopes {
		set[string(chars[v.start:v.stop])] = lang.Placeholder
	}

	var i int
	keywords := make([]string, len(set))
	for k := range set {
		keywords[i] = k
		i++
	}

	return keywords
}

func (n *trieNode) replaceWithAsterisk(chars []rune, start, stop int) {
	for i := start; i < stop; i++ {
		chars[i] = n.mask
	}
}

// WithMask customizes a Trie with keywords masked as given mask char.
func WithMask(mask rune) TrieOption {
	return func(n *trieNode) {
		n.mask = mask
	}
}

type node struct {
	children map[rune]*node
	fail     *node
	depth    int
	end      bool
}

func (n *node) add(word string) {
	chars := []rune(word)
	if len(chars) == 0 {
		return
	}

	nd := n
	var depth int
	for i, char := range chars {
		if nd.children == nil {
			child := new(node)
			child.depth = i + 1
			nd.children = map[rune]*node{char: child}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			nd = child
			depth++
		} else {
			child := new(node)
			child.depth = i + 1
			nd.children[char] = child
			nd = child
		}
	}

	nd.end = true
}

func (n *node) build() {
	var nodes []*node
	for _, child := range n.children {
		child.fail = n
		nodes = append(nodes, child)
	}
	for len(nodes) > 0 {
		nd := nodes[0]
		nodes = nodes[1:]
		for key, child := range nd.children {
			nodes = append(nodes, child)
			cur := nd
			for cur != nil {
				if cur.fail == nil {
					child.fail = n
					break
				}
				if fail, ok := cur.fail.children[key]; ok {
					child.fail = fail
					break
				}
				cur = cur.fail
			}
		}
	}
}

func (n *node) find(chars []rune) []scope {
	var scopes []scope
	size := len(chars)
	cur := n

	for i := 0; i < size; i++ {
		child, ok := cur.children[chars[i]]
		if ok {
			cur = child
		} else {
			for cur != n {
				cur = cur.fail
				if child, ok = cur.children[chars[i]]; ok {
					cur = child
					break
				}
			}

			if child == nil {
				continue
			}
		}

		for child != n {
			if child.end {
				scopes = append(scopes, scope{
					start: i + 1 - child.depth,
					stop:  i + 1,
				})
			}
			child = child.fail
		}
	}

	return scopes
}

func RemoveSpecifiedElementInStringSlice(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}
