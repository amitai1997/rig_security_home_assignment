package policy

import (
	"context"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

// Violation represents a policy violation.
type Violation struct {
	Username   string
	Permission string
	Rule       string
}

// Engine evaluates collaborators against a policy.
type Engine interface {
	// Scan returns a violation if the collaborator breaks the policy.
	Scan(ctx context.Context, c Collaborator) (*Violation, error)
}

// Collaborator mirrors githook.Collaborator to avoid import cycles.
type Collaborator struct {
	Login      string
	Permission string
}

// celEngine implements Engine using Google's CEL.
type celEngine struct {
	program cel.Program
	rule    string
}

// NewEngine compiles the given CEL expression.
func NewEngine(policyExpr string) (Engine, error) {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("login", decls.String),
			decls.NewVar("permission", decls.String),
		),
	)
	if err != nil {
		return nil, err
	}
	ast, iss := env.Compile(policyExpr)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	return &celEngine{program: prg, rule: policyExpr}, nil
}

// Scan evaluates the collaborator. If the expression evaluates to true, a violation is returned.
func (e *celEngine) Scan(ctx context.Context, c Collaborator) (*Violation, error) {
	out, _, err := e.program.ContextEval(ctx, map[string]interface{}{
		"login":      c.Login,
		"permission": c.Permission,
	})
	if err != nil {
		return nil, err
	}
	passed, ok := out.Value().(bool)
	if !ok || !passed {
		return nil, nil
	}
	return &Violation{Username: c.Login, Permission: c.Permission, Rule: e.rule}, nil
}
