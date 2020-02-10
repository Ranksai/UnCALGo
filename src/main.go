package main

import (
	"fmt"
	"strings"
)

// rule
// x0....N: node
// a b c..: 枝
// x1 := {}, x2 := {} 必ず{}, のあとには空白を入れる

const (
	SampleData                    = "&x1@cycle(&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y)"
	AtmarkSingle                  = "&x1@&x1"
	AtmarkData                    = "(&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y)@(&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y)"
	InclusionSampleData           = "(&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y)"
	InclusionSampleDataNotInclude = "&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y"
	InclusionSampleSingle         = "&x1 := {a:&x2}"
	InclusionSampleDouble         = "&x2 := {b:&x2,c:&x3}"
	InclusionSampleNone           = "&x3 := &y"
)

type Tree struct {
	Name       string
	StartNodes []string
	EndNodes   []string
	TreeNodes  Nodes
}

type Node struct {
	Name     string
	Children []ChildNode
	Eps      []string
}

type ChildNode struct {
	Name  string
	Label string
}

type ChildrenNodes = []ChildNode

type Nodes = []Node

// 与えられた d の 1 つのノードについて処理する
func calcInclusion(inclusionCurlyBracket string) (Node, error) {
	//fmt.Println("input is")
	//fmt.Println(inclusionCurlyBracket)

	// := で split
	// 左辺が対象ノード，右辺が左辺の対象ノードからどの枝を使ってどのノードに行くか書かれている ex: a:&x2 a の枝を使用して &x2 のノードに移動
	splitResult := strings.Split(strings.ReplaceAll(inclusionCurlyBracket, " ", ""), ":=")
	nodeName := splitResult[0]

	// 中括弧を削除
	splitResult[1] = strings.TrimLeft(splitResult[1], "{")
	splitResult[1] = strings.TrimRight(splitResult[1], "}")

	// , で split
	// 複数移動するノードがあったときにカンマでわかる
	splitResult = strings.Split(splitResult[1], ",")

	childrenNodes := make(ChildrenNodes, 0)
	epsList := make([]string, 0)
	// 移動するノード分確認
	for i := range splitResult {
		// 移動先を見る
		if strings.Contains(splitResult[i], ":") {
			childNode := strings.Split(splitResult[i], ":")
			childrenNodes = append(childrenNodes, ChildNode{
				Name:  childNode[1],
				Label: childNode[0],
			})
		} else {
			// &x3 = &y のようなとき ε 遷移
			epsList = append(epsList, splitResult[i])
		}
	}
	node := Node{
		Name:     nodeName,
		Children: childrenNodes,
		Eps:      epsList,
	}
	return node, nil
}

func multiCalcInclusion(inclusionCurlyBracket string) (Nodes, error) {
	nodes := make(Nodes, 0)
	// カンマ区切り
	splitResult := strings.Split(inclusionCurlyBracket, ", ")

	for i := range splitResult {
		node, _ := calcInclusion(splitResult[i])
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// cycleInput is not include "cycle"
func calcCycle(cycleInput string) (Nodes, error) {
	nodes, err := multiCalcInclusion(cycleInput)
	if err != nil {
		return nil, err
	}
	for i := range nodes {
		nodes[i].Eps = append(nodes[i].Eps, nodes[i].Name)
	}
	return nodes, nil
}

func calcAtmark(left, right string, tree Tree) (Tree, error) {
	if strings.Contains(left, ",") {
	} else {
		tree.StartNodes = append(tree.StartNodes, left)
	}
	return tree, nil
}

func checkEndNodes(tree Tree) Tree {
	childNodeName := make([]string, 0)
	endNodes := make([]string, 0)
	for _, node := range tree.TreeNodes {
		for _, childNode := range node.Children {
			childNodeName = append(childNodeName, childNode.Name)
		}
		for _, eps := range node.Eps {
			childNodeName = append(childNodeName, eps)
		}
	}
	for i := range childNodeName {
		flag := true
		for _, node := range tree.TreeNodes {
			if childNodeName[i] == node.Name {
				flag = false
				break
			}
		}
		if flag {
			endNodes = append(endNodes, childNodeName[i])
		}
	}
	tree.EndNodes = endNodes
	return tree
}

func calc(input string, tree Tree) (Tree, error) {
	if strings.Contains(input, "@") {
		splits := strings.Split(input, "@")
		// cycle 側の計算
		cycleInput := strings.Trim(splits[1], "cycle")
		cycleInput = strings.TrimLeft(cycleInput, "(")
		cycleInput = strings.TrimRight(cycleInput, ")")
		nodes, err := calcCycle(cycleInput)
		if err != nil {
			return tree, err
		}
		// @ の計算
		tree, _ = calcAtmark(splits[0], "", tree)
		tree.TreeNodes = nodes
		// endNode の確認
		tree = checkEndNodes(tree)
	} else {
		cycleInput := strings.TrimLeft(input, "(")
		cycleInput = strings.TrimRight(cycleInput, ")")
		nodes, err := calcCycle(cycleInput)
		if err != nil {
			return tree, err
		}
		tree.TreeNodes = nodes
		// endNode の確認
		tree = checkEndNodes(tree)
	}
	return tree, nil
}

func main() {
	// init tree
	tree := Tree{
		Name:       "UnCAL Tree",
		StartNodes: nil,
		EndNodes:   nil,
		TreeNodes:  nil,
	}
	// input tree data
	input := InclusionSampleDataNotInclude
	tree, _ = calc(input, tree)
	fmt.Printf("%+v", tree)
}
