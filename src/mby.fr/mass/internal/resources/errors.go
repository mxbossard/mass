package resources

import (
	"fmt"
	//"strings"
)

type ResourceNotFound struct {
	Expression string
	Kinds *[]Kind
}
func (e ResourceNotFound) Error() string {
	message := fmt.Sprintf("Resource not found: %s for kinds: %v", e.Expression, *e.Kinds)
	return message
}

type UnknownKind struct {
	Kind string
}
func (e UnknownKind) Error() string {
        return fmt.Sprintf("Unknown kind: %s", e.Kind)
}

type InconsistentExpressionType struct {
	Expression string
	ExpectedTypes *[]Kind
}
func (e InconsistentExpressionType) Error() string {
	return fmt.Sprintf("Expression type: %s is not consistent with expected kinds: %v", e.Expression, *e.ExpectedTypes)
}

type InconsistentExpression struct {
	Expression string
	ExpectedTypes *[]Kind
}
func (e InconsistentExpression) Error() string {
        return fmt.Sprintf("Expression: %s is not consistent with expected kinds: %v", e.Expression, *e.ExpectedTypes)
}

