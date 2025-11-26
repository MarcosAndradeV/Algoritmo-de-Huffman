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

func run(texto string) {
	println("Texto original: ", texto)
	freq := contarFrequencia(texto)
	raiz := Huffman(freq)
	codigos := gerarCodigos(raiz)
	codificado := codificar(texto, codigos)
	decodificado := decodificar(codificado, raiz)
	println("Texto codificado: ", codificado)
	println("Texto decodificado: ", decodificado)
	plot(raiz)
}

type No struct {
	Char  rune
	Freq  int
	Esquerdo  *No
	Direito *No
}

func contarFrequencia(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, char := range s {
		freq[char]++
	}
	return freq
}

func Huffman(freqs map[rune]int) *No {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for char, freq := range freqs {
		heap.Push(&pq, &No{
			Char: char,
			Freq: freq,
		})
	}

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*No)
		right := heap.Pop(&pq).(*No)

		parent := &No{
			Char:  0,
			Freq:  left.Freq + right.Freq,
			Esquerdo:  left,
			Direito: right,
		}
		heap.Push(&pq, parent)
	}

	return heap.Pop(&pq).(*No)
}

func gerarCodigos(root *No) map[rune]string {
	codes := make(map[rune]string)
	gerarCodigosRec(root, "", codes)
	return codes
}

func gerarCodigosRec(node *No, code string, codes map[rune]string) {
	if node == nil {
		return
	}

	if node.Esquerdo == nil && node.Direito == nil {
		codes[node.Char] = code
		return
	}

	gerarCodigosRec(node.Esquerdo, code+"0", codes)
	gerarCodigosRec(node.Direito, code+"1", codes)
}

func codificar(text string, codes map[rune]string) string {
	var encoded strings.Builder
	for _, char := range text {
		encoded.WriteString(codes[char])
	}
	return encoded.String()
}

func decodificar(binario string, raiz *No) string {
	var decoded strings.Builder
	current := raiz

	for _, bit := range binario {
		if bit == '0' {
			current = current.Esquerdo
		} else {
			current = current.Direito
		}

		if current.Esquerdo == nil && current.Direito == nil {
			decoded.WriteRune(current.Char)
			current = raiz
		}
	}

	return decoded.String()
}


func plot(raiz *No) {
	fmt.Println("\n--- Árvore de Huffman")
	plotArvore(raiz, "", true)
}

func plotArvore(no *No, prefixo string, esquerdo bool) {
	if no == nil {
		return
	}

	fmt.Print(prefixo)
	if esquerdo {
		fmt.Print("├── ")
	} else {
		fmt.Print("└── ")
	}

	if no.Char != 0 {
		fmt.Printf("'%c' (freq: %d)\n", no.Char, no.Freq)
	} else {
		fmt.Printf("* (freq: %d)\n", no.Freq)
	}

	novoPrefixo := prefixo
	if esquerdo {
		novoPrefixo += "│   "
	} else {
		novoPrefixo += "    "
	}

	if no.Esquerdo != nil {
		plotArvore(no.Esquerdo, novoPrefixo, true)
	}
	if no.Direito != nil {
		plotArvore(no.Direito, novoPrefixo, false)
	}
}

type PriorityQueue []*No

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*No))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
