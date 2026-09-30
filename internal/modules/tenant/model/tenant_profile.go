package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TenantProfile struct {
	sharedmodel.BaseModel

	TenantID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	LegalName *string `gorm:"type:text"`

	LogoAssetID *uuid.UUID `gorm:"type:uuid"`

	WebsiteURL *string `gorm:"type:text"`

	SupportEmail *string `gorm:"type:text"`

	SupportPhone *string `gorm:"type:text"`

	SecurityContactEmail *string `gorm:"type:text"`

	Locale *string `gorm:"type:text"`

	Timezone *string `gorm:"type:text"`

	CountryCode *string `gorm:"type:text"`

	Region *string `gorm:"type:text"`

	PrivacyPolicyURL *string `gorm:"type:text"`

	TermsURL *string `gorm:"type:text"`

	Metadata datatypes.JSON `gorm:"type:jsonb;not null"`
}

func (TenantProfile) TableName() string {
	return "tenant_profiles"
}
