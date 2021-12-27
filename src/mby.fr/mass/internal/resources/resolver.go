package resources

import(
	"fmt"
	"strings"
	"path/filepath"

	"mby.fr/utils/file"
	"mby.fr/utils/errorz"
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
func ResolveExpression(expressions string, expectedKinds ...Kind) (resources []Resource, aggErr errorz.Aggregated) {
	splittedExpr, exprKinds, err := splitExpressions(expressions)
	if err != nil {
		aggErr.Add(err)
		return
	}

	// If no expectedKinds, assume expected AllKind
	if len(expectedKinds) == 0 {
		expectedKinds = append(expectedKinds, AllKind)
	}

	// Check for expression kind consistency
	var notExpectedKinds []Kind
	expectAllKinds := false
	exprAllKinds := false
	for _, exprKind := range exprKinds {
		kindFound := false
		if exprKind == AllKind {
			exprAllKinds = true
		}
		for _, kind := range expectedKinds {
			if kind == AllKind {
				expectAllKinds = true
				break
			}
			if exprKind == kind || exprKind == AllKind {
				kindFound = true
				break
			}
		}
		if !kindFound {
			// An expression kind is not in expected kinds
			notExpectedKinds = append(notExpectedKinds, exprKind)
		}
	}

	if !expectAllKinds && len(notExpectedKinds) > 0 {
		// Not expecting all kinds and found an expression kind not matching expectedKinds
		err = InconsistentExpression
		aggErr.Add(err)
		return
	}

	if len(exprKinds) == 0 || exprAllKinds {
		// If no kind supplied in expr, use expectedKinds as hint
		exprKinds = expectedKinds
	}

	// resolve all expressions versus all expr kinds
	for _, expr := range splittedExpr {
		res, errors := resolveExpresionForKinds(expr, exprKinds)
		//fmt.Printf("Resolved exprs: %s with kind: %s and found: %s\n", expr, exprKind, res)
		if errors.GotError() {
			aggErr.Concat(errors)
			continue
		}
		resources = append(resources, res)
	}

	return
}

func splitExpressions(expressions string) (strippedExpressions []string, kinds []Kind, err error) {
	splittedExpr := strings.Split(expressions, " ")

	// First expr may specify type if different from others expr
	// Attempt to resolve from work dir then from workspace dir
	firstExpr := splittedExpr[0]
	firstExprIndex := 1

	splittedFirstExpr := strings.Split(firstExpr, ",")

	for _, firstExprPart := range splittedFirstExpr {
		exprKind, ok := KindFromAlias(firstExprPart)
		if ok {
			kinds = append(kinds, exprKind)
		} else {
			// first keyword not a resource type
			firstExprIndex = 0
			if len(splittedFirstExpr) > 1 {
				// Unable to determine kind but multiple kind supplied
				return nil, nil, UnknownKind
			}
		}
	}

	strippedExpressions = strings.Split(expressions, " ")[firstExprIndex:]
	return
}

// How merge errors ?
// If unkown errors return them
// If a NotFoundError for a kind return it
// If InconsistentExpressionType for all kinds it wins
// Else return all errors
func resolveExpresionForKinds(expr string, kinds []Kind) (res Resource, aggErr errorz.Aggregated) {
	errorByKind := map[Kind]error{}
	for _, kind := range kinds {
		res, err := resolveResource(expr, kind)
		if err == nil {
			return res, aggErr
		} else {
			errorByKind[kind] = err
		}
	}

	var notFound error
	var inconsistentExpressionTypeCount = 0
	var unknownErrors errorz.Aggregated
	var knownErrors errorz.Aggregated
	for _, e := range errorByKind {
		switch e {
		case ResourceNotFound:
			notFound = e
			knownErrors.Add(e)
			break
		case InconsistentExpressionType:
			inconsistentExpressionTypeCount ++
			knownErrors.Add(e)
		default:
			unknownErrors.Add(e)
		}
	}

	if unknownErrors.GotError() {
		aggErr = unknownErrors
	} else if notFound != nil {
		aggErr.Add(notFound)
	} else if len(kinds) == inconsistentExpressionTypeCount {
		// FIXME; return new InconsistentExpressionType for all kinds
		aggErr.Add(InconsistentExpressionType)
	} else {
		aggErr = knownErrors
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

