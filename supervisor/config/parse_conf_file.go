package config

import (
	"bufio"
	"bytes"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gstr"
	"io"
	"strings"
)

func ParseIni(data []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	fieldMap := make(map[string]interface{})

	a := bytes.NewReader(data)
	r := bufio.NewReader(a)
	var section string
	var lastSection string
	var haveSection bool
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		lineStr := strings.TrimSpace(string(line))
		if len(lineStr) == 0 {
			continue
		}

		if lineStr[0] == ';' || lineStr[0] == '#' {
			continue
		}

		sectionBeginPos := strings.Index(lineStr, "[")
		sectionEndPos := strings.Index(lineStr, "]")

		if sectionBeginPos >= 0 && sectionEndPos >= 2 {
			section = lineStr[sectionBeginPos+1 : sectionEndPos]

			if lastSection == "" {
				lastSection = section
			} else if lastSection != section {
				lastSection = section
				fieldMap = make(map[string]interface{})
			}
			haveSection = true
		} else if !haveSection {
			continue
		}

		if strings.Contains(lineStr, "=") && haveSection {
			values := gstr.Split(lineStr, "=")
			fieldMap[strings.TrimSpace(values[0])] = strings.TrimSpace(strings.Join(values[1:], "="))
			res[section] = fieldMap
		}
	}

	if !haveSection {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "failed to parse INI file, section not found")
	}
	return res, nil

}

func ParseEnv(s string) *map[string]string {
	result := make(map[string]string)
	start := 0
	n := len(s)
	var i int
	for {
		// find the '='
		for i = start; i < n && s[i] != '='; {
			i++
		}
		key := s[start:i]
		start = i + 1
		if s[start] == '"' {
			for i = start + 1; i < n && s[i] != '"'; {
				i++
			}
			if i < n {
				result[strings.TrimSpace(key)] = strings.TrimSpace(s[start+1 : i])
			}
			if i+1 < n && s[i+1] == ',' {
				start = i + 2
			} else {
				break
			}
		} else {
			for i = start; i < n && s[i] != ','; {
				i++
			}
			if i < n {
				result[strings.TrimSpace(key)] = strings.TrimSpace(s[start:i])
				start = i + 1
			} else {
				result[strings.TrimSpace(key)] = strings.TrimSpace(s[start:])
				break
			}
		}
	}

	return &result
}
