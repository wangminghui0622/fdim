package sensitiveword

import (
	"strings"
	"sync"
	"unicode"
)

// Filter 敏感词过滤器（基于DFA算法）
type Filter struct {
	root  *trieNode
	mutex sync.RWMutex
}

// trieNode 字典树节点
type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

// NewFilter 创建敏感词过滤器
func NewFilter() *Filter {
	return &Filter{
		root: &trieNode{
			children: make(map[rune]*trieNode),
		},
	}
}

// AddWord 添加敏感词
func (f *Filter) AddWord(word string) {
	if word == "" {
		return
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()

	word = strings.ToLower(word)
	node := f.root
	for _, char := range word {
		if node.children[char] == nil {
			node.children[char] = &trieNode{
				children: make(map[rune]*trieNode),
			}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// AddWords 批量添加敏感
func (f *Filter) AddWords(words []string) {
	for _, word := range words {
		f.AddWord(word)
	}
}

// Contains 检查文本是否包含敏感词
func (f *Filter) Contains(text string) bool {
	if text == "" {
		return false
	}
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	text = strings.ToLower(text)
	runes := []rune(text)
	length := len(runes)

	for i := 0; i < length; i++ {
		node := f.root
		for j := i; j < length; j++ {
			char := runes[j]
			// 跳过特殊字符
			if isSpecialChar(char) {
				continue
			}
			if node.children[char] == nil {
				break
			}
			node = node.children[char]
			if node.isEnd {
				return true
			}
		}
	}
	return false
}

// FindAll 查找所有敏感词
func (f *Filter) FindAll(text string) []string {
	if text == "" {
		return nil
	}
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	text = strings.ToLower(text)
	runes := []rune(text)
	length := len(runes)
	var result []string

	for i := 0; i < length; i++ {
		node := f.root
		var word []rune
		for j := i; j < length; j++ {
			char := runes[j]
			// 跳过特殊字符
			if isSpecialChar(char) {
				continue
			}
			if node.children[char] == nil {
				break
			}
			node = node.children[char]
			word = append(word, char)
			if node.isEnd {
				result = append(result, string(word))
			}
		}
	}
	return result
}

// Replace 替换敏感词
func (f *Filter) Replace(text string, replacement rune) string {
	if text == "" {
		return text
	}
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	runes := []rune(strings.ToLower(text))
	originalRunes := []rune(text)
	length := len(runes)
	result := make([]rune, length)
	copy(result, originalRunes)

	for i := 0; i < length; i++ {
		node := f.root
		matchStart := i
		matchEnd := -1
		for j := i; j < length; j++ {
			char := runes[j]
			// 跳过特殊字符
			if isSpecialChar(char) {
				continue
			}
			if node.children[char] == nil {
				break
			}
			node = node.children[char]
			if node.isEnd {
				matchEnd = j
			}
		}
		if matchEnd >= 0 {
			for k := matchStart; k <= matchEnd; k++ {
				if !isSpecialChar(runes[k]) {
					result[k] = replacement
				}
			}
			i = matchEnd
		}
	}
	return string(result)
}

// isSpecialChar 判断是否为特殊字符（用于跳过干扰字符）
func isSpecialChar(char rune) bool {
	return unicode.IsSpace(char) || unicode.IsPunct(char) || unicode.IsSymbol(char)
}

// DefaultFilter 默认敏感词过滤器
var DefaultFilter *Filter

func init() {
	DefaultFilter = NewFilter()
	// 添加默认敏感词列表
	defaultWords := []string{
		// 这里可以添加默认的敏感词
		// 实际使用时应该从配置文件或数据库加载
	}
	DefaultFilter.AddWords(defaultWords)
}

// LoadWordsFromSlice 从切片加载敏感词
func LoadWordsFromSlice(words []string) {
	DefaultFilter.AddWords(words)
}

// ContainsSensitive 检查是否包含敏感词
func ContainsSensitive(text string) bool {
	return DefaultFilter.Contains(text)
}

// FindSensitiveWords 查找敏感词
func FindSensitiveWords(text string) []string {
	return DefaultFilter.FindAll(text)
}

// ReplaceSensitive 替换敏感词
func ReplaceSensitive(text string) string {
	return DefaultFilter.Replace(text, '*')
}
