package Gim

import "strings"

//我们通过树结构查询，如果中间某一层的节点都不满足条件，那么就说明没有匹配到的路由，查询结束
type node struct {
	pattern  string  //待匹配路由
	part     string  //路由中一部分
	children []*node //子节点
	isWild   bool    //是否精确匹配 part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node)matchChild(part string)*node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}


//找到所有匹配的子节点,用于查找
func (n *node)matchChildren(part string)[]*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
func (n *node)insert(pattern string,parts []string,height int) {
	if len(parts) == height {
		//映射完成,记录pattern
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {//找不到子节点实列,新建子节点
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//路由匹配
func (n *node)search(parts []string,height int)*node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}


