package fidelity

type FidelityIssue struct {
	Field      string `json:"field"`
	Severity   string `json:"severity"`
	Problem    string `json:"problem"`
	Suggestion string `json:"suggestion"`
}

type FidelityResult struct {
	Passed bool            `json:"passed"`
	Issues []FidelityIssue `json:"issues"`
}
