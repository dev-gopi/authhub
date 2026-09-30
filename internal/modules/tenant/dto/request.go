package dto

type CreateTenantRequest struct {
	APILabel string `json:"api_label" validate:"required,min=2,max=63"`

	DisplayName string `json:"display_name" validate:"required,min=2,max=255"`

	Profile TenantProfileRequest `json:"profile" validate:"required"`

	PrimaryAdmin PrimaryAdminRequest `json:"primary_admin" validate:"required"`
}

type TenantProfileRequest struct {
	LegalName string `json:"legal_name" validate:"omitempty,max=255"`

	WebsiteURL string `json:"website_url" validate:"omitempty,url"`

	SupportEmail string `json:"support_email" validate:"omitempty,email,max=320"`

	SupportPhone string `json:"support_phone" validate:"omitempty,max=64"`

	SecurityContactEmail string `json:"security_contact_email" validate:"omitempty,email,max=320"`

	Locale string `json:"locale" validate:"omitempty,max=32"`

	Timezone string `json:"timezone" validate:"omitempty,max=64"`

	CountryCode string `json:"country_code" validate:"omitempty,len=2"`

	Region string `json:"region" validate:"omitempty,max=128"`

	PrivacyPolicyURL string `json:"privacy_policy_url" validate:"omitempty,url"`

	TermsURL string `json:"terms_url" validate:"omitempty,url"`
}

type PrimaryAdminRequest struct {
	Username string `json:"username" validate:"required,min=3,max=128"`

	Email string `json:"email" validate:"required,email,max=320"`

	DisplayName string `json:"display_name" validate:"required,min=2,max=255"`
}
