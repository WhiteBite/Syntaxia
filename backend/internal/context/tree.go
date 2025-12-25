package context

import (
	"path/filepath"
	"sort"
	"strings"
)

// treeNode represents a node in the file tree
type treeNode struct {
	name     string
	children map[string]*treeNode
	isFile   bool
}

// buildSimpleTree builds a tree representation of file paths
func (s *Service) buildSimpleTree(paths []string) string {
	root := &treeNode{name: ".", children: make(map[string]*treeNode)}

	// Build tree structure
	for _, path := range paths {
		parts := strings.Split(filepath.ToSlash(path), "/")
		current := root

		for i, part := range parts {
			if part == "" || part == "." {
				continue
			}

			if _, exists := current.children[part]; !exists {
				current.children[part] = &treeNode{
					name:     part,
					children: make(map[string]*treeNode),
					isFile:   i == len(parts)-1,
				}
			}
			current = current.children[part]
		}
	}

	// Generate tree string
	var builder strings.Builder
	s.walkTree(root, "", true, &builder)
	return builder.String()
}

// walkTree recursively walks the tree and builds the string representation
func (s *Service) walkTree(node *treeNode, prefix string, isLast bool, builder *strings.Builder) {
	if node.name != "." {
		if isLast {
			builder.WriteString(prefix + "└─ " + node.name + "\n")
			prefix += "   "
		} else {
			builder.WriteString(prefix + "├─ " + node.name + "\n")
			prefix += "│  "
		}
	}

	// Sort children for deterministic output
	childNames := make([]string, 0, len(node.children))
	for name := range node.children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)

	for i, name := range childNames {
		child := node.children[name]
		isLastChild := i == len(childNames)-1
		s.walkTree(child, prefix, isLastChild, builder)
	}
}
