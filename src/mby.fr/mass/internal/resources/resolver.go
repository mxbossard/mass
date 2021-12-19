package resources

import(
	"fmt"
	"strings"

	"mby.fr/utils/file"
	"mby.fr/mass/internal/settings"
)

const EnvPrefix = "env/"
const ProjectPrefix = "project/"
const ImagePrefix = "image/"

var InconsistentExpressionPrefix error = fmt.Errorf("Expression prefix and kind are not consistent")
var InconsistentExpression error = fmt.Errorf("Expression and kind are not consistent")

// Resolve complex resource expression
// project1/image2 project1/image3
// project1 image2 image3
// project1
// .
func ResolveExpression(expressions string, resourceKind string) (resources []Resource, err error) {
	settingsService, err := settings.GetSettingsService()
	if err != nil {
                return
        }
	workspaceDir := settingsService.WorkspaceDir()

	workDir, err := file.WorkDirPath()
	if err != nil {
                return
        }

	relativeWorkPath := strings.Replace(workDir, workspaceDir, "", 1)

	splittedExpr := strings.Split(expressions, " ")

	// First expr may specify context if different from others expr
	// Attempt to resolve from work dir then from workspace dir
	firstExpr := splittedExpr[0]
	tryExpr := relativeWorkPath + "/" + firstExpr
	firstRes, err := resolveResource(tryExpr, resourceKind)
	if err != nil {
		err = nil
		tryExpr = firstExpr
		firstRes, err = resolveResource(tryExpr, resourceKind)
	}
	if err != nil {
		//err = fmt.Errorf("Unable to found resource for expr: %s !", firstExpr)
		err = ResourceNotFound
		return
	}

	for _, expr := range splittedExpr[1:] {
		res, ok := resolveResource(expr, resourceKind)
		_ = res
		_ = ok
	}

	_ = firstRes

	return
}

func checkConsistency(expr, resourceKind string) (err error) {
	if resourceKind == ImageKind {
		// Image kind cannot be resolved with absolute expr if not referencing a project.
		if strings.HasPrefix(expr, "/") || strings.HasPrefix(expr, ImagePrefix) {
		       slashCount := strings.Count(expr, "/")
			if slashCount != 2 {
				err = InconsistentExpression
				return
			}
		}
	}

	if strings.HasPrefix(expr, EnvPrefix) {
                if resourceKind != "" && resourceKind != EnvKind {
                        err = InconsistentExpressionPrefix
                        return
                }
	}

	if strings.HasPrefix(expr, ProjectPrefix) {
		if resourceKind != "" && resourceKind != ProjectKind {
			err = InconsistentExpressionPrefix
			return
		}
	}

	if strings.HasPrefix(expr, ImagePrefix) {
		if resourceKind != "" && resourceKind != ImageKind {
			err = InconsistentExpressionPrefix
			return
		}
	}

	return
}

// Resolve simple resource expression
// project1/image2
func resolveResource(expr, resourceKind string) (r Resource, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	workspaceDir := ss.WorkspaceDir()

	workDir, err := file.WorkDirPath()
	if err != nil {
		return
	}

	err = checkConsistency(expr, resourceKind)
	if err != nil {
		return
	}

	// dot or empty expr
	if expr == "" || expr == "." {
		return resolveDotResource(resourceKind)
	}

	// absolute expr
	if strings.HasPrefix(expr, "/") {
		//expr = expr[1:] // Strip first /
		// Resolve resource from workspace dir only
		return resolveResourceFrom(workspaceDir, expr, resourceKind)
	}

	fromDir := workDir

	// env dedicated expr
	if strings.HasPrefix(expr, EnvPrefix) {
		// Resolve env only
		expr = expr[len(EnvPrefix):] // Strip prefix
		resourceKind = EnvKind
	}

	// project dedicated expr
	if strings.HasPrefix(expr, ProjectPrefix) {
		// Resolve project only
		expr = expr[len(ProjectPrefix):] // Strip prefix
		resourceKind = ProjectKind
	}

	// image dedicated expr
	if strings.HasPrefix(expr, ImagePrefix) {
		// Resolve image only
		expr = expr[len(ImagePrefix):] // Strip prefix
		resourceKind = ImageKind
		fromDir = workspaceDir
	}

	if resourceKind == EnvKind || resourceKind == ProjectKind {
		fromDir = workspaceDir
	}

	// Resolve resource from work dir
	r, err = resolveResourceFrom(fromDir, expr, resourceKind)
	//if err == ResourceNotFound {
	//	// Continue resolving
	//	err = nil
	//} else if err != nil {
	//	return
	//}

	//// Resolve resource from workspace dir
	//r, err = resolveResourceFrom(workspaceDir, expr, resourceKind)
	return
}

func resolveResourceFrom(fromDir, expr, resourceKind string) (r Resource, err error) {
	resources, err := scanResourcesFrom(fromDir, resourceKind)
	if err != nil {
		return
	}
	for _, res := range resources {
		if res.MatchExpression(expr) && (resourceKind == "" || res.Kind() == resourceKind) {
			r = res
			return
		}
	}
	err = ResourceNotFound
	return
}

func scanResourcesFrom(fromDir, resourceKind string) (r []Resource, err error) {
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

// Resolve resource in working dir
func resolveDotResource(resourceKind string) (r Resource, err error) {
	workDir, err := file.WorkDirPath()
	if err != nil {
		return
	}
	r, err = Read(workDir)
	if err != nil {
		return
	}
	if resourceKind != "" && r.Kind() != resourceKind {
		r = nil
		err = ResourceNotFound
		return
	}
	return

}

