package app

import (
	"math/rand"
	"strconv"
	"strings"
)

type Generator struct {
	words []string
	punct []string
	rnd   *rand.Rand
}

func NewGenerator(source rand.Source) *Generator {
	return &Generator{
		words: defaultWords,
		punct: defaultPunct,
		rnd:   rand.New(source),
	}
}

func (g *Generator) Build(count int, opts Options) []rune {
	if count <= 0 {
		return nil
	}
	words := make([]string, 0, count)
	for i := 0; i < count; i++ {
		words = append(words, g.nextWord(opts))
	}
	return []rune(strings.Join(words, " "))
}

func (g *Generator) Extend(target []rune, count int, opts Options) []rune {
	if count <= 0 {
		return target
	}
	words := make([]string, 0, count)
	for i := 0; i < count; i++ {
		words = append(words, g.nextWord(opts))
	}
	if len(target) == 0 {
		return []rune(strings.Join(words, " "))
	}
	return append(append(target, ' '), []rune(strings.Join(words, " "))...)
}

func (g *Generator) nextWord(opts Options) string {
	if opts.Numbers && g.rnd.Intn(10) == 0 {
		return strconv.Itoa(g.rnd.Intn(9999) + 1)
	}
	word := g.words[g.rnd.Intn(len(g.words))]
	if opts.Punctuation && g.rnd.Intn(5) == 0 {
		word += g.punct[g.rnd.Intn(len(g.punct))]
	}
	return word
}

var defaultPunct = []string{
	".", ",", ";", ":", "!", "?", "@", "#", "&", "%", "^", "*", "(", ")",
}

var defaultWords = []string{
	"about", "above", "across", "after", "again", "air", "all", "almost", "also", "always",
	"among", "an", "and", "another", "answer", "any", "around", "as", "ask", "at",
	"away", "back", "base", "be", "because", "been", "before", "begin", "below", "between",
	"both", "bring", "build", "but", "by", "call", "can", "change", "check", "close",
	"come", "common", "consider", "could", "course", "create", "day", "decide", "did", "different",
	"do", "does", "done", "down", "each", "early", "easy", "end", "enough", "even",
	"every", "example", "eye", "face", "fact", "family", "far", "fast", "feel", "few",
	"find", "first", "for", "form", "found", "from", "full", "get", "give", "go",
	"good", "great", "group", "grow", "had", "hand", "hard", "has", "have", "he",
	"head", "hear", "help", "her", "here", "high", "him", "his", "hold", "home",
	"how", "if", "important", "in", "into", "is", "it", "its", "just", "keep",
	"kind", "know", "large", "last", "late", "lead", "learn", "leave", "left", "less",
	"let", "life", "like", "line", "little", "long", "look", "made", "make", "man",
	"many", "may", "mean", "might", "more", "most", "move", "much", "must", "my",
	"name", "near", "need", "never", "new", "next", "no", "not", "now", "number",
	"of", "off", "often", "old", "on", "once", "one", "only", "open", "or",
	"order", "other", "our", "out", "over", "own", "part", "people", "place", "plan",
	"play", "point", "possible", "problem", "program", "public", "put", "question", "quick", "read",
	"real", "really", "right", "run", "said", "same", "saw", "say", "school", "see",
	"seem", "set", "she", "should", "show", "side", "small", "so", "some", "something",
	"sound", "start", "state", "still", "story", "such", "system", "take", "tell", "than",
	"that", "the", "their", "them", "then", "there", "these", "they", "thing", "think",
	"this", "those", "time", "to", "too", "try", "turn", "two", "under", "up",
	"use", "very", "want", "was", "way", "we", "well", "went", "were", "what",
	"when", "where", "which", "while", "who", "why", "will", "with", "word", "work",
	"world", "would", "write", "year", "you", "your",
}
