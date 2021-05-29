package gateway

import (
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由模式, 用于区分一个规则的 trie 路径, 如 /posts/:cates
	part     string  // 路由模式中的一部分, 以 `/` 划分, 如 :cates
	children []*node // 子节点列表, 例如 [docs, blogs, tutorials]
	isWild   bool    // 是否精确匹配, part 包含 `:` 或 `*` 时为 true
}

func newNode() *node {
	return &node{
		children: make([]*node, 0),
	}
}

// 将一个路由规则递归插入到 trie 树
// height = 0 when initialize
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 递归结束
		// 只有最后一层节点, 才会设置 pattern, 上层节点均为 pattern = ""
		// 用于判断路由规则是否被匹配成功
		n.pattern = pattern
		return
	}

	// 递归逻辑
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// 不存在则新建
		child = newNode()
		child.part = part
		child.isWild = (part[0] == ':' || part[0] == '*')

		n.children = append(n.children, child)
	}

	// 递归插入
	child.insert(pattern, parts, height+1)
}

// 返回一个匹配的路由规则
//
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 递归结束
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		// 递归搜索
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// 找到第一个与 part 模式匹配成功的子节点, 用于插入
// (路由注册)
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 找到所有与 part 模式匹配成功的子节点,用于查找
// (路由查询)
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
