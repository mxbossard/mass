package resources

import (
	"fmt"
	"sort"
	"strings"
)

type Kind int

const (
	AllKind = Kind(iota)
	EnvKind
	ProjectKind
	ImageKind
)

func (k Kind) String() (s string) {
	switch k {
	case AllKind:
		s = "all"
	case EnvKind:
		s = "env"
	case ProjectKind:
		s = "project"
	case ImageKind:
		s = "image"
	}
	return
}

/*
func (k Kind) Match(o Kind) bool {
	return o == AllKind || k == o
}
*/

func (k Kind) MarshalYAML() (interface{}, error) {
	var s string
	switch k {
	case EnvKind:
		s = "env"
	case ProjectKind:
		s = "project"
	case ImageKind:
		s = "image"
	default:
		return "", fmt.Errorf("Unable to marshal kind: %s", k)
	}
	return s, nil
}

func (k *Kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	switch s {
	case "env":
		*k = EnvKind
	case "project":
		*k = ProjectKind
	case "image":
		*k = ImageKind
	default:
		return fmt.Errorf("Unable to unmarshal kind: %s", s)
	}
	return nil
}

type KindSet map[Kind]Kind

func (s KindSet) String() string {
	kinds := make([]string, 0, len(s))
	for _, k := range s.Kinds() {
		kinds = append(kinds, k.String())
	}
	sort.Strings(kinds)
	return strings.Join(kinds, ",")
}

func (s KindSet) Kinds() []Kind {
	kinds := make([]Kind, 0, len(s))
	for k := range s {
		kinds = append(kinds, k)
	}
	return kinds
}

func (s KindSet) Contains(kind Kind) bool {
	for _, k := range s.Kinds() {
		if k == AllKind || k == kind {
			return true
		}
	}
	return false
}

var allKindSet = &KindSet{AllKind: AllKind}

func NewKindSet(kinds ...Kind) *KindSet {
	set := KindSet{}
	for _, k := range kinds {
		if k == AllKind {
			return allKindSet
		}
		set[k] = k
	}

	if len(kinds) == 0 {
		return allKindSet
	}
	return &set
}

var kindAlias = map[Kind][]string{
	EnvKind:     []string{EnvKind.String()[0:1], EnvKind.String(), EnvKind.String() + "s"},
	ProjectKind: []string{ProjectKind.String()[0:1], ProjectKind.String(), ProjectKind.String() + "s"},
	ImageKind:   []string{ImageKind.String()[0:1], ImageKind.String(), ImageKind.String() + "s"},
	AllKind:     []string{AllKind.String()},
}

func KindExists(k Kind) bool {
	return k == EnvKind || k == ProjectKind || k == ImageKind || k == AllKind
}

func KindFromAlias(alias string) (Kind, bool) {
	for k, v := range kindAlias {
		for _, a := range v {
			if alias == a {
				return k, true
			}
		}
	}
	return AllKind, false
}

func IsKindIn(kind Kind, kinds []Kind) bool {
	for _, k := range kinds {
		if k == AllKind || k == kind {
			return true
		}
	}
	return false
}
