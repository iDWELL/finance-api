package finance

// PartnerParagraph13bHandling mirrors iX-Haus N13BOPTION.
// An empty value (the zero value) is treated as "not defined" on the server
// side and persisted as SQL NULL.
type PartnerParagraph13bHandling string

const (
	PartnerParagraph13bHandlingNotDefined PartnerParagraph13bHandling = ""
	PartnerParagraph13bHandlingMustNot    PartnerParagraph13bHandling = "MUST_NOT"
	PartnerParagraph13bHandlingMay        PartnerParagraph13bHandling = "MAY"
	PartnerParagraph13bHandlingMust       PartnerParagraph13bHandling = "MUST"
)
