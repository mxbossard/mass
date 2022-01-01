package resources

import(
	"fmt"
	//"os"
	//"sync"
	//"path/filepath"
        //"gopkg.in/yaml.v2"
)

//type Kind string
//func (k Kind) String() string {
//	return string(k)
//}
//
//const AllKind Kind = "all"
//const EnvKind Kind = "env"
//const ProjectKind Kind = "project"
//const ImageKind Kind = "image"

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
type KindSet map[Kind]struct{}
var exists = struct{}{}
var allKindSet = map[Kind]struct{}{AllKind: exists}

func Kinds(kinds ...Kind) KindSet {
	set := KindSet{}
	for _, k := range kinds {
		if k == AllKind {
			return allKindSet
		}
		set[k] = exists
	}

	//if len(kinds) == 0 {
	//	return allKindSet
	//}
	return set
}

var kindAlias = map[Kind][]string {
	EnvKind: []string{EnvKind.String()[0:1], EnvKind.String(), EnvKind.String() + "s"},
	ProjectKind: []string{ProjectKind.String()[0:1], ProjectKind.String(), ProjectKind.String() + "s"},
	ImageKind: []string{ImageKind.String()[0:1], ImageKind.String(), ImageKind.String() + "s"},
	AllKind: []string{AllKind.String()},
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

