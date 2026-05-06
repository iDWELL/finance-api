# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
