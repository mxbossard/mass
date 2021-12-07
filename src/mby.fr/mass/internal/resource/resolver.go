package resource

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
func ResolveExpression(expressions string) (resources []Resource, err error) {
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
	firstRes, ok := resolveResource(tryExpr)
	if !ok {
		tryExpr = firstExpr
		firstRes, ok = resolveResource(tryExpr)
	}
	if !ok {
		err = fmt.Errorf("Unable to found resource for expr: %s !", firstExpr)
		return
	}

	for _, expr := range splittedExpr[1:] {
		res, ok := resolveResource(expr)
		_ = res
		_ = ok
	}

	_ = firstRes

	return
}

// Resolve simple resource expression
// project1/image2
func resolveResource(expr string) (r Resource, ok bool) {
	return
}
