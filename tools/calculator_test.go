package tools

import (
    "context"
    "math"
    "strconv"
    "testing"
)

func TestCalculator_Call(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"Basic addition", "2 + 2", "4", false},
        {"Exponentiation (fixes #870)", "30 ** 0.23", "2.944", false}, // ~2.9436567476686945
        {"Multiplication", "5 * 3", "15", false},
        {"Division", "10 / 2", "5", false}, // Adjusted to expect "5" not "5.0"
        {"Invalid exponentiation", "30 **", "error from evaluator: invalid exponentiation expression: 30 **", false},
        {"Empty input", "", "", true},
        {"Invalid syntax", "2 + + 2", "error from evaluator: invalid syntax in expression: 2 + + 2", false},
    }

    calc := Calculator{}

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calc.Call(context.Background(), tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Errorf("Call(%q) error = nil, want error", tt.input)
                }
                return
            }
            if err != nil {
                t.Errorf("Call(%q) error = %v, want nil", tt.input, err)
                return
            }

            if tt.input == "30 ** 0.23" {
                gotFloat, _ := strconv.ParseFloat(got, 64)
                wantFloat, _ := strconv.ParseFloat(tt.want, 64)
                if !approxEqual(gotFloat, wantFloat, 0.001) {
                    t.Errorf("Call(%q) = %v, want %v (approx)", tt.input, got, tt.want)
                }
            } else if got != tt.want {
                t.Errorf("Call(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}

func approxEqual(a, b, tol float64) bool {
    return math.Abs(a-b) <= tol
}