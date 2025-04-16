// Code generated by ogen, DO NOT EDIT.

package api

import (
	"fmt"

	"github.com/google/uuid"
)

func (s *DefaultErrorStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// Merged schema.
// Ref: #/components/schemas/CreateItemRequest
type CreateItemRequest struct {
	ID          OptUUID           `json:"id"`
	Name        string            `json:"name"`
	Description OptNilString      `json:"description"`
	Variants    []ItemVariantBase `json:"variants"`
}

// GetID returns the value of ID.
func (s *CreateItemRequest) GetID() OptUUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *CreateItemRequest) GetName() string {
	return s.Name
}

// GetDescription returns the value of Description.
func (s *CreateItemRequest) GetDescription() OptNilString {
	return s.Description
}

// GetVariants returns the value of Variants.
func (s *CreateItemRequest) GetVariants() []ItemVariantBase {
	return s.Variants
}

// SetID sets the value of ID.
func (s *CreateItemRequest) SetID(val OptUUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *CreateItemRequest) SetName(val string) {
	s.Name = val
}

// SetDescription sets the value of Description.
func (s *CreateItemRequest) SetDescription(val OptNilString) {
	s.Description = val
}

// SetVariants sets the value of Variants.
func (s *CreateItemRequest) SetVariants(val []ItemVariantBase) {
	s.Variants = val
}

// Ref: #/components/schemas/CreateItemResponse
type CreateItemResponse struct {
	Data ItemFull `json:"data"`
}

// GetData returns the value of Data.
func (s *CreateItemResponse) GetData() ItemFull {
	return s.Data
}

// SetData sets the value of Data.
func (s *CreateItemResponse) SetData(val ItemFull) {
	s.Data = val
}

// Ref: #/components/schemas/CreateOrganizationRequest
type CreateOrganizationRequest struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Subdomain string    `json:"subdomain"`
}

// GetID returns the value of ID.
func (s *CreateOrganizationRequest) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *CreateOrganizationRequest) GetName() string {
	return s.Name
}

// GetSubdomain returns the value of Subdomain.
func (s *CreateOrganizationRequest) GetSubdomain() string {
	return s.Subdomain
}

// SetID sets the value of ID.
func (s *CreateOrganizationRequest) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *CreateOrganizationRequest) SetName(val string) {
	s.Name = val
}

// SetSubdomain sets the value of Subdomain.
func (s *CreateOrganizationRequest) SetSubdomain(val string) {
	s.Subdomain = val
}

// Ref: #/components/schemas/CreateOrganizationResponse
type CreateOrganizationResponse struct {
	Data Organization `json:"data"`
}

// GetData returns the value of Data.
func (s *CreateOrganizationResponse) GetData() Organization {
	return s.Data
}

// SetData sets the value of Data.
func (s *CreateOrganizationResponse) SetData(val Organization) {
	s.Data = val
}

// Ref: #/components/schemas/CreateOrganizationUnitRequest
type CreateOrganizationUnitRequest struct {
	Name    string       `json:"name"`
	Alias   string       `json:"alias"`
	Address OptNilString `json:"address"`
}

// GetName returns the value of Name.
func (s *CreateOrganizationUnitRequest) GetName() string {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *CreateOrganizationUnitRequest) GetAlias() string {
	return s.Alias
}

// GetAddress returns the value of Address.
func (s *CreateOrganizationUnitRequest) GetAddress() OptNilString {
	return s.Address
}

// SetName sets the value of Name.
func (s *CreateOrganizationUnitRequest) SetName(val string) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *CreateOrganizationUnitRequest) SetAlias(val string) {
	s.Alias = val
}

// SetAddress sets the value of Address.
func (s *CreateOrganizationUnitRequest) SetAddress(val OptNilString) {
	s.Address = val
}

// Ref: #/components/schemas/CreateOrganizationUnitResponse
type CreateOrganizationUnitResponse struct {
	Data Unit `json:"data"`
}

// GetData returns the value of Data.
func (s *CreateOrganizationUnitResponse) GetData() Unit {
	return s.Data
}

// SetData sets the value of Data.
func (s *CreateOrganizationUnitResponse) SetData(val Unit) {
	s.Data = val
}

// Ref: #/components/schemas/CreateStorageGroupRequest
type CreateStorageGroupRequest struct {
	ParentId OptNilUUID `json:"parentId"`
	Name     string     `json:"name"`
	Alias    string     `json:"alias"`
	UnitId   uuid.UUID  `json:"unitId"`
}

// GetParentId returns the value of ParentId.
func (s *CreateStorageGroupRequest) GetParentId() OptNilUUID {
	return s.ParentId
}

// GetName returns the value of Name.
func (s *CreateStorageGroupRequest) GetName() string {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *CreateStorageGroupRequest) GetAlias() string {
	return s.Alias
}

// GetUnitId returns the value of UnitId.
func (s *CreateStorageGroupRequest) GetUnitId() uuid.UUID {
	return s.UnitId
}

// SetParentId sets the value of ParentId.
func (s *CreateStorageGroupRequest) SetParentId(val OptNilUUID) {
	s.ParentId = val
}

// SetName sets the value of Name.
func (s *CreateStorageGroupRequest) SetName(val string) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *CreateStorageGroupRequest) SetAlias(val string) {
	s.Alias = val
}

// SetUnitId sets the value of UnitId.
func (s *CreateStorageGroupRequest) SetUnitId(val uuid.UUID) {
	s.UnitId = val
}

// Ref: #/components/schemas/CreateStorageGroupResponse
type CreateStorageGroupResponse struct {
	Data StorageGroup `json:"data"`
}

// GetData returns the value of Data.
func (s *CreateStorageGroupResponse) GetData() StorageGroup {
	return s.Data
}

// SetData sets the value of Data.
func (s *CreateStorageGroupResponse) SetData(val StorageGroup) {
	s.Data = val
}

// DefaultErrorStatusCode wraps Error with StatusCode.
type DefaultErrorStatusCode struct {
	StatusCode int
	Response   Error
}

// GetStatusCode returns the value of StatusCode.
func (s *DefaultErrorStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *DefaultErrorStatusCode) GetResponse() Error {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *DefaultErrorStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *DefaultErrorStatusCode) SetResponse(val Error) {
	s.Response = val
}

// DeleteItemOK is response for DeleteItem operation.
type DeleteItemOK struct{}

// DeleteOrganizationOK is response for DeleteOrganization operation.
type DeleteOrganizationOK struct{}

// DeleteOrganizationUnitOK is response for DeleteOrganizationUnit operation.
type DeleteOrganizationUnitOK struct{}

// DeleteStorageGroupOK is response for DeleteStorageGroup operation.
type DeleteStorageGroupOK struct{}

// Represents error object.
// Ref: #/components/schemas/Error
type Error struct {
	ErrorID string `json:"error_id"`
	Message string `json:"message"`
}

// GetErrorID returns the value of ErrorID.
func (s *Error) GetErrorID() string {
	return s.ErrorID
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// SetErrorID sets the value of ErrorID.
func (s *Error) SetErrorID(val string) {
	s.ErrorID = val
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

// Ref: #/components/schemas/GetItemByIdResponse
type GetItemByIdResponse struct {
	Data ItemFull `json:"data"`
}

// GetData returns the value of Data.
func (s *GetItemByIdResponse) GetData() ItemFull {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetItemByIdResponse) SetData(val ItemFull) {
	s.Data = val
}

// Ref: #/components/schemas/GetItemsResponse
type GetItemsResponse struct {
	Data []Item `json:"data"`
}

// GetData returns the value of Data.
func (s *GetItemsResponse) GetData() []Item {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetItemsResponse) SetData(val []Item) {
	s.Data = val
}

// Ref: #/components/schemas/GetOrganizationByIdResponse
type GetOrganizationByIdResponse struct {
	Data Organization `json:"data"`
}

// GetData returns the value of Data.
func (s *GetOrganizationByIdResponse) GetData() Organization {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetOrganizationByIdResponse) SetData(val Organization) {
	s.Data = val
}

// Ref: #/components/schemas/GetOrganizationUnitByIdResponse
type GetOrganizationUnitByIdResponse struct {
	Data Unit `json:"data"`
}

// GetData returns the value of Data.
func (s *GetOrganizationUnitByIdResponse) GetData() Unit {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetOrganizationUnitByIdResponse) SetData(val Unit) {
	s.Data = val
}

// Ref: #/components/schemas/GetOrganizationUnitsResponse
type GetOrganizationUnitsResponse struct {
	Data []Unit `json:"data"`
}

// GetData returns the value of Data.
func (s *GetOrganizationUnitsResponse) GetData() []Unit {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetOrganizationUnitsResponse) SetData(val []Unit) {
	s.Data = val
}

// Ref: #/components/schemas/GetOrganizationsResponse
type GetOrganizationsResponse struct {
	Data []Organization `json:"data"`
}

// GetData returns the value of Data.
func (s *GetOrganizationsResponse) GetData() []Organization {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetOrganizationsResponse) SetData(val []Organization) {
	s.Data = val
}

// Ref: #/components/schemas/GetStorageGroupByIdResponse
type GetStorageGroupByIdResponse struct {
	Data StorageGroup `json:"data"`
}

// GetData returns the value of Data.
func (s *GetStorageGroupByIdResponse) GetData() StorageGroup {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetStorageGroupByIdResponse) SetData(val StorageGroup) {
	s.Data = val
}

// Ref: #/components/schemas/GetStorageGroupsResponse
type GetStorageGroupsResponse struct {
	Data []StorageGroup `json:"data"`
}

// GetData returns the value of Data.
func (s *GetStorageGroupsResponse) GetData() []StorageGroup {
	return s.Data
}

// SetData sets the value of Data.
func (s *GetStorageGroupsResponse) SetData(val []StorageGroup) {
	s.Data = val
}

// Merged schema.
// Ref: #/components/schemas/Item
type Item struct {
	// Merged property.
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// Merged property.
	Description NilString `json:"description"`
}

// GetID returns the value of ID.
func (s *Item) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *Item) GetName() string {
	return s.Name
}

// GetDescription returns the value of Description.
func (s *Item) GetDescription() NilString {
	return s.Description
}

// SetID sets the value of ID.
func (s *Item) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *Item) SetName(val string) {
	s.Name = val
}

// SetDescription sets the value of Description.
func (s *Item) SetDescription(val NilString) {
	s.Description = val
}

// Merged schema.
// Ref: #/components/schemas/ItemFull
type ItemFull struct {
	// Merged property.
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// Merged property.
	Description NilString     `json:"description"`
	Variants    []ItemVariant `json:"variants"`
}

// GetID returns the value of ID.
func (s *ItemFull) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *ItemFull) GetName() string {
	return s.Name
}

// GetDescription returns the value of Description.
func (s *ItemFull) GetDescription() NilString {
	return s.Description
}

// GetVariants returns the value of Variants.
func (s *ItemFull) GetVariants() []ItemVariant {
	return s.Variants
}

// SetID sets the value of ID.
func (s *ItemFull) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *ItemFull) SetName(val string) {
	s.Name = val
}

// SetDescription sets the value of Description.
func (s *ItemFull) SetDescription(val NilString) {
	s.Description = val
}

// SetVariants sets the value of Variants.
func (s *ItemFull) SetVariants(val []ItemVariant) {
	s.Variants = val
}

// Merged schema.
// Ref: #/components/schemas/ItemVariant
type ItemVariant struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// Merged property.
	Article NilString `json:"article"`
	// Merged property.
	Ean13 NilInt `json:"ean13"`
}

// GetID returns the value of ID.
func (s *ItemVariant) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *ItemVariant) GetName() string {
	return s.Name
}

// GetArticle returns the value of Article.
func (s *ItemVariant) GetArticle() NilString {
	return s.Article
}

// GetEan13 returns the value of Ean13.
func (s *ItemVariant) GetEan13() NilInt {
	return s.Ean13
}

// SetID sets the value of ID.
func (s *ItemVariant) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *ItemVariant) SetName(val string) {
	s.Name = val
}

// SetArticle sets the value of Article.
func (s *ItemVariant) SetArticle(val NilString) {
	s.Article = val
}

// SetEan13 sets the value of Ean13.
func (s *ItemVariant) SetEan13(val NilInt) {
	s.Ean13 = val
}

// Ref: #/components/schemas/ItemVariantBase
type ItemVariantBase struct {
	Name    string       `json:"name"`
	Article OptNilString `json:"article"`
	Ean13   OptNilInt    `json:"ean13"`
}

// GetName returns the value of Name.
func (s *ItemVariantBase) GetName() string {
	return s.Name
}

// GetArticle returns the value of Article.
func (s *ItemVariantBase) GetArticle() OptNilString {
	return s.Article
}

// GetEan13 returns the value of Ean13.
func (s *ItemVariantBase) GetEan13() OptNilInt {
	return s.Ean13
}

// SetName sets the value of Name.
func (s *ItemVariantBase) SetName(val string) {
	s.Name = val
}

// SetArticle sets the value of Article.
func (s *ItemVariantBase) SetArticle(val OptNilString) {
	s.Article = val
}

// SetEan13 sets the value of Ean13.
func (s *ItemVariantBase) SetEan13(val OptNilInt) {
	s.Ean13 = val
}

// Ref: #/components/schemas/ItemVariantPatchInItem
type ItemVariantPatchInItem struct {
	ID      uuid.UUID    `json:"id"`
	Name    OptString    `json:"name"`
	Article OptNilString `json:"article"`
	Ean13   OptNilInt    `json:"ean13"`
}

// GetID returns the value of ID.
func (s *ItemVariantPatchInItem) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *ItemVariantPatchInItem) GetName() OptString {
	return s.Name
}

// GetArticle returns the value of Article.
func (s *ItemVariantPatchInItem) GetArticle() OptNilString {
	return s.Article
}

// GetEan13 returns the value of Ean13.
func (s *ItemVariantPatchInItem) GetEan13() OptNilInt {
	return s.Ean13
}

// SetID sets the value of ID.
func (s *ItemVariantPatchInItem) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *ItemVariantPatchInItem) SetName(val OptString) {
	s.Name = val
}

// SetArticle sets the value of Article.
func (s *ItemVariantPatchInItem) SetArticle(val OptNilString) {
	s.Article = val
}

// SetEan13 sets the value of Ean13.
func (s *ItemVariantPatchInItem) SetEan13(val OptNilInt) {
	s.Ean13 = val
}

// NewNilInt returns new NilInt with value set to v.
func NewNilInt(v int) NilInt {
	return NilInt{
		Value: v,
	}
}

// NilInt is nullable int.
type NilInt struct {
	Value int
	Null  bool
}

// SetTo sets value to v.
func (o *NilInt) SetTo(v int) {
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o NilInt) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *NilInt) SetToNull() {
	o.Null = true
	var v int
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o NilInt) Get() (v int, ok bool) {
	if o.Null {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o NilInt) Or(d int) int {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewNilString returns new NilString with value set to v.
func NewNilString(v string) NilString {
	return NilString{
		Value: v,
	}
}

// NilString is nullable string.
type NilString struct {
	Value string
	Null  bool
}

// SetTo sets value to v.
func (o *NilString) SetTo(v string) {
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o NilString) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *NilString) SetToNull() {
	o.Null = true
	var v string
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o NilString) Get() (v string, ok bool) {
	if o.Null {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o NilString) Or(d string) string {
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

// NewOptNilInt returns new OptNilInt with value set to v.
func NewOptNilInt(v int) OptNilInt {
	return OptNilInt{
		Value: v,
		Set:   true,
	}
}

// OptNilInt is optional nullable int.
type OptNilInt struct {
	Value int
	Set   bool
	Null  bool
}

// IsSet returns true if OptNilInt was set.
func (o OptNilInt) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptNilInt) Reset() {
	var v int
	o.Value = v
	o.Set = false
	o.Null = false
}

// SetTo sets value to v.
func (o *OptNilInt) SetTo(v int) {
	o.Set = true
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o OptNilInt) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *OptNilInt) SetToNull() {
	o.Set = true
	o.Null = true
	var v int
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptNilInt) Get() (v int, ok bool) {
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptNilInt) Or(d int) int {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptNilString returns new OptNilString with value set to v.
func NewOptNilString(v string) OptNilString {
	return OptNilString{
		Value: v,
		Set:   true,
	}
}

// OptNilString is optional nullable string.
type OptNilString struct {
	Value string
	Set   bool
	Null  bool
}

// IsSet returns true if OptNilString was set.
func (o OptNilString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptNilString) Reset() {
	var v string
	o.Value = v
	o.Set = false
	o.Null = false
}

// SetTo sets value to v.
func (o *OptNilString) SetTo(v string) {
	o.Set = true
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o OptNilString) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *OptNilString) SetToNull() {
	o.Set = true
	o.Null = true
	var v string
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptNilString) Get() (v string, ok bool) {
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptNilString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptNilUUID returns new OptNilUUID with value set to v.
func NewOptNilUUID(v uuid.UUID) OptNilUUID {
	return OptNilUUID{
		Value: v,
		Set:   true,
	}
}

// OptNilUUID is optional nullable uuid.UUID.
type OptNilUUID struct {
	Value uuid.UUID
	Set   bool
	Null  bool
}

// IsSet returns true if OptNilUUID was set.
func (o OptNilUUID) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptNilUUID) Reset() {
	var v uuid.UUID
	o.Value = v
	o.Set = false
	o.Null = false
}

// SetTo sets value to v.
func (o *OptNilUUID) SetTo(v uuid.UUID) {
	o.Set = true
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o OptNilUUID) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *OptNilUUID) SetToNull() {
	o.Set = true
	o.Null = true
	var v uuid.UUID
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptNilUUID) Get() (v uuid.UUID, ok bool) {
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptNilUUID) Or(d uuid.UUID) uuid.UUID {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
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
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Subdomain string    `json:"subdomain"`
}

// GetID returns the value of ID.
func (s *Organization) GetID() uuid.UUID {
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
func (s *Organization) SetID(val uuid.UUID) {
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

// Ref: #/components/schemas/PatchItemRequest
type PatchItemRequest struct {
	Name        OptString                `json:"name"`
	Description OptNilString             `json:"description"`
	Variants    []ItemVariantPatchInItem `json:"variants"`
}

// GetName returns the value of Name.
func (s *PatchItemRequest) GetName() OptString {
	return s.Name
}

// GetDescription returns the value of Description.
func (s *PatchItemRequest) GetDescription() OptNilString {
	return s.Description
}

// GetVariants returns the value of Variants.
func (s *PatchItemRequest) GetVariants() []ItemVariantPatchInItem {
	return s.Variants
}

// SetName sets the value of Name.
func (s *PatchItemRequest) SetName(val OptString) {
	s.Name = val
}

// SetDescription sets the value of Description.
func (s *PatchItemRequest) SetDescription(val OptNilString) {
	s.Description = val
}

// SetVariants sets the value of Variants.
func (s *PatchItemRequest) SetVariants(val []ItemVariantPatchInItem) {
	s.Variants = val
}

// Ref: #/components/schemas/PatchItemResponse
type PatchItemResponse struct {
	Data ItemFull `json:"data"`
}

// GetData returns the value of Data.
func (s *PatchItemResponse) GetData() ItemFull {
	return s.Data
}

// SetData sets the value of Data.
func (s *PatchItemResponse) SetData(val ItemFull) {
	s.Data = val
}

// Ref: #/components/schemas/PatchOrganizationRequest
type PatchOrganizationRequest struct {
	Name      OptString `json:"name"`
	Subdomain OptString `json:"subdomain"`
}

// GetName returns the value of Name.
func (s *PatchOrganizationRequest) GetName() OptString {
	return s.Name
}

// GetSubdomain returns the value of Subdomain.
func (s *PatchOrganizationRequest) GetSubdomain() OptString {
	return s.Subdomain
}

// SetName sets the value of Name.
func (s *PatchOrganizationRequest) SetName(val OptString) {
	s.Name = val
}

// SetSubdomain sets the value of Subdomain.
func (s *PatchOrganizationRequest) SetSubdomain(val OptString) {
	s.Subdomain = val
}

// Ref: #/components/schemas/PatchOrganizationResponse
type PatchOrganizationResponse struct {
	Data []Organization `json:"data"`
}

// GetData returns the value of Data.
func (s *PatchOrganizationResponse) GetData() []Organization {
	return s.Data
}

// SetData sets the value of Data.
func (s *PatchOrganizationResponse) SetData(val []Organization) {
	s.Data = val
}

// Ref: #/components/schemas/PatchOrganizationUnitRequest
type PatchOrganizationUnitRequest struct {
	Name    OptString    `json:"name"`
	Alias   OptString    `json:"alias"`
	Address OptNilString `json:"address"`
}

// GetName returns the value of Name.
func (s *PatchOrganizationUnitRequest) GetName() OptString {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *PatchOrganizationUnitRequest) GetAlias() OptString {
	return s.Alias
}

// GetAddress returns the value of Address.
func (s *PatchOrganizationUnitRequest) GetAddress() OptNilString {
	return s.Address
}

// SetName sets the value of Name.
func (s *PatchOrganizationUnitRequest) SetName(val OptString) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *PatchOrganizationUnitRequest) SetAlias(val OptString) {
	s.Alias = val
}

// SetAddress sets the value of Address.
func (s *PatchOrganizationUnitRequest) SetAddress(val OptNilString) {
	s.Address = val
}

// Ref: #/components/schemas/PatchOrganizationUnitResponse
type PatchOrganizationUnitResponse struct {
	Data []Unit `json:"data"`
}

// GetData returns the value of Data.
func (s *PatchOrganizationUnitResponse) GetData() []Unit {
	return s.Data
}

// SetData sets the value of Data.
func (s *PatchOrganizationUnitResponse) SetData(val []Unit) {
	s.Data = val
}

// Ref: #/components/schemas/PatchStorageGroupRequest
type PatchStorageGroupRequest struct {
	ParentId OptNilUUID `json:"parentId"`
	Name     OptString  `json:"name"`
	Alias    OptString  `json:"alias"`
	UnitId   OptUUID    `json:"unitId"`
}

// GetParentId returns the value of ParentId.
func (s *PatchStorageGroupRequest) GetParentId() OptNilUUID {
	return s.ParentId
}

// GetName returns the value of Name.
func (s *PatchStorageGroupRequest) GetName() OptString {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *PatchStorageGroupRequest) GetAlias() OptString {
	return s.Alias
}

// GetUnitId returns the value of UnitId.
func (s *PatchStorageGroupRequest) GetUnitId() OptUUID {
	return s.UnitId
}

// SetParentId sets the value of ParentId.
func (s *PatchStorageGroupRequest) SetParentId(val OptNilUUID) {
	s.ParentId = val
}

// SetName sets the value of Name.
func (s *PatchStorageGroupRequest) SetName(val OptString) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *PatchStorageGroupRequest) SetAlias(val OptString) {
	s.Alias = val
}

// SetUnitId sets the value of UnitId.
func (s *PatchStorageGroupRequest) SetUnitId(val OptUUID) {
	s.UnitId = val
}

// Ref: #/components/schemas/PatchStorageGroupResponse
type PatchStorageGroupResponse struct {
	Data []StorageGroup `json:"data"`
}

// GetData returns the value of Data.
func (s *PatchStorageGroupResponse) GetData() []StorageGroup {
	return s.Data
}

// SetData sets the value of Data.
func (s *PatchStorageGroupResponse) SetData(val []StorageGroup) {
	s.Data = val
}

// Merged schema.
// Ref: #/components/schemas/StorageGroup
type StorageGroup struct {
	ID uuid.UUID `json:"id"`
	// Merged property.
	ParentId NilUUID   `json:"parentId"`
	Name     string    `json:"name"`
	Alias    string    `json:"alias"`
	UnitId   uuid.UUID `json:"unitId"`
}

// GetID returns the value of ID.
func (s *StorageGroup) GetID() uuid.UUID {
	return s.ID
}

// GetParentId returns the value of ParentId.
func (s *StorageGroup) GetParentId() NilUUID {
	return s.ParentId
}

// GetName returns the value of Name.
func (s *StorageGroup) GetName() string {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *StorageGroup) GetAlias() string {
	return s.Alias
}

// GetUnitId returns the value of UnitId.
func (s *StorageGroup) GetUnitId() uuid.UUID {
	return s.UnitId
}

// SetID sets the value of ID.
func (s *StorageGroup) SetID(val uuid.UUID) {
	s.ID = val
}

// SetParentId sets the value of ParentId.
func (s *StorageGroup) SetParentId(val NilUUID) {
	s.ParentId = val
}

// SetName sets the value of Name.
func (s *StorageGroup) SetName(val string) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *StorageGroup) SetAlias(val string) {
	s.Alias = val
}

// SetUnitId sets the value of UnitId.
func (s *StorageGroup) SetUnitId(val uuid.UUID) {
	s.UnitId = val
}

// Merged schema.
// Ref: #/components/schemas/Unit
type Unit struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Alias string    `json:"alias"`
	// Merged property.
	Address NilString `json:"address"`
}

// GetID returns the value of ID.
func (s *Unit) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *Unit) GetName() string {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *Unit) GetAlias() string {
	return s.Alias
}

// GetAddress returns the value of Address.
func (s *Unit) GetAddress() NilString {
	return s.Address
}

// SetID sets the value of ID.
func (s *Unit) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *Unit) SetName(val string) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *Unit) SetAlias(val string) {
	s.Alias = val
}

// SetAddress sets the value of Address.
func (s *Unit) SetAddress(val NilString) {
	s.Address = val
}

// Merged schema.
// Ref: #/components/schemas/UpdateItemRequest
type UpdateItemRequest struct {
	ID          OptUUID      `json:"id"`
	Name        string       `json:"name"`
	Description OptNilString `json:"description"`
	// Merged property.
	Variants []UpdateItemRequestVariantsItem `json:"variants"`
}

// GetID returns the value of ID.
func (s *UpdateItemRequest) GetID() OptUUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *UpdateItemRequest) GetName() string {
	return s.Name
}

// GetDescription returns the value of Description.
func (s *UpdateItemRequest) GetDescription() OptNilString {
	return s.Description
}

// GetVariants returns the value of Variants.
func (s *UpdateItemRequest) GetVariants() []UpdateItemRequestVariantsItem {
	return s.Variants
}

// SetID sets the value of ID.
func (s *UpdateItemRequest) SetID(val OptUUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *UpdateItemRequest) SetName(val string) {
	s.Name = val
}

// SetDescription sets the value of Description.
func (s *UpdateItemRequest) SetDescription(val OptNilString) {
	s.Description = val
}

// SetVariants sets the value of Variants.
func (s *UpdateItemRequest) SetVariants(val []UpdateItemRequestVariantsItem) {
	s.Variants = val
}

// Merged schema.
type UpdateItemRequestVariantsItem struct {
	ID uuid.UUID `json:"id"`
	// Merged property.
	Name string `json:"name"`
	// Merged property.
	Article OptNilString `json:"article"`
	// Merged property.
	Ean13 OptNilInt `json:"ean13"`
}

// GetID returns the value of ID.
func (s *UpdateItemRequestVariantsItem) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *UpdateItemRequestVariantsItem) GetName() string {
	return s.Name
}

// GetArticle returns the value of Article.
func (s *UpdateItemRequestVariantsItem) GetArticle() OptNilString {
	return s.Article
}

// GetEan13 returns the value of Ean13.
func (s *UpdateItemRequestVariantsItem) GetEan13() OptNilInt {
	return s.Ean13
}

// SetID sets the value of ID.
func (s *UpdateItemRequestVariantsItem) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *UpdateItemRequestVariantsItem) SetName(val string) {
	s.Name = val
}

// SetArticle sets the value of Article.
func (s *UpdateItemRequestVariantsItem) SetArticle(val OptNilString) {
	s.Article = val
}

// SetEan13 sets the value of Ean13.
func (s *UpdateItemRequestVariantsItem) SetEan13(val OptNilInt) {
	s.Ean13 = val
}

// Ref: #/components/schemas/UpdateItemResponse
type UpdateItemResponse struct {
	Data ItemFull `json:"data"`
}

// GetData returns the value of Data.
func (s *UpdateItemResponse) GetData() ItemFull {
	return s.Data
}

// SetData sets the value of Data.
func (s *UpdateItemResponse) SetData(val ItemFull) {
	s.Data = val
}

// Ref: #/components/schemas/UpdateOrganizationRequest
type UpdateOrganizationRequest struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Subdomain string    `json:"subdomain"`
}

// GetID returns the value of ID.
func (s *UpdateOrganizationRequest) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *UpdateOrganizationRequest) GetName() string {
	return s.Name
}

// GetSubdomain returns the value of Subdomain.
func (s *UpdateOrganizationRequest) GetSubdomain() string {
	return s.Subdomain
}

// SetID sets the value of ID.
func (s *UpdateOrganizationRequest) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *UpdateOrganizationRequest) SetName(val string) {
	s.Name = val
}

// SetSubdomain sets the value of Subdomain.
func (s *UpdateOrganizationRequest) SetSubdomain(val string) {
	s.Subdomain = val
}

// Ref: #/components/schemas/UpdateOrganizationResponse
type UpdateOrganizationResponse struct {
	Data []Organization `json:"data"`
}

// GetData returns the value of Data.
func (s *UpdateOrganizationResponse) GetData() []Organization {
	return s.Data
}

// SetData sets the value of Data.
func (s *UpdateOrganizationResponse) SetData(val []Organization) {
	s.Data = val
}

// Ref: #/components/schemas/UpdateOrganizationUnitRequest
type UpdateOrganizationUnitRequest struct {
	Name    string       `json:"name"`
	Alias   string       `json:"alias"`
	Address OptNilString `json:"address"`
}

// GetName returns the value of Name.
func (s *UpdateOrganizationUnitRequest) GetName() string {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *UpdateOrganizationUnitRequest) GetAlias() string {
	return s.Alias
}

// GetAddress returns the value of Address.
func (s *UpdateOrganizationUnitRequest) GetAddress() OptNilString {
	return s.Address
}

// SetName sets the value of Name.
func (s *UpdateOrganizationUnitRequest) SetName(val string) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *UpdateOrganizationUnitRequest) SetAlias(val string) {
	s.Alias = val
}

// SetAddress sets the value of Address.
func (s *UpdateOrganizationUnitRequest) SetAddress(val OptNilString) {
	s.Address = val
}

// Ref: #/components/schemas/UpdateOrganizationUnitResponse
type UpdateOrganizationUnitResponse struct {
	Data []Unit `json:"data"`
}

// GetData returns the value of Data.
func (s *UpdateOrganizationUnitResponse) GetData() []Unit {
	return s.Data
}

// SetData sets the value of Data.
func (s *UpdateOrganizationUnitResponse) SetData(val []Unit) {
	s.Data = val
}

// Ref: #/components/schemas/UpdateStorageGroupRequest
type UpdateStorageGroupRequest struct {
	ParentId OptNilUUID `json:"parentId"`
	Name     string     `json:"name"`
	Alias    string     `json:"alias"`
	UnitId   uuid.UUID  `json:"unitId"`
}

// GetParentId returns the value of ParentId.
func (s *UpdateStorageGroupRequest) GetParentId() OptNilUUID {
	return s.ParentId
}

// GetName returns the value of Name.
func (s *UpdateStorageGroupRequest) GetName() string {
	return s.Name
}

// GetAlias returns the value of Alias.
func (s *UpdateStorageGroupRequest) GetAlias() string {
	return s.Alias
}

// GetUnitId returns the value of UnitId.
func (s *UpdateStorageGroupRequest) GetUnitId() uuid.UUID {
	return s.UnitId
}

// SetParentId sets the value of ParentId.
func (s *UpdateStorageGroupRequest) SetParentId(val OptNilUUID) {
	s.ParentId = val
}

// SetName sets the value of Name.
func (s *UpdateStorageGroupRequest) SetName(val string) {
	s.Name = val
}

// SetAlias sets the value of Alias.
func (s *UpdateStorageGroupRequest) SetAlias(val string) {
	s.Alias = val
}

// SetUnitId sets the value of UnitId.
func (s *UpdateStorageGroupRequest) SetUnitId(val uuid.UUID) {
	s.UnitId = val
}

// Ref: #/components/schemas/UpdateStorageGroupResponse
type UpdateStorageGroupResponse struct {
	Data []StorageGroup `json:"data"`
}

// GetData returns the value of Data.
func (s *UpdateStorageGroupResponse) GetData() []StorageGroup {
	return s.Data
}

// SetData sets the value of Data.
func (s *UpdateStorageGroupResponse) SetData(val []StorageGroup) {
	s.Data = val
}
