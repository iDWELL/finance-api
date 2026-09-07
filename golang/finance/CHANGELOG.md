# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.10.0] - 2026-09-07

### Added

- `RealEstateObject.UnsentFields()` — the optional fields that were not present in the JSON object the value was unmarshalled from, so that the receiver can keep their current value instead of clearing them. Returns nil for a value that was assembled in Go, whose fields are all authoritative

### Changed

- The optional `RealEstateObject` fields (`AccountingArea`, `UserAccount`, `Description`, `AlternativeAddresses`, `ZipCode`, `City`, `BankAccounts`, `ManagementEnd`) are omitted from the marshalled JSON while they are at their zero value, instead of being sent as `null`. A sender that does not own a field now leaves it out of the payload, where sending `null` used to clear the value written by the sender that does own it

## [1.9.0] - 2026-08-27

### Added

- `RealEstateObject.ManagementEnd` (`date.NullableDate`) — mirrors iX-Haus `OBJ.DWVERWALTUNGSENDE`; independent of `Active`

## [1.8.0] - 2026-08-04

### Added

- `ObjectTenantOwnerRole` enum with values `UNSPECIFIED`, `TENANT` and `OWNER`
- `ObjectTenantOwner.Role` — the role of the person on the unit. A tenant/owner record identifies one person, on one unit, in one role, because the same person can be the tenant of one unit and the owner of another, or both on the same unit. Sources that list tenants and owners together without telling them apart leave the field empty, which `Validate()` normalizes to `UNSPECIFIED`, so payloads written before the field existed stay valid

## [1.7.0] - 2026-07-24

### Added

- `RevokeUserAccess()` — revokes portal access from users by email via `POST /user/access/revoke`; the user records themselves are kept
- `UserAccessRevokeRequest` / `UserAccessRevokeResult` structs

## [1.6.0] - 2026-07-15

### Added

- `UpdateClientCompanyCostCenter()` — single-record GraphQL mutation to update an existing `ClientCompanyCostCenter`'s number/description; currency is optional and left unchanged when omitted
- `RealEstateObjectType` enum values `FIBU` (Finanzbuchhaltung / financial accounting object) and `VWO` (Verwaltungsobjekt / management-administration object)

## [1.5.0] - 2026-07-10

### Added

- `Partner.Paragraph13bHandling` (`PartnerParagraph13bHandling` string enum: `""`, `MUST_NOT`, `MAY`, `MUST`) — mirrors iX-Haus `N13BOPTION`
- `NeedUpsertTaxExemption()` now also returns true when this field is set

## [1.4.1] - 2026-07-09

### Fixed

- `CreateClientCompanyCostCenter()` — GraphQL variables `$number`/`$description` declared as `String!` while the API schema expects `NonEmptyText!`, so every call was rejected with a 400 before execution

## [1.4.0] - 2026-07-08

### Added

- `CreateClientCompanyCostCenter()` — single-record GraphQL mutation to create a `ClientCompanyCostCenter` (currency defaults to EUR)
- `ClientCompanyCostCenter` struct

## [1.3.2] - 2026-06-05

### Fixed

- `CreateRealEstateObjectCostCenter()` — corrected GraphQL field names: `objectInstanceId` → `objectInstanceRowId`, `clientCompanyCostCenterId` → `clientCompanyCostCenterRowId`

## [1.3.1] - 2026-06-01

_Re-tag of [1.3.0] — no code changes._

## [1.3.0] - 2026-05-31

### Added

- `GLAccount.CostCenterHandling` (`GLAccountCostCenterHandling` string enum: `""`, `MUST_NOT`, `MUST`, `MAY`) — mirrors iX-Haus `Sachkonto.KSTSTBEHANDLUNG`
- `GLAccount.ProjectHandling` (`GLAccountProjectHandling`) — mirrors iX-Haus `Sachkonto.PRJBEHANDLUNG`
- `GLAccount.OrderHandling` (`GLAccountOrderHandling`) — mirrors iX-Haus `Sachkonto.AuftragBehandlung`
- `GLAccount.OrderHandlingAmountThreshold` (`*money.Amount`) — mirrors iX-Haus `Sachkonto.AuftragBehandlungBetragGrenze`

## [1.2.0] - 2026-05-26

### Added

- `CreateRealEstateObjectCostCenter()` — single-record GraphQL mutation to associate a `ClientCompanyCostCenter` with a `RealEstateObject`
- `RealEstateObjectCostCenter` struct

## [1.1.1] - 2026-05-06

### Changed

- `GLAccount.DefaultVATCode` type changed from `nullable.Type[int]` to `nullable.Type[int64]`

## [1.1.0] - 2026-05-06

### Added

- `GLAccount.DefaultVATCode` (`nullable.Type[int]`) — optional default VAT code for a general ledger account

## [1.0.1] - 2026-05-04

### Removed

- `Partner.Paragraph13bApplicable` (`nullable.Type[bool]`) — unused field removed from `Partner` struct and from `NeedUpsertTaxExemption()` check

## [1.0.0] - 2026-05-01

### Added

- `WithBaseURL()`, `WithHTTPClient()` — context helpers for configuring API base URL and HTTP client
- `GetCurrentCompany()` — fetch current company data (name, legal form, address, VAT) via GraphQL
- `PostPartners()` — bulk upsert business partners with validation, normalization, and flexible error handling; `Partner.Validate()`, `Partner.Normalize()`
- `PostBankAccounts()` — bulk upsert bank accounts; `BankAccount.Validate()`, `BankAccount.Normalize()`
- `PostGLAccounts()` — bulk upsert general ledger accounts; `GLAccount.Validate()`
- `PostRealEstateObjects()` — bulk upsert real estate objects; `RealEstateObject.Validate()`; `RealEstateObjectType` enum (`WEG`, `HI`, `SUB`, `KREIS`, `MANDANT`, `MRG`, `MHV`, `SEV`, `HBH`) with `IsVirtual()`
- `PostObjectGroups()` — bulk upsert object groups linking real estate objects
- `PostObjectRoles()` — bulk upsert user–object role assignments
- `PostObjectTenantOwners()` — bulk upsert tenant/owner records per real estate object; `ObjectTenantOwner.Validate()`
- `PostObjectInstancesWithIDProp()` — generic bulk upsert for arbitrary object instances with a custom ID property
- `Invoice` with `Validate()` — invoice model covering partner, dates, amounts, VAT, currency, delivery and payment details; `AccountingItem` for individual booking lines
- `UploadDocument()` — upload PDF, PNG, JPEG, or TIFF with optional invoice data and tags; returns document ID
- `DownloadDocumentPDF()` — download document as PDF; options: `WithAuditTrail()`, `WithAuditTrailLang()`, `WithEmbedXML()`
- `ImportState` enum (`UNCHANGED`, `UPDATED`, `CREATED`, `ERROR`) returned by all `Post*` import results
