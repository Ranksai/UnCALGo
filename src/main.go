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
	AtmarkData                    = "&(&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y)@(&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y)"
	InclusionSampleData           = "(&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y)"
	InclusionSampleDataNotInclude = "&x1 := {a:&x2}, &x2 := {b:&x2,c:&x3}, &x3 := &y"
	InclusionSampleSingle         = "&x1 := {a:&x2}"
	InclusionSampleDouble         = "&x2 := {b:&x2,c:&x3}"
	InclusionSampleNone           = "&x3 := &y"
)

type Node struct {
	Name     string
	Id       int
	Parent   []Node
	Children map[string]string
}

type Nodes = []Node

// 与えられた d の 1 つのノードについて処理する
func calcInclusion(inclusionCurlyBracket string) (Node, error) {
	fmt.Println("input is")
	fmt.Println(inclusionCurlyBracket)

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

	moveMap := make(map[string]string, 0)

	// 移動するノード分確認
	for i := range splitResult {
		// 移動先を見る
		if strings.Contains(splitResult[i], ":") {
			moveNode := strings.Split(splitResult[i], ":")
			moveMap[moveNode[0]] = moveNode[1]
		} else {
			// TODO: ε 遷移が複数あるときに対応できない 今だと最後の ε 遷移のみ
			// &x3 = &y のようなとき ε 遷移
			moveMap["ε"] = splitResult[i]
		}
	}
	node := Node{
		Name:     nodeName,
		Id:       0,
		Parent:   nil,
		Children: moveMap,
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

func main() {
	result, err := multiCalcInclusion(InclusionSampleDataNotInclude)
	if err != nil {
		fmt.Printf("Err func main: %+w", err)
	}
	fmt.Println("print node")
	fmt.Println(result)
}
