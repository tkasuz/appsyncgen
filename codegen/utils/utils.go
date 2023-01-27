package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(path string) (*string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	str := string(b)
	return &str, nil
}

func ToYaml(data any) string {
	y, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalln(err)
	}
	return string(y)
}

func ToJson(data any) string {
	y, err := json.MarshalIndent(&data, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	return string(y)
}

func CreateFile(dir string, filename string) (file *os.File, err error) {
	err = os.MkdirAll(filepath.Join(dir), os.ModePerm)
	if err != nil {
		return nil, err
	}
	file, err = os.Create(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}
	return file, err
}

type Pair struct {
	First  any
	Second any
}

func (p Pair) IsSame(pp Pair) bool {
	if (p.First == pp.First && p.Second == pp.Second) || (p.First == pp.Second && p.Second == pp.First) {
		return true
	} else {
		return false
	}
}

type PairList []*Pair

func (l PairList) HasSamePair(p Pair) bool {
	for _, it := range l {
		if it.IsSame(p) {
			return true
		}
	}
	return false
}

func (l PairList) ForFirst(v any) any {
	for _, it := range l {
		if it.First == v {
			return it.Second
		}
	}
	return nil
}

func (l PairList) ForPair(s any) any {
	for _, it := range l {
		if it.First == s {
			return it.Second
		} else if it.Second == s {
			return it.First
		}
	}
	return nil
}

func FormatGraphqlSchema(file *os.File) error {
	lines := make([]string, 0)
	f, err := os.Open(file.Name())
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) != 0 {
			lines = append(lines, s.Text())
		}
	}
	newFile, err := os.Create(file.Name())
	if err != nil {
		return err
	}
	for _, line := range lines {
		newFile.WriteString(fmt.Sprintln(line))
		if line == "}" {
			newFile.WriteString("\n")
		}
	}
	return nil
}
