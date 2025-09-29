package elements

import (
	"context"
	"database/sql"
	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

/**
 * 
 * ELEMENT STORE
 * 
 * 
 **/

const (
	CreateElementMySQLQuery = `
		INSERT INTO 
		elements (
			slug, level, active, locked, parent_id, definition_id, status_id, created_by_profile_id,
			main_entity, name, alternate_name, description, same_as, url
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	UpdateElementMetaMySQLQuery = `
		UPDATE elements SET active = ?, locked = ?, status_id = ? WHERE id = ?
	`
	
	UpdateElementDescriptiveMySQLQuery = `
		UPDATE elements SET alternate_name = ?, description = ?, same_as = ? WHERE element_id = ?
	`
	
	DeleteElementMySQLQuery = "UPDATE elements SET deleted_at = NOW() WHERE id = ?"

	CreatePersonSchemaMySQLQuery = `
		INSERT INTO schema_persons (birth_date, birth_place, nationality, given_name, alternate_name, additional_name, family_name, gender, prefix, suffix)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	UpdatePersonSchemaNameMySQLQuery = `
		UPDATE schema_persons SET given_name = ?, alternate_name = ?, additional_name = ?, family_name = ?, prefix = ?, suffix = ?
		WHERE element_id = ?
	`

	CreateContactPointSchemaMySQLQuery = `
		INSERT schema_contact_points (element_id, email, fax_number, phone, url)
		VALUES (?, ?, ?, ?, ?)
	`
	UpdateContactPointSchemaMySQLQuery = `
		UPDATE schema_contact_points SET 
		email = ?, fax_number = ?, phone = ?, url = ?
		WHERE element_id = ?
	`

	CreatePostalAddressSchemaMySQLQuery = `
		INSERT schema_postal_addresses (element_id, country, locality, region, postal_code, street, number, address)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	UpdatePostalAddressSchemaMySQLQuery = `
		UPDATE schema_postal_addresses SET 
		country = ?, locality = ?, region = ?, postal_code = ?, street = ?, number = ?, address = ?
		WHERE element_id = ?
	`
	
	CreateWebResourceSchemaMySQLQuery = `
		INSERT schema_webresources (element_id, name, scope, body)
		VALUES (?, ?, ?, ?)
	`
	UpdateWebResourceSchemaMySQLQuery = `
		UPDATE schema_webresources SET
		name = ?, scope = ?, body = ? 
		WHERE element_id = ?
	`

	CreateTextSchemaMySQLQuery = `
		INSERT schema_text (element_id, headline, description, body)
		VALUES (?, ?, ?, ?)
	`
	UpdateTextSchemaMySQLQuery = `
		UPDATE schema_text SET
		headline = ?, description = ?, body = ? 
		WHERE element_id = ?
	`

	CreateArticleSchemaMySQLQuery = `
		INSERT schema_article (
			about, alternate_name, alternative_headline, article_body, article_section,
			author, genre, headline, text			
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	UpdateArticleSchemaMySQLQuery = `
		UPDATE schema_article SET 
			about = ?, alternate_name = ?, alternative_headline = ?, article_body = ?,
			author = ?, genre = ?, headline = ?, text = ?
		WHERE element_id = ?
	`

	CreateTimeTrackingSchemaMySQLQuery = `
		INSERT schema_time_tracking (element_id, doortime, start_date, previous_start_date, end_date)
		VALUES (?, ?, ?, ?, ?)
	`
	UpdateTimeTrackingSchemaMySQLQuery = `
		UPDATE schema_time_tracking SET
			doortime = ?,
			previous_start_date = ?,
			start_date = ?,
			end_date = ? 
		WHERE element_id = ?
	`

	CreateFileSchemaMySQLQuery = `
		INSERT INTO schema_files (element_id, width, height, filename, original_filename, filepath)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	UpdateFileSchemaMySQLQuery = `
		UPDATE schema_files SET filename = ?, filepath = ? WHERE id = ?
	`
)

type ElementStore struct {
	adapter *sql.DB
}

func NewElementStore(adapter *sql.DB) *ElementStore {
	return &ElementStore{adapter: adapter}
}

func (s *ElementStore) Create(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	res, err := s.adapter.ExecContext(ctx,
		CreateElementMySQLQuery,
		element.Slug, element.Level, element.Active, element.Locked, element.ParentID, element.DefinitionID, element.StatusID, element.CreatedByProfileID,
		element.MainEntity, element.Name, element.AlternateName, element.Description, element.SameAs, element.Url,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	element.ID = uint(id)
	return nil
}

func (s *ElementStore) CreateFileSchema(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	_, err := s.adapter.ExecContext(ctx,
		CreateFileSchemaMySQLQuery,
		element.ID,
		element.File.Filename, element.File.OriginalFilename, element.File.Filepath,
		element.File.Width, element.File.Height,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *ElementStore) CreateTextSchema(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	_, err := s.adapter.ExecContext(ctx,
		CreateTextSchemaMySQLQuery,
		element.ID,
		element.Text.Headline, element.Text.Description, element.Text.Body,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *ElementStore) CreatePersonSchema(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	_, err := s.adapter.ExecContext(ctx,
		CreatePersonSchemaMySQLQuery,
		element.Person.Prefix, element.Person.Suffix, element.Person.Gender,
		element.Person.GivenName, element.Person.AlternateName, element.Person.AdditionalName, element.Person.FamilyName,
		element.Person.BirthDate, element.Person.BirthPlace, element.Person.Nationality,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *ElementStore) CreateContactPointSchema(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	_, err := s.adapter.ExecContext(ctx,
		CreateContactPointSchemaMySQLQuery,
		element.ID,
		element.ContactPoint.Email, element.ContactPoint.FaxNumber, element.ContactPoint.Phone, element.ContactPoint.Url,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *ElementStore) CreatePostalAddressSchema(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	_, err := s.adapter.ExecContext(ctx,
		CreatePostalAddressSchemaMySQLQuery,
		element.ID,
		element.PostalAddress.Country, element.PostalAddress.Locality, element.PostalAddress.Region, element.PostalAddress.PostalCode,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *ElementStore) CreateTimeTrackingSchema(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	_, err := s.adapter.ExecContext(ctx,
		CreateTimeTrackingSchemaMySQLQuery,
		element.ID,
		element.TimeTracking.DoorTime, element.TimeTracking.PreviousStartDate, element.TimeTracking.EndDate,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *ElementStore) CreateWebResourceSchema(ctx context.Context, p actions.Payload) error {
	element := p.(*Element)
	_, err := s.adapter.ExecContext(ctx,
		CreateWebResourceSchemaMySQLQuery,
		element.ID,
		element.WebResource.Name, element.WebResource.Scope, element.WebResource.Body, 
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *ElementStore) Delete(ctx context.Context, p actions.Payload) error {
	payload := p.(*actions.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteElementMySQLQuery,
		payload.ID, // element.ID
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}
