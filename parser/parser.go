package parser

import (
	"github.com/aleury/goini"
)

func Parse(name, input string) *goini.File {
	file := &goini.File{
		Name:     name,
		Sections: make([]goini.Section, 0),
	}
	l := lex(name, input)
	section := goini.Section{}
	key := ""
	for {
		item := l.nextItem()
		if item.Type == itemEOF {
			file.Sections = append(file.Sections, section)
			break
		}
		switch item.Type {
		case itemSection:
			if len(section.KeyValuePairs) > 0 {
				file.Sections = append(file.Sections, section)
			}
			section.Name = item.Value
			section.KeyValuePairs = make([]goini.KeyValuePair, 0)
		case itemKey:
			key = item.Value
		case itemValue:
			section.KeyValuePairs = append(section.KeyValuePairs, goini.KeyValuePair{
				Key:   key,
				Value: item.Value,
			})
			key = ""
		}
	}
	return file
}
