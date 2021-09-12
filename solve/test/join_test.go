package solve_test

import "testing"

func TestSolveJoins_basic(t *testing.T) {
	const level = `
	0j1  0  0j1
	`
	solutions := []string{
		"000",
		"111",
		"222",
	}

	testSolveJoinsAbstract(t, level, solutions, 3)
}
