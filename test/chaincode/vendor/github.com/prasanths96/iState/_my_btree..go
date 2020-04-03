//

package istate

import (
	"fmt"
	"reflect"
	"strings"
)

type BST struct {
	Parent *BSTNode
	Smallest
}

type BSTNode struct {
	Left  *BSTNode
	Right *BSTNode
	Data  []ComparatorInterface // To store multiple data of same value in single node
}

type ComparatorInterface interface {
	Compare(interface{}) int // Returns -1 if value is lessthan input, 0 for equal, 1 for greaterthan
	GetElem() interface{}
}

func NewBST(elem ComparatorInterface) (bst *BST) {
	bstNode := &BSTNode{
		Left:  nil,
		Right: nil,
		Data:  []ComparatorInterface{elem},
	}
	bst = &BST{
		Parent: bstNode,
	}
	return
}

//
func (tree *BST) Add(elem ComparatorInterface) {
	comparison := tree.Parent.Data[0].Compare(elem.GetElem())
	switch comparison {
	case -1:
		// Add to left
		switch leftNode := tree.Parent.Left; leftNode == nil {
		case true:
			tree.Parent.Left = &BSTNode{
				Left:  nil,
				Right: nil,
				Data:  []ComparatorInterface{elem},
			}
		default:
			bst := &BST{
				Parent: leftNode,
			}
			bst.Add(elem)
		}
	case 0:
		// Add here
		tree.Parent.Data = append(tree.Parent.Data, elem)
	case 1:
		// Add to right
		switch rightNode := tree.Parent.Right; rightNode == nil {
		case true:
			tree.Parent.Right = &BSTNode{
				Left:  nil,
				Right: nil,
				Data:  []ComparatorInterface{elem},
			}
		default:
			bst := &BST{
				Parent: rightNode,
			}
			bst.Add(elem)
		}
	}
}

//
func (tree *BST) GetSmallestElem() (elem []interface{}) {
	switch leftNode = tree.Parent.Left; leftNode == nil {
	case true:
		elem = make([]interface{}, len(tree.Parent.Data), len(tree.Parent.Data))
		for i := 0; i < len(tree.Parent.Data); i++ {
			elem[i] = tree.Parent.Data[i].GetElem()
		}
	default:
		bst := &BST{
			Parent: leftNode,
		}
		elem = bst.GetSmallest()
	}
	return
}

//
func (tree *BST) GetLargestElem() (elem []interface{}) {
	switch rightNode := tree.Parent.Right; rightNode == nil {
	case true:
		elem = make([]interface{}, len(tree.Parent.Data), len(tree.Parent.Data))
		for i := 0; i < len(tree.Parent.Data); i++ {
			elem[i] = tree.Parent.Data[i].GetElem()
		}
	default:
		bst := &BST{
			Parent: rightNode,
		}
		elem = bst.GetLargest()
	}
	return
}

//
func (tree *BST) GetLeftMostNode() (node *BSTNode) {
	switch node = tree.Parent.Left; node == nil {
	case true:
		return
	default:
		bst := &BST{
			Parent: node,
		}
		node = bst.GetLeftMostNode()
	}
	return
}

//
func (tree *BST) GetRightMostNode() (node *BSTNode) {
	switch node = tree.Parent.Right; node == nil {
	case true:
		return
	default:
		bst := &BST{
			Parent: node,
		}
		node = bst.GetRightMostNode()
	}
	return
}

//
func (tree *BST) GetOrderedList() (elem []interface{}) {
	node := tree.GetLeftMostNode()

}

//
func (node *BSTNode) GetLeftNode() (node *BSTNode) {
	node = node.Left
	return
}

//
func (node *BSTNode) GetRightNode() (node *BSTNode) {
	node = node.Right
	return
}
