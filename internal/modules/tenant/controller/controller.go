package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	rootmiddleware "github.com/dev-gopi/authhub/internal/modules/rootauth/middleware"
	"github.com/dev-gopi/authhub/internal/modules/tenant/dto"
	tenantentity "github.com/dev-gopi/authhub/internal/modules/tenant/entity"
	tenantservice "github.com/dev-gopi/authhub/internal/modules/tenant/service"
	tenantvalidator "github.com/dev-gopi/authhub/internal/modules/tenant/validator"
	"github.com/dev-gopi/authhub/internal/shared/response"
)

type Controller struct {
	service   tenantservice.Interface
	validator *tenantvalidator.Validator
}

func New(service tenantservice.Interface, validator *tenantvalidator.Validator) *Controller {
	return &Controller{service: service, validator: validator}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	auth, ok := rootmiddleware.FromRootAuthContext(r.Context())
	if !ok || auth == nil || !auth.IsRootAdmin {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var req dto.CreateTenantRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if err := c.validator.ValidateCreate(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid tenant request")
		return
	}

	result, err := c.service.Create(r.Context(), req, auth)
	if err != nil {
		switch {
		case errors.Is(err, tenantentity.ErrTenantAlreadyExists), errors.Is(err, tenantentity.ErrPrimaryAdminConflict):
			response.Error(w, http.StatusConflict, "CONFLICT", err.Error())
		case errors.Is(err, tenantentity.ErrUnauthorized):
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Root administrator access required")
		case errors.Is(err, tenantentity.ErrProvisioningFailed):
			response.Error(w, http.StatusServiceUnavailable, "TENANT_PROVISIONING_FAILED", "Tenant remains in provisioning state")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}
		return
	}
	response.Success(w, http.StatusCreated, result)
}
