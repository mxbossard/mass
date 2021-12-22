package resources

import(
	"fmt"
	"strings"
	"path/filepath"

	"mby.fr/utils/file"
	"mby.fr/mass/internal/settings"
)

const EnvPrefix = "env/"
const ProjectPrefix = "project/"
const ImagePrefix = "image/"

var InvalidArgument error = fmt.Errorf("Invalid argument error")
var InconsistentExpressionType error = fmt.Errorf("Expression type and kind are not consistent")
var InconsistentExpression error = fmt.Errorf("Expression and kind are not consistent")

// Resolve complex resource expression
// type
// type name1 name2 name3
// type1/name1 type2/name2 type3/name3
// all
// 
// projects
// projects p1 p2
// p p1 p2
// envs
// env e1 e2
// images
// image p1/i1 p2
// project1/image2 project1/image3
// project1 image2 image3
// project1
func ResolveExpression(expressions string, expectedKind Kind) (resources []Resource, err error) {
	//settingsService, err := settings.GetSettingsService()
	//if err != nil {
        //        return
        //}
	//workspaceDir := settingsService.WorkspaceDir()

	//workDir, err := file.WorkDirPath()
	//if err != nil {
        //        return
        //}

	//relativeWorkPath := strings.Replace(workDir, workspaceDir, "", 1)

	splittedExpr := strings.Split(expressions, " ")

	// First expr may specify type if different from others expr
	// Attempt to resolve from work dir then from workspace dir
	firstExpr := splittedExpr[0]
	firstExprIndex := 1
	exprKind, ok := KindFromAlias(firstExpr)
	if ok {
		if expectedKind != AllKind && exprKind != expectedKind {
			err = InconsistentExpression
			return
		}
	} else {
		// first keyword not a resource type we need to consume it as an expr
		exprKind = expectedKind
		firstExprIndex = 0
	}

	if exprKind == AllKind {
		// If no kind supplied in expr, use expectedKind as hint
		exprKind = expectedKind
	}

	for _, expr := range splittedExpr[firstExprIndex:] {
		res, err := resolveResource(expr, exprKind)
		if err != nil {
			return resources, err
		}
		resources = append(resources, res)
	}

	return
}

// Resolve simple resource expression with an expected kind. Fail if kind missmatch
// env1 env/env1 e/env1 envs/env1
// project1 project/project1 projects/project1
// project1/image2 image/project1/image2
//
func resolveResource(expr string, expectedKind Kind) (r Resource, err error) {
	if expr == "" || ! KindExists(expectedKind) {
		err = InvalidArgument
		return
	}

	kind, name := splitExpression(expr)
	if kind != AllKind && expectedKind != AllKind && kind != expectedKind {
		err = InconsistentExpressionType
		return
	}
	if expectedKind == AllKind {
		expectedKind = kind
	}
	if expectedKind == AllKind {
		err = InconsistentExpression
		return
	}

	// Resolve contextual resource
	r, err = resolveContextualResource(name, expectedKind)

	return
}

func splitExpression(expr string) (kind Kind, name string) {
	kind = AllKind
	name = expr
	res := strings.Split(name, "/")
	if len(res) > 1 {
		switch res[0] {
			case "p", "project", "projects":
				kind = ProjectKind
				name = strings.Join(res[1:], "/")
			case "e", "env", "envs":
				kind = EnvKind
				name = strings.Join(res[1:], "/")
			case "i", "image", "images":
				kind = ImageKind
				name = strings.Join(res[1:], "/")
		}
	}
	return
}

func splitImageName(name string) (project, image string) {
	res := strings.Split(name, "/")
	if len(res) == 1 {
		image = name
	} else {
		project = res[0]
		image = strings.Join(res[1:], "/")
	}
	return
}

// Return one resource matching name and kind in context (context depends on user working dir)
// project1 project2
// env1 env2
// project1/image1 image2
func resolveContextualResource(name string, kind Kind) (r Resource, err error) {
	if name == "" || ! KindExists(kind) {
		err = InvalidArgument
		return
	}

	ss, err := settings.GetSettingsService()
        if err != nil {
                return
        }
        workspaceDir := ss.WorkspaceDir()

	switch kind {
	case EnvKind, ProjectKind:
		return resolveResourceFrom(workspaceDir, name, kind)
	case ImageKind:
		// Special case: if image name does not have project prefix, can be resolve from work dir
		project, image := splitImageName(name)
		if project == "" {
			workDir, err := file.WorkDirPath()
			if err != nil {
				return r, err
			}
			image = filepath.Base(workDir) + "/" + image
			return resolveResourceFrom(workDir, image, kind)
		}
		return resolveResourceFrom(workspaceDir, name, kind)

	}

	return
}

// Return one resource matching name and kind in fromDir tree
func resolveResourceFrom(fromDir, name string, kind Kind) (r Resource, err error) {
	if name == "" || fromDir == "" || ! KindExists(kind) {
		err = InvalidArgument
		return
	}

	resources, err := scanResourcesFrom(fromDir, kind)
	if err != nil {
		return
	}
	for _, res := range resources {
		if res.Name() == name && res.Kind() == kind {
			r = res
			return
		}
	}
	err = ResourceNotFound
	return
}

func isResourceMatchingExpr(r Resource, expr string) bool {
	return r.Name() == expr
}

func scanResourcesFrom(fromDir string, resourceKind Kind) (r []Resource, err error) {
	if resourceKind == "" || resourceKind == EnvKind {
		envs, err := ScanEnvs(fromDir)
		if err != nil {
			return r, err
		}
		for _, e := range envs {
			r = append(r, e)
		}
	}

	if resourceKind == "" || resourceKind == ProjectKind {
		projects, err := ScanProjects(fromDir)
		if err != nil {
			return r, err
		}
		for _, p := range projects {
			r = append(r, p)
		}
	}

	if resourceKind == "" || resourceKind == ImageKind {
		images, err := ScanImages(fromDir)
		if err != nil {
			return r, err
		}
		for _, i := range images {
			r = append(r, i)
		}
	}

	return
}

