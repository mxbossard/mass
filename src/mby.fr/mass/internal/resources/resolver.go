package resources

import(
	"fmt"
	"strings"

	"mby.fr/utils/file"
	"mby.fr/mass/internal/settings"
)

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
		err = fmt.Errorf("Unable to found resource for expr: %s !", firstExpr)
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

// Resolve simple resource expression
// project1/image2
func resolveResource(expr, resourceKind string) (r Resource, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	wksDir := ss.WorkspaceDir()
	projects, err := ScanProjects(wksDir)
	if err != nil {
		return
	}
	for _, p := range projects {
		if p.Name() == expr && p.Kind() == resourceKind {
			r = p
			return
		}
	}
	err = fmt.Errorf("Resource not found")
	return
}
