// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/hay-kot/homebox/backend/internal/data/ent/attachment"
	"github.com/hay-kot/homebox/backend/internal/data/ent/authtokens"
	"github.com/hay-kot/homebox/backend/internal/data/ent/document"
	"github.com/hay-kot/homebox/backend/internal/data/ent/group"
	"github.com/hay-kot/homebox/backend/internal/data/ent/groupinvitationtoken"
	"github.com/hay-kot/homebox/backend/internal/data/ent/item"
	"github.com/hay-kot/homebox/backend/internal/data/ent/itemfield"
	"github.com/hay-kot/homebox/backend/internal/data/ent/label"
	"github.com/hay-kot/homebox/backend/internal/data/ent/location"
	"github.com/hay-kot/homebox/backend/internal/data/ent/maintenanceentry"
	"github.com/hay-kot/homebox/backend/internal/data/ent/notifier"
	"github.com/hay-kot/homebox/backend/internal/data/ent/schema"
	"github.com/hay-kot/homebox/backend/internal/data/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	attachmentMixin := schema.Attachment{}.Mixin()
	attachmentMixinFields0 := attachmentMixin[0].Fields()
	_ = attachmentMixinFields0
	attachmentFields := schema.Attachment{}.Fields()
	_ = attachmentFields
	// attachmentDescCreatedAt is the schema descriptor for created_at field.
	attachmentDescCreatedAt := attachmentMixinFields0[1].Descriptor()
	// attachment.DefaultCreatedAt holds the default value on creation for the created_at field.
	attachment.DefaultCreatedAt = attachmentDescCreatedAt.Default.(func() time.Time)
	// attachmentDescUpdatedAt is the schema descriptor for updated_at field.
	attachmentDescUpdatedAt := attachmentMixinFields0[2].Descriptor()
	// attachment.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	attachment.DefaultUpdatedAt = attachmentDescUpdatedAt.Default.(func() time.Time)
	// attachment.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	attachment.UpdateDefaultUpdatedAt = attachmentDescUpdatedAt.UpdateDefault.(func() time.Time)
	// attachmentDescPrimary is the schema descriptor for primary field.
	attachmentDescPrimary := attachmentFields[1].Descriptor()
	// attachment.DefaultPrimary holds the default value on creation for the primary field.
	attachment.DefaultPrimary = attachmentDescPrimary.Default.(bool)
	// attachmentDescID is the schema descriptor for id field.
	attachmentDescID := attachmentMixinFields0[0].Descriptor()
	// attachment.DefaultID holds the default value on creation for the id field.
	attachment.DefaultID = attachmentDescID.Default.(func() uuid.UUID)
	authrolesFields := schema.AuthRoles{}.Fields()
	_ = authrolesFields
	authtokensMixin := schema.AuthTokens{}.Mixin()
	authtokensMixinFields0 := authtokensMixin[0].Fields()
	_ = authtokensMixinFields0
	authtokensFields := schema.AuthTokens{}.Fields()
	_ = authtokensFields
	// authtokensDescCreatedAt is the schema descriptor for created_at field.
	authtokensDescCreatedAt := authtokensMixinFields0[1].Descriptor()
	// authtokens.DefaultCreatedAt holds the default value on creation for the created_at field.
	authtokens.DefaultCreatedAt = authtokensDescCreatedAt.Default.(func() time.Time)
	// authtokensDescUpdatedAt is the schema descriptor for updated_at field.
	authtokensDescUpdatedAt := authtokensMixinFields0[2].Descriptor()
	// authtokens.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	authtokens.DefaultUpdatedAt = authtokensDescUpdatedAt.Default.(func() time.Time)
	// authtokens.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	authtokens.UpdateDefaultUpdatedAt = authtokensDescUpdatedAt.UpdateDefault.(func() time.Time)
	// authtokensDescExpiresAt is the schema descriptor for expires_at field.
	authtokensDescExpiresAt := authtokensFields[1].Descriptor()
	// authtokens.DefaultExpiresAt holds the default value on creation for the expires_at field.
	authtokens.DefaultExpiresAt = authtokensDescExpiresAt.Default.(func() time.Time)
	// authtokensDescID is the schema descriptor for id field.
	authtokensDescID := authtokensMixinFields0[0].Descriptor()
	// authtokens.DefaultID holds the default value on creation for the id field.
	authtokens.DefaultID = authtokensDescID.Default.(func() uuid.UUID)
	documentMixin := schema.Document{}.Mixin()
	documentMixinFields0 := documentMixin[0].Fields()
	_ = documentMixinFields0
	documentFields := schema.Document{}.Fields()
	_ = documentFields
	// documentDescCreatedAt is the schema descriptor for created_at field.
	documentDescCreatedAt := documentMixinFields0[1].Descriptor()
	// document.DefaultCreatedAt holds the default value on creation for the created_at field.
	document.DefaultCreatedAt = documentDescCreatedAt.Default.(func() time.Time)
	// documentDescUpdatedAt is the schema descriptor for updated_at field.
	documentDescUpdatedAt := documentMixinFields0[2].Descriptor()
	// document.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	document.DefaultUpdatedAt = documentDescUpdatedAt.Default.(func() time.Time)
	// document.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	document.UpdateDefaultUpdatedAt = documentDescUpdatedAt.UpdateDefault.(func() time.Time)
	// documentDescTitle is the schema descriptor for title field.
	documentDescTitle := documentFields[0].Descriptor()
	// document.TitleValidator is a validator for the "title" field. It is called by the builders before save.
	document.TitleValidator = func() func(string) error {
		validators := documentDescTitle.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(title string) error {
			for _, fn := range fns {
				if err := fn(title); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// documentDescPath is the schema descriptor for path field.
	documentDescPath := documentFields[1].Descriptor()
	// document.PathValidator is a validator for the "path" field. It is called by the builders before save.
	document.PathValidator = func() func(string) error {
		validators := documentDescPath.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(_path string) error {
			for _, fn := range fns {
				if err := fn(_path); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// documentDescID is the schema descriptor for id field.
	documentDescID := documentMixinFields0[0].Descriptor()
	// document.DefaultID holds the default value on creation for the id field.
	document.DefaultID = documentDescID.Default.(func() uuid.UUID)
	groupMixin := schema.Group{}.Mixin()
	groupMixinFields0 := groupMixin[0].Fields()
	_ = groupMixinFields0
	groupFields := schema.Group{}.Fields()
	_ = groupFields
	// groupDescCreatedAt is the schema descriptor for created_at field.
	groupDescCreatedAt := groupMixinFields0[1].Descriptor()
	// group.DefaultCreatedAt holds the default value on creation for the created_at field.
	group.DefaultCreatedAt = groupDescCreatedAt.Default.(func() time.Time)
	// groupDescUpdatedAt is the schema descriptor for updated_at field.
	groupDescUpdatedAt := groupMixinFields0[2].Descriptor()
	// group.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	group.DefaultUpdatedAt = groupDescUpdatedAt.Default.(func() time.Time)
	// group.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	group.UpdateDefaultUpdatedAt = groupDescUpdatedAt.UpdateDefault.(func() time.Time)
	// groupDescName is the schema descriptor for name field.
	groupDescName := groupFields[0].Descriptor()
	// group.NameValidator is a validator for the "name" field. It is called by the builders before save.
	group.NameValidator = func() func(string) error {
		validators := groupDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// groupDescID is the schema descriptor for id field.
	groupDescID := groupMixinFields0[0].Descriptor()
	// group.DefaultID holds the default value on creation for the id field.
	group.DefaultID = groupDescID.Default.(func() uuid.UUID)
	groupinvitationtokenMixin := schema.GroupInvitationToken{}.Mixin()
	groupinvitationtokenMixinFields0 := groupinvitationtokenMixin[0].Fields()
	_ = groupinvitationtokenMixinFields0
	groupinvitationtokenFields := schema.GroupInvitationToken{}.Fields()
	_ = groupinvitationtokenFields
	// groupinvitationtokenDescCreatedAt is the schema descriptor for created_at field.
	groupinvitationtokenDescCreatedAt := groupinvitationtokenMixinFields0[1].Descriptor()
	// groupinvitationtoken.DefaultCreatedAt holds the default value on creation for the created_at field.
	groupinvitationtoken.DefaultCreatedAt = groupinvitationtokenDescCreatedAt.Default.(func() time.Time)
	// groupinvitationtokenDescUpdatedAt is the schema descriptor for updated_at field.
	groupinvitationtokenDescUpdatedAt := groupinvitationtokenMixinFields0[2].Descriptor()
	// groupinvitationtoken.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	groupinvitationtoken.DefaultUpdatedAt = groupinvitationtokenDescUpdatedAt.Default.(func() time.Time)
	// groupinvitationtoken.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	groupinvitationtoken.UpdateDefaultUpdatedAt = groupinvitationtokenDescUpdatedAt.UpdateDefault.(func() time.Time)
	// groupinvitationtokenDescExpiresAt is the schema descriptor for expires_at field.
	groupinvitationtokenDescExpiresAt := groupinvitationtokenFields[1].Descriptor()
	// groupinvitationtoken.DefaultExpiresAt holds the default value on creation for the expires_at field.
	groupinvitationtoken.DefaultExpiresAt = groupinvitationtokenDescExpiresAt.Default.(func() time.Time)
	// groupinvitationtokenDescUses is the schema descriptor for uses field.
	groupinvitationtokenDescUses := groupinvitationtokenFields[2].Descriptor()
	// groupinvitationtoken.DefaultUses holds the default value on creation for the uses field.
	groupinvitationtoken.DefaultUses = groupinvitationtokenDescUses.Default.(int)
	// groupinvitationtokenDescID is the schema descriptor for id field.
	groupinvitationtokenDescID := groupinvitationtokenMixinFields0[0].Descriptor()
	// groupinvitationtoken.DefaultID holds the default value on creation for the id field.
	groupinvitationtoken.DefaultID = groupinvitationtokenDescID.Default.(func() uuid.UUID)
	itemMixin := schema.Item{}.Mixin()
	itemMixinFields0 := itemMixin[0].Fields()
	_ = itemMixinFields0
	itemMixinFields1 := itemMixin[1].Fields()
	_ = itemMixinFields1
	itemFields := schema.Item{}.Fields()
	_ = itemFields
	// itemDescCreatedAt is the schema descriptor for created_at field.
	itemDescCreatedAt := itemMixinFields0[1].Descriptor()
	// item.DefaultCreatedAt holds the default value on creation for the created_at field.
	item.DefaultCreatedAt = itemDescCreatedAt.Default.(func() time.Time)
	// itemDescUpdatedAt is the schema descriptor for updated_at field.
	itemDescUpdatedAt := itemMixinFields0[2].Descriptor()
	// item.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	item.DefaultUpdatedAt = itemDescUpdatedAt.Default.(func() time.Time)
	// item.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	item.UpdateDefaultUpdatedAt = itemDescUpdatedAt.UpdateDefault.(func() time.Time)
	// itemDescName is the schema descriptor for name field.
	itemDescName := itemMixinFields1[0].Descriptor()
	// item.NameValidator is a validator for the "name" field. It is called by the builders before save.
	item.NameValidator = func() func(string) error {
		validators := itemDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// itemDescDescription is the schema descriptor for description field.
	itemDescDescription := itemMixinFields1[1].Descriptor()
	// item.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	item.DescriptionValidator = itemDescDescription.Validators[0].(func(string) error)
	// itemDescImportRef is the schema descriptor for import_ref field.
	itemDescImportRef := itemFields[0].Descriptor()
	// item.ImportRefValidator is a validator for the "import_ref" field. It is called by the builders before save.
	item.ImportRefValidator = itemDescImportRef.Validators[0].(func(string) error)
	// itemDescNotes is the schema descriptor for notes field.
	itemDescNotes := itemFields[1].Descriptor()
	// item.NotesValidator is a validator for the "notes" field. It is called by the builders before save.
	item.NotesValidator = itemDescNotes.Validators[0].(func(string) error)
	// itemDescQuantity is the schema descriptor for quantity field.
	itemDescQuantity := itemFields[2].Descriptor()
	// item.DefaultQuantity holds the default value on creation for the quantity field.
	item.DefaultQuantity = itemDescQuantity.Default.(int)
	// itemDescInsured is the schema descriptor for insured field.
	itemDescInsured := itemFields[3].Descriptor()
	// item.DefaultInsured holds the default value on creation for the insured field.
	item.DefaultInsured = itemDescInsured.Default.(bool)
	// itemDescArchived is the schema descriptor for archived field.
	itemDescArchived := itemFields[4].Descriptor()
	// item.DefaultArchived holds the default value on creation for the archived field.
	item.DefaultArchived = itemDescArchived.Default.(bool)
	// itemDescAssetID is the schema descriptor for asset_id field.
	itemDescAssetID := itemFields[5].Descriptor()
	// item.DefaultAssetID holds the default value on creation for the asset_id field.
	item.DefaultAssetID = itemDescAssetID.Default.(int)
	// itemDescSerialNumber is the schema descriptor for serial_number field.
	itemDescSerialNumber := itemFields[6].Descriptor()
	// item.SerialNumberValidator is a validator for the "serial_number" field. It is called by the builders before save.
	item.SerialNumberValidator = itemDescSerialNumber.Validators[0].(func(string) error)
	// itemDescModelNumber is the schema descriptor for model_number field.
	itemDescModelNumber := itemFields[7].Descriptor()
	// item.ModelNumberValidator is a validator for the "model_number" field. It is called by the builders before save.
	item.ModelNumberValidator = itemDescModelNumber.Validators[0].(func(string) error)
	// itemDescManufacturer is the schema descriptor for manufacturer field.
	itemDescManufacturer := itemFields[8].Descriptor()
	// item.ManufacturerValidator is a validator for the "manufacturer" field. It is called by the builders before save.
	item.ManufacturerValidator = itemDescManufacturer.Validators[0].(func(string) error)
	// itemDescLifetimeWarranty is the schema descriptor for lifetime_warranty field.
	itemDescLifetimeWarranty := itemFields[9].Descriptor()
	// item.DefaultLifetimeWarranty holds the default value on creation for the lifetime_warranty field.
	item.DefaultLifetimeWarranty = itemDescLifetimeWarranty.Default.(bool)
	// itemDescWarrantyDetails is the schema descriptor for warranty_details field.
	itemDescWarrantyDetails := itemFields[11].Descriptor()
	// item.WarrantyDetailsValidator is a validator for the "warranty_details" field. It is called by the builders before save.
	item.WarrantyDetailsValidator = itemDescWarrantyDetails.Validators[0].(func(string) error)
	// itemDescPurchasePrice is the schema descriptor for purchase_price field.
	itemDescPurchasePrice := itemFields[14].Descriptor()
	// item.DefaultPurchasePrice holds the default value on creation for the purchase_price field.
	item.DefaultPurchasePrice = itemDescPurchasePrice.Default.(float64)
	// itemDescSetPrice is the schema descriptor for set_price field.
	itemDescSetPrice := itemFields[15].Descriptor()
	// item.DefaultSetPrice holds the default value on creation for the set_price field.
	item.DefaultSetPrice = itemDescSetPrice.Default.(float64)
	// itemDescOriginalPrice is the schema descriptor for original_price field.
	itemDescOriginalPrice := itemFields[16].Descriptor()
	// item.DefaultOriginalPrice holds the default value on creation for the original_price field.
	item.DefaultOriginalPrice = itemDescOriginalPrice.Default.(float64)
	// itemDescSoldPrice is the schema descriptor for sold_price field.
	itemDescSoldPrice := itemFields[19].Descriptor()
	// item.DefaultSoldPrice holds the default value on creation for the sold_price field.
	item.DefaultSoldPrice = itemDescSoldPrice.Default.(float64)
	// itemDescSoldNotes is the schema descriptor for sold_notes field.
	itemDescSoldNotes := itemFields[20].Descriptor()
	// item.SoldNotesValidator is a validator for the "sold_notes" field. It is called by the builders before save.
	item.SoldNotesValidator = itemDescSoldNotes.Validators[0].(func(string) error)
	// itemDescID is the schema descriptor for id field.
	itemDescID := itemMixinFields0[0].Descriptor()
	// item.DefaultID holds the default value on creation for the id field.
	item.DefaultID = itemDescID.Default.(func() uuid.UUID)
	itemfieldMixin := schema.ItemField{}.Mixin()
	itemfieldMixinFields0 := itemfieldMixin[0].Fields()
	_ = itemfieldMixinFields0
	itemfieldMixinFields1 := itemfieldMixin[1].Fields()
	_ = itemfieldMixinFields1
	itemfieldFields := schema.ItemField{}.Fields()
	_ = itemfieldFields
	// itemfieldDescCreatedAt is the schema descriptor for created_at field.
	itemfieldDescCreatedAt := itemfieldMixinFields0[1].Descriptor()
	// itemfield.DefaultCreatedAt holds the default value on creation for the created_at field.
	itemfield.DefaultCreatedAt = itemfieldDescCreatedAt.Default.(func() time.Time)
	// itemfieldDescUpdatedAt is the schema descriptor for updated_at field.
	itemfieldDescUpdatedAt := itemfieldMixinFields0[2].Descriptor()
	// itemfield.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	itemfield.DefaultUpdatedAt = itemfieldDescUpdatedAt.Default.(func() time.Time)
	// itemfield.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	itemfield.UpdateDefaultUpdatedAt = itemfieldDescUpdatedAt.UpdateDefault.(func() time.Time)
	// itemfieldDescName is the schema descriptor for name field.
	itemfieldDescName := itemfieldMixinFields1[0].Descriptor()
	// itemfield.NameValidator is a validator for the "name" field. It is called by the builders before save.
	itemfield.NameValidator = func() func(string) error {
		validators := itemfieldDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// itemfieldDescDescription is the schema descriptor for description field.
	itemfieldDescDescription := itemfieldMixinFields1[1].Descriptor()
	// itemfield.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	itemfield.DescriptionValidator = itemfieldDescDescription.Validators[0].(func(string) error)
	// itemfieldDescTextValue is the schema descriptor for text_value field.
	itemfieldDescTextValue := itemfieldFields[1].Descriptor()
	// itemfield.TextValueValidator is a validator for the "text_value" field. It is called by the builders before save.
	itemfield.TextValueValidator = itemfieldDescTextValue.Validators[0].(func(string) error)
	// itemfieldDescBooleanValue is the schema descriptor for boolean_value field.
	itemfieldDescBooleanValue := itemfieldFields[3].Descriptor()
	// itemfield.DefaultBooleanValue holds the default value on creation for the boolean_value field.
	itemfield.DefaultBooleanValue = itemfieldDescBooleanValue.Default.(bool)
	// itemfieldDescTimeValue is the schema descriptor for time_value field.
	itemfieldDescTimeValue := itemfieldFields[4].Descriptor()
	// itemfield.DefaultTimeValue holds the default value on creation for the time_value field.
	itemfield.DefaultTimeValue = itemfieldDescTimeValue.Default.(func() time.Time)
	// itemfieldDescID is the schema descriptor for id field.
	itemfieldDescID := itemfieldMixinFields0[0].Descriptor()
	// itemfield.DefaultID holds the default value on creation for the id field.
	itemfield.DefaultID = itemfieldDescID.Default.(func() uuid.UUID)
	labelMixin := schema.Label{}.Mixin()
	labelMixinFields0 := labelMixin[0].Fields()
	_ = labelMixinFields0
	labelMixinFields1 := labelMixin[1].Fields()
	_ = labelMixinFields1
	labelFields := schema.Label{}.Fields()
	_ = labelFields
	// labelDescCreatedAt is the schema descriptor for created_at field.
	labelDescCreatedAt := labelMixinFields0[1].Descriptor()
	// label.DefaultCreatedAt holds the default value on creation for the created_at field.
	label.DefaultCreatedAt = labelDescCreatedAt.Default.(func() time.Time)
	// labelDescUpdatedAt is the schema descriptor for updated_at field.
	labelDescUpdatedAt := labelMixinFields0[2].Descriptor()
	// label.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	label.DefaultUpdatedAt = labelDescUpdatedAt.Default.(func() time.Time)
	// label.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	label.UpdateDefaultUpdatedAt = labelDescUpdatedAt.UpdateDefault.(func() time.Time)
	// labelDescName is the schema descriptor for name field.
	labelDescName := labelMixinFields1[0].Descriptor()
	// label.NameValidator is a validator for the "name" field. It is called by the builders before save.
	label.NameValidator = func() func(string) error {
		validators := labelDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// labelDescDescription is the schema descriptor for description field.
	labelDescDescription := labelMixinFields1[1].Descriptor()
	// label.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	label.DescriptionValidator = labelDescDescription.Validators[0].(func(string) error)
	// labelDescColor is the schema descriptor for color field.
	labelDescColor := labelFields[0].Descriptor()
	// label.ColorValidator is a validator for the "color" field. It is called by the builders before save.
	label.ColorValidator = labelDescColor.Validators[0].(func(string) error)
	// labelDescID is the schema descriptor for id field.
	labelDescID := labelMixinFields0[0].Descriptor()
	// label.DefaultID holds the default value on creation for the id field.
	label.DefaultID = labelDescID.Default.(func() uuid.UUID)
	locationMixin := schema.Location{}.Mixin()
	locationMixinFields0 := locationMixin[0].Fields()
	_ = locationMixinFields0
	locationMixinFields1 := locationMixin[1].Fields()
	_ = locationMixinFields1
	locationFields := schema.Location{}.Fields()
	_ = locationFields
	// locationDescCreatedAt is the schema descriptor for created_at field.
	locationDescCreatedAt := locationMixinFields0[1].Descriptor()
	// location.DefaultCreatedAt holds the default value on creation for the created_at field.
	location.DefaultCreatedAt = locationDescCreatedAt.Default.(func() time.Time)
	// locationDescUpdatedAt is the schema descriptor for updated_at field.
	locationDescUpdatedAt := locationMixinFields0[2].Descriptor()
	// location.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	location.DefaultUpdatedAt = locationDescUpdatedAt.Default.(func() time.Time)
	// location.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	location.UpdateDefaultUpdatedAt = locationDescUpdatedAt.UpdateDefault.(func() time.Time)
	// locationDescName is the schema descriptor for name field.
	locationDescName := locationMixinFields1[0].Descriptor()
	// location.NameValidator is a validator for the "name" field. It is called by the builders before save.
	location.NameValidator = func() func(string) error {
		validators := locationDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// locationDescDescription is the schema descriptor for description field.
	locationDescDescription := locationMixinFields1[1].Descriptor()
	// location.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	location.DescriptionValidator = locationDescDescription.Validators[0].(func(string) error)
	// locationDescID is the schema descriptor for id field.
	locationDescID := locationMixinFields0[0].Descriptor()
	// location.DefaultID holds the default value on creation for the id field.
	location.DefaultID = locationDescID.Default.(func() uuid.UUID)
	maintenanceentryMixin := schema.MaintenanceEntry{}.Mixin()
	maintenanceentryMixinFields0 := maintenanceentryMixin[0].Fields()
	_ = maintenanceentryMixinFields0
	maintenanceentryFields := schema.MaintenanceEntry{}.Fields()
	_ = maintenanceentryFields
	// maintenanceentryDescCreatedAt is the schema descriptor for created_at field.
	maintenanceentryDescCreatedAt := maintenanceentryMixinFields0[1].Descriptor()
	// maintenanceentry.DefaultCreatedAt holds the default value on creation for the created_at field.
	maintenanceentry.DefaultCreatedAt = maintenanceentryDescCreatedAt.Default.(func() time.Time)
	// maintenanceentryDescUpdatedAt is the schema descriptor for updated_at field.
	maintenanceentryDescUpdatedAt := maintenanceentryMixinFields0[2].Descriptor()
	// maintenanceentry.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	maintenanceentry.DefaultUpdatedAt = maintenanceentryDescUpdatedAt.Default.(func() time.Time)
	// maintenanceentry.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	maintenanceentry.UpdateDefaultUpdatedAt = maintenanceentryDescUpdatedAt.UpdateDefault.(func() time.Time)
	// maintenanceentryDescName is the schema descriptor for name field.
	maintenanceentryDescName := maintenanceentryFields[3].Descriptor()
	// maintenanceentry.NameValidator is a validator for the "name" field. It is called by the builders before save.
	maintenanceentry.NameValidator = func() func(string) error {
		validators := maintenanceentryDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// maintenanceentryDescDescription is the schema descriptor for description field.
	maintenanceentryDescDescription := maintenanceentryFields[4].Descriptor()
	// maintenanceentry.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	maintenanceentry.DescriptionValidator = maintenanceentryDescDescription.Validators[0].(func(string) error)
	// maintenanceentryDescCost is the schema descriptor for cost field.
	maintenanceentryDescCost := maintenanceentryFields[5].Descriptor()
	// maintenanceentry.DefaultCost holds the default value on creation for the cost field.
	maintenanceentry.DefaultCost = maintenanceentryDescCost.Default.(float64)
	// maintenanceentryDescID is the schema descriptor for id field.
	maintenanceentryDescID := maintenanceentryMixinFields0[0].Descriptor()
	// maintenanceentry.DefaultID holds the default value on creation for the id field.
	maintenanceentry.DefaultID = maintenanceentryDescID.Default.(func() uuid.UUID)
	notifierMixin := schema.Notifier{}.Mixin()
	notifierMixinFields0 := notifierMixin[0].Fields()
	_ = notifierMixinFields0
	notifierFields := schema.Notifier{}.Fields()
	_ = notifierFields
	// notifierDescCreatedAt is the schema descriptor for created_at field.
	notifierDescCreatedAt := notifierMixinFields0[1].Descriptor()
	// notifier.DefaultCreatedAt holds the default value on creation for the created_at field.
	notifier.DefaultCreatedAt = notifierDescCreatedAt.Default.(func() time.Time)
	// notifierDescUpdatedAt is the schema descriptor for updated_at field.
	notifierDescUpdatedAt := notifierMixinFields0[2].Descriptor()
	// notifier.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	notifier.DefaultUpdatedAt = notifierDescUpdatedAt.Default.(func() time.Time)
	// notifier.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	notifier.UpdateDefaultUpdatedAt = notifierDescUpdatedAt.UpdateDefault.(func() time.Time)
	// notifierDescName is the schema descriptor for name field.
	notifierDescName := notifierFields[0].Descriptor()
	// notifier.NameValidator is a validator for the "name" field. It is called by the builders before save.
	notifier.NameValidator = func() func(string) error {
		validators := notifierDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// notifierDescURL is the schema descriptor for url field.
	notifierDescURL := notifierFields[1].Descriptor()
	// notifier.URLValidator is a validator for the "url" field. It is called by the builders before save.
	notifier.URLValidator = func() func(string) error {
		validators := notifierDescURL.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(url string) error {
			for _, fn := range fns {
				if err := fn(url); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// notifierDescIsActive is the schema descriptor for is_active field.
	notifierDescIsActive := notifierFields[2].Descriptor()
	// notifier.DefaultIsActive holds the default value on creation for the is_active field.
	notifier.DefaultIsActive = notifierDescIsActive.Default.(bool)
	// notifierDescID is the schema descriptor for id field.
	notifierDescID := notifierMixinFields0[0].Descriptor()
	// notifier.DefaultID holds the default value on creation for the id field.
	notifier.DefaultID = notifierDescID.Default.(func() uuid.UUID)
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields0[1].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields0[2].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.NameValidator is a validator for the "name" field. It is called by the builders before save.
	user.NameValidator = func() func(string) error {
		validators := userDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[1].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = func() func(string) error {
		validators := userDescEmail.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(email string) error {
			for _, fn := range fns {
				if err := fn(email); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescPassword is the schema descriptor for password field.
	userDescPassword := userFields[2].Descriptor()
	// user.PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	user.PasswordValidator = func() func(string) error {
		validators := userDescPassword.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(password string) error {
			for _, fn := range fns {
				if err := fn(password); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescIsSuperuser is the schema descriptor for is_superuser field.
	userDescIsSuperuser := userFields[3].Descriptor()
	// user.DefaultIsSuperuser holds the default value on creation for the is_superuser field.
	user.DefaultIsSuperuser = userDescIsSuperuser.Default.(bool)
	// userDescSuperuser is the schema descriptor for superuser field.
	userDescSuperuser := userFields[4].Descriptor()
	// user.DefaultSuperuser holds the default value on creation for the superuser field.
	user.DefaultSuperuser = userDescSuperuser.Default.(bool)
	// userDescID is the schema descriptor for id field.
	userDescID := userMixinFields0[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
