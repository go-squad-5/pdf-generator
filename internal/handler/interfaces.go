package handler

type PDFService interface {
	GenerateQuizReport(sessionID string) ([]byte, error)
}

type EmailService interface {
	SendQuizReportByEmail(sessionID string) error
}
