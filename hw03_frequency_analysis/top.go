package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var re = regexp.MustCompile(`([[:punct:]]*)([^[:punct:]]*)`)

func normalize(str string) string {
	// оставляем всё кроме первых и последних знаков пунктуации
	var sb strings.Builder
	matches := re.FindAllStringSubmatch(str, -1)
	for i := 0; i < len(matches); i++ {
		level := matches[i]
		if i > 0 && level[2] != "" {
			sb.WriteString(level[1])
		}
		sb.WriteString(level[2])
	}

	if sb.Len() == 0 && utf8.RuneCountInString(str) > 1 {
		// специальный случай: больше одного минуса подряд - это слово
		allHyphen := true
		for _, r := range str {
			if r != '-' {
				allHyphen = false
				break
			}
		}
		if allHyphen {
			sb.WriteString(str)
		}
	}
	return strings.ToLower(sb.String())
}

func Top10(src string) []string {
	// считаем количество появления слов
	words := make(map[string]int)
	for _, s := range strings.Fields(src) {
		s = normalize(s)
		if s != "" {
			words[s]++
		}
	}
	// список слов для сортировки
	sorted := make([]string, 0, len(words))
	for k := range words {
		sorted = append(sorted, k)
	}
	// сортируем по количеству и лексографически
	sort.Slice(sorted, func(i, j int) bool {
		c1 := words[sorted[i]]
		c2 := words[sorted[j]]
		return c1 > c2 || (c1 == c2 && sorted[i] < sorted[j])
	})
	return sorted[:min(len(sorted), 10)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
