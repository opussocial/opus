package elements

import (
	"time"
)

// Person see : https://schema.org/Person
type PersonSchema struct {

	ID                  uint `json:"id"`
	ElementID			uint `json:"elementId"`

	// BirthDate see : https://schema.org/birthDate
	// Date of birth.
	// types : Date
	BirthDate time.Time `json:"birthDate,omitempty"`

	// BirthPlace see : https://schema.org/birthPlace
	// The place where the person was born.
	// types : Place
	BirthPlace string `json:"birthPlace,omitempty"`

	// AlternateName see : https://schema.org/alternateName
	// An alias for the item.
	// types : Text
	AlternateName string `json:"alternateName,omitempty"`

	// GivenName see : https://schema.org/givenName
	// Given name. In the U.S., the first name of a Person. This can be used along with familyName instead of the name property.
	// types : Text
	GivenName string `json:"givenName,omitempty"`

	// AdditionalName see : https://schema.org/additionalName
	// An additional name for a Person, can be used for a middle name.
	// types : Text
	AdditionalName string `json:"additionalName,omitempty"`

	// FamilyName see : https://schema.org/familyName
	// Family name. In the U.S., the last name of an Person. This can be used along with givenName instead of the name property.
	// types : Text
	FamilyName string `json:"familyName,omitempty"`

	// Gender see : https://schema.org/gender
	// Gender of the person. While http://schema.org/Male and http://schema.org/Female may be used, text strings are also acceptable for people who do not identify as a binary gender.
	// types : GenderType Text
	Gender string `json:"gender,omitempty"`

	// HonorificPrefix see : https://schema.org/honorificPrefix
	// An honorific prefix preceding a Person&#39;s name such as Dr/Mrs/Mr.
	// types : Text
	Prefix string `json:"prefix,omitempty"`

	// HonorificSuffix see : https://schema.org/honorificSuffix
	// An honorific suffix preceding a Person&#39;s name such as M.D. /PhD/MSCSW.
	// types : Text
	Suffix string `json:"suffix,omitempty"`

	// Nationality see : https://schema.org/nationality
	// Nationality of the person.
	// types : Country
	Nationality string `json:"nationality,omitempty"`
}

func (p *PersonSchema) Validate() error {
	return nil
}

func (p *PersonSchema) Process() error {
	return nil
}

// ContactPoint see : https://schema.org/ContactPoint
type ContactPointSchema struct {

	ID                  uint
	ElementID			uint

	// Email see : https://schema.org/email
	// Email address.
	// types : Text
	Email string `json:"email,omitempty"`

	// FaxNumber see : https://schema.org/faxNumber
	// The fax number.
	// types : Text
	FaxNumber string `json:"faxNumber,omitempty"`

	// Telephone see : https://schema.org/telephone
	// The telephone number.
	// types : Text
	Phone string `json:"phone,omitempty"`

	// Url see : https://schema.org/url
	// URL of the item.
	// types : URL
	Url string `json:"url,omitempty"`
}

func (p *ContactPointSchema) Validate() error {
	return nil
}

func (p *ContactPointSchema) Process() error {
	return nil
}

type WebResourceSchema struct {

	ID                  uint
	ElementID			uint

	// Name see : https://schema.org/name
	// The name of the item.Web
	// types : Text
	Name string `json:"name,omitempty"`

	// link or embed
	// types : Text
	Scope string `json:"scope,omitempty"`

	// types : Text
	Body string `json:"body,omitempty"`
}

func (p *WebResourceSchema) Validate() error {
	return nil
}

func (p *WebResourceSchema) Process() error {
	return nil
}

type TextSchema struct {

	ID                  uint
	ElementID			uint

	// Headline see : https://schema.org/headline
	// Headline of the article.
	// types : Text
	Headline string `json:"headline,omitempty"`

	// Description see : https://schema.org/description
	// A description of the item.
	// types : Text
	Description string `json:"description,omitempty"`

	// ArticleBody see : https://schema.org/articleBody
	// The actual body of the article.
	// types : Text
	Body string `json:"body,omitempty"`
}

func (p *TextSchema) Validate() error {
	return nil
}

func (p *TextSchema) Process() error {
	return nil
}

// Article see : https://schema.org/Article
type ArticleSchema struct {

	ID                  uint
	ElementID			uint

	// About see : https://schema.org/about
	// The subject matter of the content. Inverse property: subjectOf (see: https://schema.orghttps://pending.schema.org/subjectOf).
	// types : Thing
	About string `json:"about,omitempty"`

	// AlternateName see : https://schema.org/alternateName
	// An alias for the item.
	// types : Text
	AlternateName string `json:"alternateName,omitempty"`

	// AlternativeHeadline see : https://schema.org/alternativeHeadline
	// A secondary title of the CreativeWork.
	// types : Text
	AlternativeHeadline string `json:"alternativeHeadline,omitempty"`

	// ArticleBody see : https://schema.org/articleBody
	// The actual body of the article.
	// types : Text
	ArticleBody string `json:"articleBody,omitempty"`

	// ArticleSection see : https://schema.org/articleSection
	// Articles may belong to one or more &#39;sections&#39; in a magazine or newspaper, such as Sports, Lifestyle, etc.
	// types : Text
	ArticleSection string `json:"articleSection,omitempty"`

	// Author see : https://schema.org/author
	// The author of this content or rating. Please note that author is special in that HTML 5 provides a special mechanism for indicating authorship via the rel tag. That is equivalent to this and may be used interchangeably.
	// types : Organization Person
	Author string `json:"author,omitempty"`

	// Genre see : https://schema.org/genre
	// Genre of the creative work, broadcast channel or group.
	// types : Text URL
	Genre string `json:"genre,omitempty"`

	// Headline see : https://schema.org/headline
	// Headline of the article.
	// types : Text
	Headline string `json:"headline,omitempty"`

	// Text see : https://schema.org/text
	// The textual content of this CreativeWork.
	// types : Text
	Text string `json:"text,omitempty"`
}

func (p *ArticleSchema) Validate() error {
	return nil
}

func (p *ArticleSchema) Process() error {
	return nil
}

type PostalAddressSchema struct {

	ID                  uint
	ElementID			uint

	// AddressCountry see : https://schema.org/addressCountry
	// The country. For example, USA. You can also provide the two-letter ISO 3166-1 alpha-2 country code (see: https://schema.orghttp://en.wikipedia.org/wiki/ISO_3166-1).
	// types : Country Text
	Country string `json:"country,omitempty"`

	// AddressLocality see : https://schema.org/addressLocality
	// The locality. For example, Mountain View.
	// types : Text
	Locality string `json:"locality,omitempty"`

	// AddressRegion see : https://schema.org/addressRegion
	// The region. For example, CA.
	// types : Text
	Region string `json:"region,omitempty"`

	// PostalCode see : https://schema.org/postalCode
	// The postal code. For example, 94043.
	// types : Text
	PostalCode string `json:"postalCode,omitempty"`

	// StreetAddress see : https://schema.org/streetAddress
	// The street address. For example, 1600 Amphitheatre Pkwy.
	// types : Text
	Street string `json:"street,omitempty"`
	Number string `json:"number,omitempty"`
	Address string `json:"address,omitempty"`
}

func (p *PostalAddressSchema) Validate() error {
	return nil
}

func (p *PostalAddressSchema) Process() error {
	return nil
}

// Event see : https://schema.org/Event
type TimeTrackingSchema struct {

	ID                  uint
	ElementID			uint

	// DoorTime see : https://schema.org/doorTime
	// The time admission will commence.
	// types : time.Time
	DoorTime time.Time `json:"doorTime,omitempty"`

	// StartDate see : https://schema.org/startDate
	// The start date and time of the item (in ISO 8601 date format (see: https://schema.orghttp://en.wikipedia.org/wiki/ISO_8601)).
	// types : Date time.Time
	StartDate time.Time `json:"startDate,omitempty"`

	// PreviousStartDate see : https://schema.org/previousStartDate
	// Used in conjunction with eventStatus for rescheduled or cancelled events. This property contains the previously scheduled start date. For rescheduled events, the startDate property should be used for the newly scheduled start date. In the (rare) case of an event that has been postponed and rescheduled multiple times, this field may be repeated.
	// types : Date
	PreviousStartDate time.Time `json:"previousStartDate,omitempty"`

	// EndDate see : https://schema.org/endDate
	// The end date and time of the item (in ISO 8601 date format (see: https://schema.orghttp://en.wikipedia.org/wiki/ISO_8601)).
	// types : Date time.Time
	EndDate time.Time `json:"endDate,omitempty"`
}

func (p *TimeTrackingSchema) Validate() error {
	return nil
}

func (p *TimeTrackingSchema) Process() error {
	return nil
}

type FileSchema struct {

	ID                  uint
	ElementID			uint

	Width            uint   `json:"width,omitempty"`
	Height           uint   `json:"height,omitempty"`
	Filename         string `json:"filename,omitempty"`
	OriginalFilename string `json:"originalFilename,omitempty"`
	Filepath         string `json:"filepath,omitempty"`
}

func (p *FileSchema) Validate() error {
	return nil
}

func (p *FileSchema) Process() error {
	return nil
}
