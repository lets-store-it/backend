// Code generated by ogen, DO NOT EDIT.

package api

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *ErrorStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// DeleteOrgOK is response for DeleteOrg operation.
type DeleteOrgOK struct{}

// DeleteUnitOK is response for DeleteUnit operation.
type DeleteUnitOK struct{}

// Represents error object.
// Ref: #/components/schemas/Error
type Error struct {
	Message string        `json:"message"`
	Details *ErrorDetails `json:"details"`
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// GetDetails returns the value of Details.
func (s *Error) GetDetails() *ErrorDetails {
	return s.Details
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

// SetDetails sets the value of Details.
func (s *Error) SetDetails(val *ErrorDetails) {
	s.Details = val
}

type ErrorDetails struct{}

// ErrorStatusCode wraps Error with StatusCode.
type ErrorStatusCode struct {
	StatusCode int
	Response   Error
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorStatusCode) GetResponse() Error {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorStatusCode) SetResponse(val Error) {
	s.Response = val
}

// Merged schema.
type GetUnitByIdOK struct {
	CreatedAt time.Time `json:"created_at"`
	// ID of the employee who created the entity.
	CreatedBy NilUUID     `json:"created_by"`
	UpdatedAt NilDateTime `json:"updated_at"`
	// ID of the employee who last updated the entity.
	UpdatedBy NilUUID        `json:"updated_by"`
	Spaces    []StorageSpace `json:"spaces"`
}

// GetCreatedAt returns the value of CreatedAt.
func (s *GetUnitByIdOK) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCreatedBy returns the value of CreatedBy.
func (s *GetUnitByIdOK) GetCreatedBy() NilUUID {
	return s.CreatedBy
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *GetUnitByIdOK) GetUpdatedAt() NilDateTime {
	return s.UpdatedAt
}

// GetUpdatedBy returns the value of UpdatedBy.
func (s *GetUnitByIdOK) GetUpdatedBy() NilUUID {
	return s.UpdatedBy
}

// GetSpaces returns the value of Spaces.
func (s *GetUnitByIdOK) GetSpaces() []StorageSpace {
	return s.Spaces
}

// SetCreatedAt sets the value of CreatedAt.
func (s *GetUnitByIdOK) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetCreatedBy sets the value of CreatedBy.
func (s *GetUnitByIdOK) SetCreatedBy(val NilUUID) {
	s.CreatedBy = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *GetUnitByIdOK) SetUpdatedAt(val NilDateTime) {
	s.UpdatedAt = val
}

// SetUpdatedBy sets the value of UpdatedBy.
func (s *GetUnitByIdOK) SetUpdatedBy(val NilUUID) {
	s.UpdatedBy = val
}

// SetSpaces sets the value of Spaces.
func (s *GetUnitByIdOK) SetSpaces(val []StorageSpace) {
	s.Spaces = val
}

// Merged schema.
type GetUnitsOK struct {
	// Merged property.
	Items    []GetUnitsOKItemsItem `json:"items"`
	Metadata PaginationMetadata    `json:"metadata"`
}

// GetItems returns the value of Items.
func (s *GetUnitsOK) GetItems() []GetUnitsOKItemsItem {
	return s.Items
}

// GetMetadata returns the value of Metadata.
func (s *GetUnitsOK) GetMetadata() PaginationMetadata {
	return s.Metadata
}

// SetItems sets the value of Items.
func (s *GetUnitsOK) SetItems(val []GetUnitsOKItemsItem) {
	s.Items = val
}

// SetMetadata sets the value of Metadata.
func (s *GetUnitsOK) SetMetadata(val PaginationMetadata) {
	s.Metadata = val
}

// Merged schema.
type GetUnitsOKItemsItem struct {
	CreatedAt time.Time `json:"created_at"`
	// ID of the employee who created the entity.
	CreatedBy NilUUID     `json:"created_by"`
	UpdatedAt NilDateTime `json:"updated_at"`
	// ID of the employee who last updated the entity.
	UpdatedBy NilUUID `json:"updated_by"`
}

// GetCreatedAt returns the value of CreatedAt.
func (s *GetUnitsOKItemsItem) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCreatedBy returns the value of CreatedBy.
func (s *GetUnitsOKItemsItem) GetCreatedBy() NilUUID {
	return s.CreatedBy
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *GetUnitsOKItemsItem) GetUpdatedAt() NilDateTime {
	return s.UpdatedAt
}

// GetUpdatedBy returns the value of UpdatedBy.
func (s *GetUnitsOKItemsItem) GetUpdatedBy() NilUUID {
	return s.UpdatedBy
}

// SetCreatedAt sets the value of CreatedAt.
func (s *GetUnitsOKItemsItem) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetCreatedBy sets the value of CreatedBy.
func (s *GetUnitsOKItemsItem) SetCreatedBy(val NilUUID) {
	s.CreatedBy = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *GetUnitsOKItemsItem) SetUpdatedAt(val NilDateTime) {
	s.UpdatedAt = val
}

// SetUpdatedBy sets the value of UpdatedBy.
func (s *GetUnitsOKItemsItem) SetUpdatedBy(val NilUUID) {
	s.UpdatedBy = val
}

// NewNilDateTime returns new NilDateTime with value set to v.
func NewNilDateTime(v time.Time) NilDateTime {
	return NilDateTime{
		Value: v,
	}
}

// NilDateTime is nullable time.Time.
type NilDateTime struct {
	Value time.Time
	Null  bool
}

// SetTo sets value to v.
func (o *NilDateTime) SetTo(v time.Time) {
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o NilDateTime) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *NilDateTime) SetToNull() {
	o.Null = true
	var v time.Time
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o NilDateTime) Get() (v time.Time, ok bool) {
	if o.Null {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o NilDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewNilUUID returns new NilUUID with value set to v.
func NewNilUUID(v uuid.UUID) NilUUID {
	return NilUUID{
		Value: v,
	}
}

// NilUUID is nullable uuid.UUID.
type NilUUID struct {
	Value uuid.UUID
	Null  bool
}

// SetTo sets value to v.
func (o *NilUUID) SetTo(v uuid.UUID) {
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o NilUUID) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *NilUUID) SetToNull() {
	o.Null = true
	var v uuid.UUID
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o NilUUID) Get() (v uuid.UUID, ok bool) {
	if o.Null {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o NilUUID) Or(d uuid.UUID) uuid.UUID {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptInt32 returns new OptInt32 with value set to v.
func NewOptInt32(v int32) OptInt32 {
	return OptInt32{
		Value: v,
		Set:   true,
	}
}

// OptInt32 is optional int32.
type OptInt32 struct {
	Value int32
	Set   bool
}

// IsSet returns true if OptInt32 was set.
func (o OptInt32) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt32) Reset() {
	var v int32
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt32) SetTo(v int32) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt32) Get() (v int32, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt32) Or(d int32) int32 {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptUUID returns new OptUUID with value set to v.
func NewOptUUID(v uuid.UUID) OptUUID {
	return OptUUID{
		Value: v,
		Set:   true,
	}
}

// OptUUID is optional uuid.UUID.
type OptUUID struct {
	Value uuid.UUID
	Set   bool
}

// IsSet returns true if OptUUID was set.
func (o OptUUID) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptUUID) Reset() {
	var v uuid.UUID
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptUUID) SetTo(v uuid.UUID) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptUUID) Get() (v uuid.UUID, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptUUID) Or(d uuid.UUID) uuid.UUID {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/Organization
type Organization struct {
	ID        OptUUID `json:"id"`
	Name      string  `json:"name"`
	Subdomain string  `json:"subdomain"`
}

// GetID returns the value of ID.
func (s *Organization) GetID() OptUUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *Organization) GetName() string {
	return s.Name
}

// GetSubdomain returns the value of Subdomain.
func (s *Organization) GetSubdomain() string {
	return s.Subdomain
}

// SetID sets the value of ID.
func (s *Organization) SetID(val OptUUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *Organization) SetName(val string) {
	s.Name = val
}

// SetSubdomain sets the value of Subdomain.
func (s *Organization) SetSubdomain(val string) {
	s.Subdomain = val
}

// Merged schema.
// Ref: #/components/schemas/OrganizationsPagedResponse
type OrganizationsPagedResponse struct {
	// Merged property.
	Items    []OrganizationsPagedResponseItemsItem `json:"items"`
	Metadata PaginationMetadata                    `json:"metadata"`
}

// GetItems returns the value of Items.
func (s *OrganizationsPagedResponse) GetItems() []OrganizationsPagedResponseItemsItem {
	return s.Items
}

// GetMetadata returns the value of Metadata.
func (s *OrganizationsPagedResponse) GetMetadata() PaginationMetadata {
	return s.Metadata
}

// SetItems sets the value of Items.
func (s *OrganizationsPagedResponse) SetItems(val []OrganizationsPagedResponseItemsItem) {
	s.Items = val
}

// SetMetadata sets the value of Metadata.
func (s *OrganizationsPagedResponse) SetMetadata(val PaginationMetadata) {
	s.Metadata = val
}

// Merged schema.
type OrganizationsPagedResponseItemsItem struct {
	ID        OptUUID `json:"id"`
	Name      string  `json:"name"`
	Subdomain string  `json:"subdomain"`
}

// GetID returns the value of ID.
func (s *OrganizationsPagedResponseItemsItem) GetID() OptUUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *OrganizationsPagedResponseItemsItem) GetName() string {
	return s.Name
}

// GetSubdomain returns the value of Subdomain.
func (s *OrganizationsPagedResponseItemsItem) GetSubdomain() string {
	return s.Subdomain
}

// SetID sets the value of ID.
func (s *OrganizationsPagedResponseItemsItem) SetID(val OptUUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *OrganizationsPagedResponseItemsItem) SetName(val string) {
	s.Name = val
}

// SetSubdomain sets the value of Subdomain.
func (s *OrganizationsPagedResponseItemsItem) SetSubdomain(val string) {
	s.Subdomain = val
}

// Ref: #/components/schemas/PaginationMetadata
type PaginationMetadata struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
	Total  int32 `json:"total"`
}

// GetOffset returns the value of Offset.
func (s *PaginationMetadata) GetOffset() int32 {
	return s.Offset
}

// GetLimit returns the value of Limit.
func (s *PaginationMetadata) GetLimit() int32 {
	return s.Limit
}

// GetTotal returns the value of Total.
func (s *PaginationMetadata) GetTotal() int32 {
	return s.Total
}

// SetOffset sets the value of Offset.
func (s *PaginationMetadata) SetOffset(val int32) {
	s.Offset = val
}

// SetLimit sets the value of Limit.
func (s *PaginationMetadata) SetLimit(val int32) {
	s.Limit = val
}

// SetTotal sets the value of Total.
func (s *PaginationMetadata) SetTotal(val int32) {
	s.Total = val
}

// Ref: #/components/schemas/StorageSpace
type StorageSpace struct {
	CreatedAt time.Time `json:"created_at"`
	// ID of the employee who created the entity.
	CreatedBy NilUUID     `json:"created_by"`
	UpdatedAt NilDateTime `json:"updated_at"`
	// ID of the employee who last updated the entity.
	UpdatedBy NilUUID `json:"updated_by"`
}

// GetCreatedAt returns the value of CreatedAt.
func (s *StorageSpace) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCreatedBy returns the value of CreatedBy.
func (s *StorageSpace) GetCreatedBy() NilUUID {
	return s.CreatedBy
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *StorageSpace) GetUpdatedAt() NilDateTime {
	return s.UpdatedAt
}

// GetUpdatedBy returns the value of UpdatedBy.
func (s *StorageSpace) GetUpdatedBy() NilUUID {
	return s.UpdatedBy
}

// SetCreatedAt sets the value of CreatedAt.
func (s *StorageSpace) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetCreatedBy sets the value of CreatedBy.
func (s *StorageSpace) SetCreatedBy(val NilUUID) {
	s.CreatedBy = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *StorageSpace) SetUpdatedAt(val NilDateTime) {
	s.UpdatedAt = val
}

// SetUpdatedBy sets the value of UpdatedBy.
func (s *StorageSpace) SetUpdatedBy(val NilUUID) {
	s.UpdatedBy = val
}

// Ref: #/components/schemas/Unit
type Unit struct {
	CreatedAt time.Time `json:"created_at"`
	// ID of the employee who created the entity.
	CreatedBy NilUUID     `json:"created_by"`
	UpdatedAt NilDateTime `json:"updated_at"`
	// ID of the employee who last updated the entity.
	UpdatedBy NilUUID `json:"updated_by"`
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Unit) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCreatedBy returns the value of CreatedBy.
func (s *Unit) GetCreatedBy() NilUUID {
	return s.CreatedBy
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *Unit) GetUpdatedAt() NilDateTime {
	return s.UpdatedAt
}

// GetUpdatedBy returns the value of UpdatedBy.
func (s *Unit) GetUpdatedBy() NilUUID {
	return s.UpdatedBy
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Unit) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetCreatedBy sets the value of CreatedBy.
func (s *Unit) SetCreatedBy(val NilUUID) {
	s.CreatedBy = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *Unit) SetUpdatedAt(val NilDateTime) {
	s.UpdatedAt = val
}

// SetUpdatedBy sets the value of UpdatedBy.
func (s *Unit) SetUpdatedBy(val NilUUID) {
	s.UpdatedBy = val
}
