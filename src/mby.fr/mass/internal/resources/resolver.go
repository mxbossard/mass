package resources

import (
	"fmt"
	"path/filepath"
	"strings"

	"mby.fr/mass/internal/settings"
	"mby.fr/utils/errorz"
	"mby.fr/utils/filez"
)

const EnvPrefix = "env/"
const ProjectPrefix = "project/"
const ImagePrefix = "image/"

var InvalidArgument error = fmt.Errorf("Invalid argument error")

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
func ResolveExpression(expressions string, expectedKinds ...Kind) (resources []Resourcer, aggErr errorz.Aggregated) {
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
		err = InconsistentExpression{expressions, NewKindSet(expectedKinds...)}
		aggErr.Add(err)
		return
	}

	if len(exprKinds) == 0 || exprAllKinds {
		// If no kind supplied in expr, use expectedKinds as hint
		exprKinds = expectedKinds
	}

	// resolve all expressions versus all expr kinds
	for _, expr := range splittedExpr {
		res, errors := resolveExpresionForKinds(expr, *NewKindSet(exprKinds...))
		//fmt.Printf("Resolved exprs: %s with kind: %s and found: %s\n", expr, exprKind, res)
		if errors.GotError() {
			aggErr.Concat(errors)
			continue
		}
		resources = append(resources, res)
	}

	return
}

func ResolveUniqResourceExpression[T Resourcer](expression string) (resource T, err error) {
	kind := KindFromResource(resource)
	resources, aggErr := ResolveExpression(expression, kind)
	if aggErr.GotError() {
		err = aggErr.Return()
		return
	}

	if len(resources) > 1 {
		err = fmt.Errorf("Expression %s is resolved in more than 1 reource !", expression)
	} else if len(resources) == 0 {
		err = ResourceNotFound{Expression: expression, Kinds: NewKindSet(kind)}
	} else {
		var ok bool
		resource, ok = resources[0].(T)
		if !ok {
			err = fmt.Errorf("Unable to wrap resolved %v resource into supplied kind: %v", resource, kind)
		}
	}

	return
}

// Split exprssion into a list of kinds and a list of expressions
func splitExpressions(expressions string) (strippedExpressions []string, kinds []Kind, err error) {
	if expressions == "" || expressions == "." {
		return []string{expressions}, []Kind{AllKind}, nil
	}

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
				return nil, nil, UnknownKind{firstExprPart}
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
func resolveExpresionForKinds(expr string, kinds KindSet) (res Resourcer, aggErr errorz.Aggregated) {
	kindInExpr, name := splitExpression(expr)
	if kindInExpr != AllKind {
		if !kinds.Contains(kindInExpr) {
			// Kind in expr is not in kinds => Inconsistent
			err := InconsistentExpressionType{expr, &kinds}
			aggErr.Add(err)
			return
		} else {
			// Restrict kinds to kind in expr
			kinds = *NewKindSet(kindInExpr)
		}
	}

	errorByKind := map[Kind]error{}
	for _, kind := range kinds {
		res, err := resolveResource(name, kind)
		if err == nil {
			return res, aggErr
		} else {
			errorByKind[kind] = err
		}
	}

	var notFoundKinds []Kind
	var inconsistentExpressionTypeCount = 0
	var unknownErrors errorz.Aggregated
	var knownErrors errorz.Aggregated
	for _, kind := range kinds {
		e := errorByKind[kind]
		switch e := e.(type) {
		case ResourceNotFound:
			notFoundKinds = append(notFoundKinds, kind)
		case InconsistentExpressionType:
			inconsistentExpressionTypeCount++
			knownErrors.Add(e)
		default:
			unknownErrors.Add(e)
		}
	}

	if unknownErrors.GotError() {
		aggErr = unknownErrors
	} else if len(notFoundKinds) > 0 {
		notFound := ResourceNotFound{Expression: name, Kinds: NewKindSet(notFoundKinds...)}
		aggErr.Add(notFound)
	} else if len(kinds) == inconsistentExpressionTypeCount {
		aggErr.Add(InconsistentExpressionType{expr, &kinds})
	} else {
		aggErr = knownErrors
	}
	return
}

// Resolve simple resource expression with an expected kind. Fail if kind missmatch
// env1 env/env1 e/env1 envs/env1
// project1 project/project1 projects/project1
// project1/image2 image/project1/image2
func resolveResource(expr string, expectedKind Kind) (r Resourcer, err error) {
	if !KindExists(expectedKind) {
		err = InvalidArgument
		return
	}

	kind, name := splitExpression(expr)
	if kind != AllKind && expectedKind != AllKind && kind != expectedKind {
		err = InconsistentExpressionType{expr, NewKindSet(expectedKind)}
		return
	}
	if expectedKind == AllKind {
		expectedKind = kind
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

	// Filter name
	name = strings.TrimSuffix(name, "/")

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
func resolveContextualResource(name string, kind Kind) (r Resourcer, err error) {
	if !KindExists(kind) {
		err = InvalidArgument
		return
	}

	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	workspaceDir := ss.WorkspaceDir()

	workDir, err2 := filez.WorkDirPath()
	if err != nil {
		return r, err2
	}

	switch kind {
	case AllKind:
		// For AllKind always use work dir
		r, err = resolveResourceFrom(workDir, name, kind)
	case EnvKind, ProjectKind, ImageKind:
		r, err = resolveResourceFrom(workDir, name, kind)
		if err != nil && !IsResourceNotFound(err) {
			return
		} else if err != nil {
			// Try to resolve from workspace dir
			r, err = resolveResourceFrom(workspaceDir, name, kind)
		}
		if err != nil {
			return
		}
	}

	// Rewrite ResourceNotFound content
	if _, ok := err.(ResourceNotFound); err != nil && ok {
		err = ResourceNotFound{name, NewKindSet(kind), err}
	}

	return
}

// Return one resource matching name and kind in fromDir tree
// Special case for images : match the name in general but match image name in project dir
func resolveResourceFrom(fromDir, name string, kind Kind) (r Resourcer, err error) {
	if fromDir == "" || !KindExists(kind) {
		err = InvalidArgument
		return
	}

	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	// 1- Attempt to resolve resource name as path
	nameAsPath := filepath.Join(fromDir, name)

	if nameAsPath == ss.WorkspaceDir() {
		// If in root workspace dir return what ?
		err = InvalidArgument
		return
	}

	res, err := getDirResource(nameAsPath, kind)
	if err == nil && (kind == AllKind || kind == res.Kind()) {
		return res, err
	}
	// If error swallow it.

	// 2- Attempt to resolve resource name in fromDir
	res, err = getDirResource(fromDir, kind)
	if err == nil && res.Match(name, kind) {
		return res, err
	}

	// 3- Attempt to resolve resource name first part
	splittedName := strings.Split(name, PathSeparator)
	if len(splittedName) > 1 {
		firstDir := fromDir
		firstName := splittedName[0]
		remainingName := strings.Join(splittedName[1:], PathSeparator)
		if len(firstName) == 0 {
			// First name char is a spearator => resolve from root dir
			firstDir = ss.WorkspaceDir()
			firstName = splittedName[1]
			if len(splittedName) > 2 {
				remainingName = strings.Join(splittedName[2:], PathSeparator)
			} else {
				remainingName = ""
			}
		}
		firstRes, error := resolveResourceFrom(firstDir, firstName, AllKind)
		if error != nil && !IsResourceNotFound(error) {
			// Swallow ResourceNotFound
			return nil, error
		}
		if firstRes != nil && remainingName != "" {
			foundRes, error := resolveResourceFrom(firstRes.Dir(), remainingName, kind)
			if error != nil && !IsResourceNotFound(error) {
				// Swallow ResourceNotFound
				return nil, error
			} else if error == nil {
				return foundRes, nil
			}
		}
	}

	// 4- Scan to find resource
	var resources []Resourcer
	if kind == AllKind {
		resources, err = scanResourcesFrom(fromDir, AllKind, 1)
	} else {
		resources, err = scanResourcesFrom(fromDir, kind, -1)
	}
	if err != nil {
		return
	}

	// Filter found resources
	// Keep only resources matching specified kind
	// For Image, keep by Name() in general, plus keep by ImageName() if in context of a Project
	for _, res := range resources {
		if kind == AllKind || kind == res.Kind() {
			switch v := res.(type) {
			case Image:
				if v.FullName() == name {
					// Image general case
					return res, nil
				}

				fromDirRes, err := ReadResourcer(fromDir)
				if err == nil && fromDirRes.Kind() == ProjectKind && v.ImageName() == name {
					// In Project context
					return res, nil
				}

				if err != nil && !IsResourceNotFound(err) {
					// Swallow ResourceNotFound
					return r, err
				}

			default:
				if v.FullName() == name {
					return res, nil
				}
			}
		}
	}
	err = ResourceNotFound{Expression: name, Kinds: NewKindSet(kind)}
	return
}

func isResourceMatchingExpr(r Resourcer, expr string) bool {
	return r.FullName() == expr
}

// Return resource with kind in dir if it exists
func getDirResource(fromDir string, resourceKind Kind) (res Resourcer, err error) {
	r, err := ReadResourcer(fromDir)
	if err != nil {
		return
	}
	if resourceKind == AllKind || resourceKind == r.Kind() {
		res = r
	} else {
		err = ResourceNotFound{Expression: fromDir, Kinds: NewKindSet(resourceKind)}
	}

	return
}

// Return resource with kind in current dir if it exists
func getCurrentDirResource(resourceKind Kind) (res Resourcer, err error) {
	workDir, err := filez.WorkDirPath()
	if err != nil {
		return
	}
	return getDirResource(workDir, resourceKind)
}
