package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strings"
)

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao ler:", err)
	}

	return ""
}

func main() {
	for {
		print("Digite uma palavra: ")
		text := readLine()
		println()
		run(text)

		print("\nQuer digitar outra? (s/n): ")
		choice := readLine()
		if strings.ToLower(strings.TrimSpace(choice)) != "s" {
			break
		}
		println()
	}
}

func run(text string) {
	println("Texto original: ", text)
	freq := countFreqChar(text)
	root := newHuffmanTree(freq)
	codes := generateCodes(root)
	encoded := encodeText(text, codes)
	decoded := decodeBinary(encoded, root)
	println("Texto codificado: ", encoded)
	println("Texto decodificado: ", decoded)
	plotTree(root)
}

type Node struct {
	Char  rune
	Freq  int
	Left  *Node
	Right *Node
}

func countFreqChar(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, char := range s {
		freq[char]++
	}
	return freq
}

func newHuffmanTree(freq map[rune]int) *Node {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for char, frequency := range freq {
		heap.Push(&pq, &Node{
			Char: char,
			Freq: frequency,
		})
	}

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)

		parent := &Node{
			Char:  0,
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&pq, parent)
	}

	return heap.Pop(&pq).(*Node)
}

func generateCodes(root *Node) map[rune]string {
	codes := make(map[rune]string)
	generateCodesHelper(root, "", codes)
	return codes
}

func generateCodesHelper(node *Node, code string, codes map[rune]string) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		codes[node.Char] = code
		return
	}

	generateCodesHelper(node.Left, code+"0", codes)
	generateCodesHelper(node.Right, code+"1", codes)
}

func encodeText(text string, codes map[rune]string) string {
	var encoded strings.Builder
	for _, char := range text {
		encoded.WriteString(codes[char])
	}
	return encoded.String()
}

func decodeBinary(binary string, root *Node) string {
	var decoded strings.Builder
	current := root

	for _, bit := range binary {
		if bit == '0' {
			current = current.Left
		} else {
			current = current.Right
		}

		if current.Left == nil && current.Right == nil {
			decoded.WriteRune(current.Char)
			current = root
		}
	}

	return decoded.String()
}

func plotTree(root *Node) {
	fmt.Println("\n--- Árvore de Huffman")
	plotTreeHelper(root, "", true)
}

func plotTreeHelper(node *Node, prefix string, isLeft bool) {
	if node == nil {
		return
	}

	fmt.Print(prefix)
	if isLeft {
		fmt.Print("├── ")
	} else {
		fmt.Print("└── ")
	}

	if node.Char != 0 {
		fmt.Printf("'%c' (freq: %d)\n", node.Char, node.Freq)
	} else {
		fmt.Printf("* (freq: %d)\n", node.Freq)
	}

	newPrefix := prefix
	if isLeft {
		newPrefix += "│   "
	} else {
		newPrefix += "    "
	}

	if node.Left != nil {
		plotTreeHelper(node.Left, newPrefix, true)
	}
	if node.Right != nil {
		plotTreeHelper(node.Right, newPrefix, false)
	}
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Node))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
