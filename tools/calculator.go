package tools

import (
    "context"
    "fmt"
    "strings"

    "github.com/tmc/langchaingo/callbacks"
    "go.starlark.net/lib/math"
    "go.starlark.net/starlark"
)

type Calculator struct {
    CallbacksHandler callbacks.Handler
}

var _ Tool = Calculator{}

func (c Calculator) Description() string {
    return `Useful for getting the result of a math expression.
    The input to this tool should be a valid mathematical expression that could be executed by a Starlark evaluator.
    Supports basic arithmetic (+, -, *, /), exponentiation (e.g., "pow(30, 0.23)" or "30 ** 0.23"), and parentheses.`
}

func (c Calculator) Name() string {
    return "calculator"
}

func (c Calculator) Call(ctx context.Context, input string) (string, error) {
    if c.CallbacksHandler != nil {
        c.CallbacksHandler.HandleToolStart(ctx, input)
    }

    expr := strings.TrimSpace(input)
    if expr == "" {
        return "", fmt.Errorf("empty expression provided")
    }

    if strings.Contains(expr, "**") {
        parts := strings.Split(expr, "**")
        if len(parts) != 2 {
            return fmt.Sprintf("invalid exponentiation expression: %s", expr), nil
        }
        base := strings.TrimSpace(parts[0])
        exponent := strings.TrimSpace(parts[1])
        expr = fmt.Sprintf("pow(%s, %s)", base, exponent)
    }

    v, err := starlark.Eval(&starlark.Thread{Name: "main"}, "input", expr, math.Module.Members)
    if err != nil {
        if c.CallbacksHandler != nil {
            c.CallbacksHandler.HandleToolError(ctx, err)
        }
        return fmt.Sprintf("error from evaluator: %s", err.Error()), nil
    }

    result := v.String()
    if c.CallbacksHandler != nil {
        c.CallbacksHandler.HandleToolEnd(ctx, result)
    }

    return result, nil
}