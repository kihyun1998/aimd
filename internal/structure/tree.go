package structure

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Tree는 디렉토리 구조를 표현하는 인터페이스
type Tree interface {
	BuildTree(files []string) error
	ToMarkdown() string
}

// Node는 파일 시스템의 노드를 표현
type Node struct {
	Name     string
	IsDir    bool
	Children map[string]*Node
}

// directoryTree는 Tree 인터페이스 구현체
type directoryTree struct {
	root     *Node
	rootPath string
}

// NewDirectoryTree는 새로운 directoryTree 인스턴스를 생성
func NewDirectoryTree(rootPath string) Tree {
	return &directoryTree{
		root: &Node{
			Name:     filepath.Base(rootPath),
			IsDir:    true,
			Children: make(map[string]*Node),
		},
		rootPath: rootPath,
	}
}

func (dt *directoryTree) BuildTree(files []string) error {
	for _, file := range files {
		relPath, err := filepath.Rel(dt.rootPath, file)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		parts := strings.Split(filepath.ToSlash(relPath), "/")
		current := dt.root

		for i, part := range parts {
			isLast := i == len(parts)-1
			if _, exists := current.Children[part]; !exists {
				current.Children[part] = &Node{
					Name:     part,
					IsDir:    !isLast,
					Children: make(map[string]*Node),
				}
			}
			current = current.Children[part]
		}
	}
	return nil
}

// ToMarkdown은 트리구조를 마크다운으로 변환하는 함수
func (dt *directoryTree) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString("## Project Structure\n\n")
	sb.WriteString("```\n")
	sb.WriteString(dt.root.Name + "/\n")
	dt.writeNode(&sb, dt.root, "", true)
	sb.WriteString("```\n\n")
	return sb.String()
}

func (dt *directoryTree) writeNode(sb *strings.Builder, node *Node, prefix string, isLast bool) {
	children := make([]*Node, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		return children[i].Name < children[j].Name
	})

	for i, child := range children {
		isLastChild := i == len(children)-1
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		sb.WriteString(prefix)
		if isLastChild {
			sb.WriteString("└── ")
		} else {
			sb.WriteString("├── ")
		}

		if child.IsDir {
			sb.WriteString(child.Name + "/\n")
		} else {
			sb.WriteString(child.Name + "\n")
		}

		if child.IsDir {
			dt.writeNode(sb, child, newPrefix, isLastChild)
		}
	}
}
