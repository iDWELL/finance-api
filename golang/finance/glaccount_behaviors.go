package finance

// GLAccountCostCenterHandling mirrors iX-Haus Sachkonto.KSTSTBEHANDLUNG.
// An empty value (the zero value) is treated as "not defined" on the server
// side and persisted as SQL NULL.
type GLAccountCostCenterHandling string

const (
	GLAccountCostCenterHandlingNotDefined GLAccountCostCenterHandling = ""
	GLAccountCostCenterHandlingMustNot    GLAccountCostCenterHandling = "MUST_NOT"
	GLAccountCostCenterHandlingMust       GLAccountCostCenterHandling = "MUST"
	GLAccountCostCenterHandlingMay        GLAccountCostCenterHandling = "MAY"
)

// GLAccountProjectHandling mirrors iX-Haus Sachkonto.PRJBEHANDLUNG.
// An empty value (the zero value) is treated as "not defined" on the server
// side and persisted as SQL NULL.
type GLAccountProjectHandling string

const (
	GLAccountProjectHandlingNotDefined GLAccountProjectHandling = ""
	GLAccountProjectHandlingMustNot    GLAccountProjectHandling = "MUST_NOT"
	GLAccountProjectHandlingMust       GLAccountProjectHandling = "MUST"
	GLAccountProjectHandlingMay        GLAccountProjectHandling = "MAY"
)

// GLAccountOrderHandling mirrors iX-Haus Sachkonto.AuftragBehandlung.
// An empty value (the zero value) is treated as "not defined" on the server
// side and persisted as SQL NULL.
type GLAccountOrderHandling string

const (
	GLAccountOrderHandlingNotDefined GLAccountOrderHandling = ""
	GLAccountOrderHandlingMustNot    GLAccountOrderHandling = "MUST_NOT"
	GLAccountOrderHandlingMust       GLAccountOrderHandling = "MUST"
	GLAccountOrderHandlingMay        GLAccountOrderHandling = "MAY"
)
