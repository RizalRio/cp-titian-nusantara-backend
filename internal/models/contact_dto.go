package models

// 🌟 DTOs (Request Payloads)
type CreateContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type CreateCollaborationRequest struct {
	OrganizationName  string `json:"organization_name" binding:"required"`
	ContactPerson     string `json:"contact_person" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	CollaborationType string `json:"collaboration_type" binding:"required"`
	Message           string `json:"message" binding:"required"`
	ProposalFileURL   string `json:"proposal_file_url"`
}