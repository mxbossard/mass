package resources

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
)

type Kind int

const (
	AllKind = Kind(iota)
	EnvKind
	ProjectKind
	ImageKind
	kindLimit
)

func TypeFromKind(kind Kind) (t reflect.Type) {
	switch kind {
	case EnvKind:
		t = reflect.TypeOf((*Env)(nil)).Elem()
	case ProjectKind:
		t = reflect.TypeOf((*Project)(nil)).Elem()
	case ImageKind:
		t = reflect.TypeOf((*Image)(nil)).Elem()
	default:
		log.Fatalf("Resource type not found for kind %s !", kind)
	}
	return
}

func KindFromResource(res Resourcer) (k Kind) {
	t := reflect.TypeOf(res)
	return KindFromType(t)
	//return KindFromType(any)(res).(type))
	// Iterate over all kinds
	/*
		for k = Kind(0); i < kindLimit; i++ {
			if (interface{})(res).(type) == k.ResourceType() {
				return
			}
		}
		log.Fatalf("Resource kind not found for type %T !", res)
	*/
}

func KindFromType(t reflect.Type) (k Kind) {
	// Iterate over all kinds
	for k = Kind(1); k < kindLimit; k++ {
		if t == k.ResourceType() {
			return
		}
	}
	log.Fatalf("Resource kind not found for type %s !", t)
	return
}

func (k Kind) ResourceType() (t reflect.Type) {
	return TypeFromKind(k)
}

func (k Kind) String() (s string) {
	switch k {
	case AllKind:
		s = "all"
	default:
		s = strings.ToLower(k.ResourceType().Name())
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
	case AllKind:
		return "", fmt.Errorf("AllKind not marshallable !")
	default:
		s = k.String()
	}
	return s, nil
}

func (k *Kind) UnmarshalYAML0(unmarshal func(interface{}) error) error {
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

func (k *Kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	// Iterate over all kinds
	for kind := Kind(1); kind < kindLimit; kind++ {
		if s == kind.String() {
			*k = kind
			return nil
		}
	}

	return fmt.Errorf("Unable to unmarshal kind: %s", s)
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
