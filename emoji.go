package main

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"strings"
)

//go:embed resources/emoji.json
var emojiFileData []byte
var emojis map[string]string

var re regexp.Regexp = *regexp.MustCompile(`:[^:\s]*:`)

func replaceShortcodesWithEmojis(text string) string {
	return re.ReplaceAllStringFunc(text, shortcodeToEmoji)
}

func convertSlackReaction(text string) string {
	var converted string
	emoji := strings.Split(text, "::")
	for _, e := range emoji {
		converted += shortcodeToEmoji(e)
	}
	return converted
}

func shortcodeToEmoji(code string) string {
	strippedCode := strings.TrimPrefix(code, ":")
	strippedCode = strings.TrimSuffix(strippedCode, ":")
	emoji, found := emojis[strippedCode]
	if found {
		return emoji
	} else {
		return code
	}
}

func emojiToShortcode(emoji string) string {
	var partCodes []string
	for _, r := range emoji {
		if r == '\ufe0f' { // ignore unicode variation selector
			continue
		}
		for code, e := range emojis {
			if string(r) == e {
				partCodes = append(partCodes, code)
				continue
			}
		}
	}
	return strings.Join(partCodes, "::")
}

func init() {
	json.Unmarshal(emojiFileData, &emojis)
}
