package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ContactImportExportService handles vCard import/export for contacts
type ContactImportExportService struct {
	db *sql.DB
}

// ContactGroup represents a contact group/category
type ContactGroup struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Color       string      `json:"color,omitempty"`
	UserID      uuid.UUID   `json:"user_id"`
	Contacts    []uuid.UUID `json:"contacts"`
	IsPublic    bool        `json:"is_public"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ContactCategory represents a predefined category
type ContactCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
	IsDefault   bool   `json:"is_default"`
}

// VCardContact represents a contact in vCard format
type VCardContact struct {
	UID          string            `json:"uid"`
	FullName     string            `json:"full_name"`
	FirstName    string            `json:"first_name,omitempty"`
	LastName     string            `json:"last_name,omitempty"`
	Email        []string          `json:"email"`
	Phone        []PhoneNumber     `json:"phone"`
	Address      []Address         `json:"address"`
	Organization string            `json:"organization,omitempty"`
	Title        string            `json:"title,omitempty"`
	Photo        *ContactPhoto     `json:"photo,omitempty"`
	Birthday     *time.Time        `json:"birthday,omitempty"`
	Notes        string            `json:"notes,omitempty"`
	Categories   []string          `json:"categories"`
	URLs         []string          `json:"urls"`
	SocialMedia  map[string]string `json:"social_media,omitempty"`
	CustomFields map[string]string `json:"custom_fields,omitempty"`
}

// PhoneNumber represents a phone number with type
type PhoneNumber struct {
	Number string `json:"number"`
	Type   string `json:"type"` // "home", "work", "mobile", "fax", etc.
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Type       string `json:"type"` // "home", "work", etc.
}

// ContactPhoto represents a contact photo
type ContactPhoto struct {
	Data        []byte `json:"data"`
	ContentType string `json:"content_type"`
	Filename    string `json:"filename,omitempty"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	TotalContacts  int      `json:"total_contacts"`
	ImportedCount  int      `json:"imported_count"`
	SkippedCount   int      `json:"skipped_count"`
	ErrorCount     int      `json:"error_count"`
	DuplicateCount int      `json:"duplicate_count"`
	NewGroups      []string `json:"new_groups"`
	Errors         []string `json:"errors"`
}

// ExportOptions represents options for contact export
type ExportOptions struct {
	Format         string      `json:"format"` // "vcard", "csv", "json"
	IncludePhotos  bool        `json:"include_photos"`
	GroupIDs       []uuid.UUID `json:"group_ids,omitempty"`   // Export specific groups
	ContactIDs     []uuid.UUID `json:"contact_ids,omitempty"` // Export specific contacts
	IncludePrivate bool        `json:"include_private"`
	MaxPhotoSize   int         `json:"max_photo_size,omitempty"` // Max photo size in bytes
}

// NewContactImportExportService creates a new contact import/export service
func NewContactImportExportService(db *sql.DB) *ContactImportExportService {
	return &ContactImportExportService{db: db}
}

// ImportVCard imports contacts from vCard data
func (cies *ContactImportExportService) ImportVCard(vCardData string, userID uuid.UUID, options map[string]interface{}) (*ImportResult, error) {
	// Parse vCard data
	vCards := cies.parseVCardData(vCardData)

	result := &ImportResult{
		TotalContacts: len(vCards),
		NewGroups:     []string{},
		Errors:        []string{},
	}

	groupMap := make(map[string]uuid.UUID)

	for _, vCard := range vCards {
		err := cies.importVCardContact(vCard, userID, result, groupMap, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to import contact %s: %v", vCard.FullName, err))
			result.ErrorCount++
		}
	}

	return result, nil
}

// ExportContacts exports contacts to various formats
func (cies *ContactImportExportService) ExportContacts(userID uuid.UUID, options ExportOptions) ([]byte, error) {
	// Get contacts to export
	contacts, err := cies.getContactsForExport(userID, options)
	if err != nil {
		return nil, err
	}

	switch options.Format {
	case "vcard":
		return cies.exportToVCard(contacts, options)
	case "csv":
		return cies.exportToCSV(contacts, options)
	case "json":
		return cies.exportToJSON(contacts, options)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// CreateContactGroup creates a new contact group
func (cies *ContactImportExportService) CreateContactGroup(group *ContactGroup) error {
	group.ID = uuid.New()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	query := `
		INSERT INTO contact_groups (id, name, description, color, user_id, contacts, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := cies.db.Exec(query,
		group.ID, group.Name, group.Description, group.Color, group.UserID,
		group.Contacts, group.IsPublic, group.CreatedAt, group.UpdatedAt,
	)

	return err
}

// GetContactGroups gets all contact groups for a user
func (cies *ContactImportExportService) GetContactGroups(userID uuid.UUID) ([]ContactGroup, error) {
	query := `
		SELECT id, name, description, color, user_id, contacts, is_public, created_at, updated_at
		FROM contact_groups
		WHERE user_id = $1 OR is_public = true
		ORDER BY name
	`

	rows, err := cies.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []ContactGroup
	for rows.Next() {
		var group ContactGroup
		err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.Color, &group.UserID,
			&group.Contacts, &group.IsPublic, &group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			continue
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// AddContactsToGroup adds contacts to a group
func (cies *ContactImportExportService) AddContactsToGroup(groupID uuid.UUID, contactIDs []uuid.UUID, userID uuid.UUID) error {
	// Verify group ownership
	var groupOwner uuid.UUID
	err := cies.db.QueryRow("SELECT user_id FROM contact_groups WHERE id = $1", groupID).Scan(&groupOwner)
	if err != nil {
		return err
	}

	if groupOwner != userID {
		return fmt.Errorf("user does not own this group")
	}

	// Get current contacts in group
	var currentContacts []uuid.UUID
	err = cies.db.QueryRow("SELECT contacts FROM contact_groups WHERE id = $1", groupID).Scan(&currentContacts)
	if err != nil {
		return err
	}

	// Add new contacts (avoid duplicates)
	contactMap := make(map[uuid.UUID]bool)
	for _, contactID := range currentContacts {
		contactMap[contactID] = true
	}

	for _, contactID := range contactIDs {
		contactMap[contactID] = true
	}

	// Convert back to slice
	var allContacts []uuid.UUID
	for contactID := range contactMap {
		allContacts = append(allContacts, contactID)
	}

	// Update group
	query := `UPDATE contact_groups SET contacts = $1, updated_at = $2 WHERE id = $3`
	_, err = cies.db.Exec(query, allContacts, time.Now(), groupID)

	return err
}

// RemoveContactsFromGroup removes contacts from a group
func (cies *ContactImportExportService) RemoveContactsFromGroup(groupID uuid.UUID, contactIDs []uuid.UUID, userID uuid.UUID) error {
	// Verify group ownership
	var groupOwner uuid.UUID
	err := cies.db.QueryRow("SELECT user_id FROM contact_groups WHERE id = $1", groupID).Scan(&groupOwner)
	if err != nil {
		return err
	}

	if groupOwner != userID {
		return fmt.Errorf("user does not own this group")
	}

	// Get current contacts in group
	var currentContacts []uuid.UUID
	err = cies.db.QueryRow("SELECT contacts FROM contact_groups WHERE id = $1", groupID).Scan(&currentContacts)
	if err != nil {
		return err
	}

	// Create map for removal
	removeMap := make(map[uuid.UUID]bool)
	for _, contactID := range contactIDs {
		removeMap[contactID] = true
	}

	// Filter out contacts to remove
	var remainingContacts []uuid.UUID
	for _, contactID := range currentContacts {
		if !removeMap[contactID] {
			remainingContacts = append(remainingContacts, contactID)
		}
	}

	// Update group
	query := `UPDATE contact_groups SET contacts = $1, updated_at = $2 WHERE id = $3`
	_, err = cies.db.Exec(query, remainingContacts, time.Now(), groupID)

	return err
}

// GetContactCategories returns predefined contact categories
func (cies *ContactImportExportService) GetContactCategories() []ContactCategory {
	return []ContactCategory{
		{ID: "family", Name: "Family", Description: "Family members and relatives", Icon: "👨‍👩‍👧‍👦", Color: "#FF6B6B", IsDefault: true},
		{ID: "friends", Name: "Friends", Description: "Personal friends and acquaintances", Icon: "👥", Color: "#4ECDC4", IsDefault: true},
		{ID: "work", Name: "Work", Description: "Colleagues and professional contacts", Icon: "💼", Color: "#45B7D1", IsDefault: true},
		{ID: "business", Name: "Business", Description: "Business partners and clients", Icon: "🤝", Color: "#96CEB4", IsDefault: true},
		{ID: "medical", Name: "Medical", Description: "Healthcare providers", Icon: "⚕️", Color: "#FFEAA7", IsDefault: false},
		{ID: "emergency", Name: "Emergency", Description: "Emergency contacts", Icon: "🚨", Color: "#D63031", IsDefault: false},
		{ID: "school", Name: "School", Description: "Teachers, classmates, alumni", Icon: "🎓", Color: "#A29BFE", IsDefault: false},
		{ID: "sports", Name: "Sports", Description: "Sports teams and coaches", Icon: "⚽", Color: "#FD79A8", IsDefault: false},
		{ID: "religious", Name: "Religious", Description: "Religious organization contacts", Icon: "⛪", Color: "#E17055", IsDefault: false},
		{ID: "political", Name: "Political", Description: "Political contacts and organizations", Icon: "🏛️", Color: "#636E72", IsDefault: false},
	}
}

// GetContactsByGroup gets contacts in a specific group
func (cies *ContactImportExportService) GetContactsByGroup(groupID uuid.UUID, userID uuid.UUID) ([]map[string]interface{}, error) {
	query := `
		SELECT c.id, c.first_name, c.last_name, c.email, c.phone, c.company, c.position,
		       c.groups, c.created_at
		FROM contacts c
		JOIN contact_groups cg ON cg.user_id = c.user_id
		WHERE cg.id = $1 AND (cg.user_id = $2 OR cg.is_public = true)
		ORDER BY c.first_name, c.last_name
	`

	rows, err := cies.db.Query(query, groupID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []map[string]interface{}
	for rows.Next() {
		var contact struct {
			ID        uuid.UUID
			FirstName *string
			LastName  *string
			Email     *string
			Phone     *string
			Company   *string
			Position  *string
			Groups    []uuid.UUID
			CreatedAt time.Time
		}

		err := rows.Scan(
			&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email,
			&contact.Phone, &contact.Company, &contact.Position, &contact.Groups, &contact.CreatedAt,
		)
		if err != nil {
			continue
		}

		contacts = append(contacts, map[string]interface{}{
			"id":         contact.ID,
			"first_name": contact.FirstName,
			"last_name":  contact.LastName,
			"email":      contact.Email,
			"phone":      contact.Phone,
			"company":    contact.Company,
			"position":   contact.Position,
			"groups":     contact.Groups,
			"created_at": contact.CreatedAt,
		})
	}

	return contacts, nil
}

// SearchContactsInGroups searches contacts within specific groups
func (cies *ContactImportExportService) SearchContactsInGroups(userID uuid.UUID, groupIDs []uuid.UUID, query string) ([]map[string]interface{}, error) {
	if len(groupIDs) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Build IN clause for group IDs
	placeholders := make([]string, len(groupIDs))
	args := []interface{}{userID, query}
	for i, groupID := range groupIDs {
		placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
		args = append(args, groupID)
	}

	searchQuery := fmt.Sprintf(`
		SELECT DISTINCT c.id, c.first_name, c.last_name, c.email, c.phone, c.company,
		       c.position, c.groups, c.created_at
		FROM contacts c
		JOIN contact_groups cg ON cg.user_id = c.user_id
		WHERE cg.user_id = $1
		AND cg.id IN (%s)
		AND (
			c.first_name ILIKE $2 OR
			c.last_name ILIKE $2 OR
			c.email ILIKE $2 OR
			c.company ILIKE $2
		)
		ORDER BY c.first_name, c.last_name
		LIMIT 50
	`, strings.Join(placeholders, ","))

	rows, err := cies.db.Query(searchQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []map[string]interface{}
	for rows.Next() {
		var contact struct {
			ID        uuid.UUID
			FirstName *string
			LastName  *string
			Email     *string
			Phone     *string
			Company   *string
			Position  *string
			Groups    []uuid.UUID
			CreatedAt time.Time
		}

		err := rows.Scan(
			&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email,
			&contact.Phone, &contact.Company, &contact.Position, &contact.Groups, &contact.CreatedAt,
		)
		if err != nil {
			continue
		}

		contacts = append(contacts, map[string]interface{}{
			"id":         contact.ID,
			"first_name": contact.FirstName,
			"last_name":  contact.LastName,
			"email":      contact.Email,
			"phone":      contact.Phone,
			"company":    contact.Company,
			"position":   contact.Position,
			"groups":     contact.Groups,
			"created_at": contact.CreatedAt,
		})
	}

	return contacts, nil
}

// Helper methods for vCard parsing and generation

func (cies *ContactImportExportService) parseVCardData(vCardData string) []VCardContact {
	// Split vCard data into individual cards
	vCardSeparator := regexp.MustCompile(`(?m)^BEGIN:VCARD`)
	locations := vCardSeparator.FindAllStringIndex(vCardData, -1)

	var vCards []VCardContact

	for i, loc := range locations {
		start := loc[0]
		end := len(vCardData)
		if i < len(locations)-1 {
			end = locations[i+1][0]
		}

		vCardText := vCardData[start:end]
		vCard := cies.parseSingleVCard(vCardText)
		if vCard != nil {
			vCards = append(vCards, *vCard)
		}
	}

	return vCards
}

func (cies *ContactImportExportService) parseSingleVCard(vCardText string) *VCardContact {
	lines := strings.Split(vCardText, "\n")
	vCard := &VCardContact{
		Email:        []string{},
		Phone:        []PhoneNumber{},
		Address:      []Address{},
		Categories:   []string{},
		URLs:         []string{},
		SocialMedia:  make(map[string]string),
		CustomFields: make(map[string]string),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "BEGIN:") || strings.HasPrefix(line, "END:") {
			continue
		}

		// Handle line continuations
		for strings.HasSuffix(line, "=") {
			// This is a simple implementation - real vCard parsing is more complex
			break
		}

		cies.parseVCardLine(line, vCard)
	}

	return vCard
}

func (cies *ContactImportExportService) parseVCardLine(line string, vCard *VCardContact) {
	// Simple vCard line parsing (real implementation would be more robust)
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	property := strings.ToUpper(parts[0])
	value := parts[1]

	switch {
	case property == "FN":
		vCard.FullName = value
	case property == "N":
		// Parse name components
		nameParts := strings.Split(value, ";")
		if len(nameParts) >= 2 {
			vCard.LastName = nameParts[0]
			vCard.FirstName = nameParts[1]
		}
	case strings.HasPrefix(property, "EMAIL"):
		vCard.Email = append(vCard.Email, value)
	case strings.HasPrefix(property, "TEL"):
		phoneType := cies.extractVCardType(property)
		vCard.Phone = append(vCard.Phone, PhoneNumber{
			Number: value,
			Type:   phoneType,
		})
	case property == "ORG":
		vCard.Organization = value
	case property == "TITLE":
		vCard.Title = value
	case property == "NOTE":
		vCard.Notes = value
	case property == "CATEGORIES":
		vCard.Categories = strings.Split(value, ",")
	case strings.HasPrefix(property, "URL"):
		vCard.URLs = append(vCard.URLs, value)
	case property == "UID":
		vCard.UID = value
	case property == "BDAY":
		// Parse birthday
		if t, err := time.Parse("20060102", value); err == nil {
			vCard.Birthday = &t
		}
	}
}

func (cies *ContactImportExportService) extractVCardType(property string) string {
	// Extract type from property (e.g., "TEL;TYPE=WORK" -> "work")
	typeRegex := regexp.MustCompile(`TYPE=([^;:]+)`)
	matches := typeRegex.FindStringSubmatch(property)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return "other"
}

func (cies *ContactImportExportService) importVCardContact(vCard VCardContact, userID uuid.UUID, result *ImportResult, groupMap map[string]uuid.UUID, options map[string]interface{}) error {
	// Check for duplicate contacts
	isDuplicate, existingContactID := cies.isDuplicateContact(vCard, userID)
	if isDuplicate && options["skipDuplicates"].(bool) {
		result.SkippedCount++
		return nil
	}

	// Create or update contact
	var contactID uuid.UUID
	if isDuplicate {
		contactID = existingContactID
		result.DuplicateCount++
	} else {
		contactID = uuid.New()
	}

	// Create contact data
	contactData := map[string]interface{}{
		"id":         contactID,
		"user_id":    userID,
		"first_name": vCard.FirstName,
		"last_name":  vCard.LastName,
		"email":      vCard.Email,
		"phone":      vCard.Phone,
		"company":    vCard.Organization,
		"position":   vCard.Title,
		"notes":      vCard.Notes,
		"birthday":   vCard.Birthday,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	// Handle groups/categories
	if len(vCard.Categories) > 0 {
		var groupIDs []uuid.UUID
		for _, category := range vCard.Categories {
			groupID, exists := groupMap[category]
			if !exists {
				// Create new group
				group := &ContactGroup{
					Name:        category,
					Description: fmt.Sprintf("Imported from vCard - %s", category),
					UserID:      userID,
					IsPublic:    false,
				}
				err := cies.CreateContactGroup(group)
				if err != nil {
					continue
				}
				groupID = group.ID
				groupMap[category] = groupID
				result.NewGroups = append(result.NewGroups, category)
			}
			groupIDs = append(groupIDs, groupID)
		}
		contactData["groups"] = groupIDs
	}

	// Save contact (would call contact service)
	// This is a simplified implementation
	result.ImportedCount++

	return nil
}

func (cies *ContactImportExportService) isDuplicateContact(vCard VCardContact, userID uuid.UUID) (bool, uuid.UUID) {
	// Check for duplicates based on email/name
	if len(vCard.Email) > 0 {
		var existingID uuid.UUID
		err := cies.db.QueryRow(`
			SELECT id FROM contacts
			WHERE user_id = $1 AND email @> $2
			LIMIT 1
		`, userID, vCard.Email).Scan(&existingID)
		if err == nil {
			return true, existingID
		}
	}

	// Check by name
	if vCard.FullName != "" {
		var existingID uuid.UUID
		err := cies.db.QueryRow(`
			SELECT id FROM contacts
			WHERE user_id = $1 AND (
				first_name || ' ' || last_name = $2 OR
				full_name = $2
			)
			LIMIT 1
		`, userID, vCard.FullName).Scan(&existingID)
		if err == nil {
			return true, existingID
		}
	}

	return false, uuid.Nil
}

func (cies *ContactImportExportService) getContactsForExport(userID uuid.UUID, options ExportOptions) ([]VCardContact, error) {
	// Build query based on options
	query := `
		SELECT id, first_name, last_name, email, phone, company, position,
		       address, notes, birthday, groups, created_at
		FROM contacts
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	// Filter by specific contacts
	if len(options.ContactIDs) > 0 {
		placeholders := make([]string, len(options.ContactIDs))
		for i, contactID := range options.ContactIDs {
			placeholders[i] = fmt.Sprintf("$%d", argCount+1)
			args = append(args, contactID)
			argCount++
		}
		query += fmt.Sprintf(" AND id IN (%s)", strings.Join(placeholders, ","))
	}

	// Filter by groups
	if len(options.GroupIDs) > 0 {
		placeholders := make([]string, len(options.GroupIDs))
		for i, groupID := range options.GroupIDs {
			placeholders[i] = fmt.Sprintf("$%d", argCount+1)
			args = append(args, groupID)
			argCount++
		}
		query += fmt.Sprintf(" AND groups && ARRAY[%s]", strings.Join(placeholders, ","))
	}

	query += " ORDER BY first_name, last_name"

	rows, err := cies.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []VCardContact
	for rows.Next() {
		var contact VCardContact
		var firstName, lastName, email, phone, company, position, address, notes *string
		var birthday *time.Time
		var groups []uuid.UUID
		var createdAt time.Time

		err := rows.Scan(
			&contact.UID, &firstName, &lastName, &email, &phone, &company, &position,
			&address, &notes, &birthday, &groups, &createdAt,
		)
		if err != nil {
			continue
		}

		contact.FullName = cies.buildFullName(firstName, lastName)
		contact.FirstName = *firstName
		contact.LastName = *lastName
		if email != nil {
			contact.Email = []string{*email}
		}
		if phone != nil {
			contact.Phone = []PhoneNumber{{Number: *phone, Type: "mobile"}}
		}
		contact.Organization = *company
		contact.Title = *position
		contact.Notes = *notes
		contact.Birthday = birthday

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (cies *ContactImportExportService) exportToVCard(contacts []VCardContact, options ExportOptions) ([]byte, error) {
	var vCardData strings.Builder

	for _, contact := range contacts {
		vCardData.WriteString(cies.generateVCard(contact, options))
		vCardData.WriteString("\n")
	}

	return []byte(vCardData.String()), nil
}

func (cies *ContactImportExportService) generateVCard(contact VCardContact, options ExportOptions) string {
	var vCard strings.Builder

	vCard.WriteString("BEGIN:VCARD\n")
	vCard.WriteString("VERSION:3.0\n")

	if contact.UID != "" {
		vCard.WriteString(fmt.Sprintf("UID:%s\n", contact.UID))
	}

	if contact.FullName != "" {
		vCard.WriteString(fmt.Sprintf("FN:%s\n", contact.FullName))
	}

	if contact.LastName != "" || contact.FirstName != "" {
		vCard.WriteString(fmt.Sprintf("N:%s;%s;;;\n", contact.LastName, contact.FirstName))
	}

	for _, email := range contact.Email {
		vCard.WriteString(fmt.Sprintf("EMAIL:%s\n", email))
	}

	for _, phone := range contact.Phone {
		vCard.WriteString(fmt.Sprintf("TEL;TYPE=%s:%s\n", strings.ToUpper(phone.Type), phone.Number))
	}

	if contact.Organization != "" {
		vCard.WriteString(fmt.Sprintf("ORG:%s\n", contact.Organization))
	}

	if contact.Title != "" {
		vCard.WriteString(fmt.Sprintf("TITLE:%s\n", contact.Title))
	}

	if contact.Birthday != nil {
		vCard.WriteString(fmt.Sprintf("BDAY:%s\n", contact.Birthday.Format("20060102")))
	}

	if contact.Notes != "" {
		vCard.WriteString(fmt.Sprintf("NOTE:%s\n", contact.Notes))
	}

	if len(contact.Categories) > 0 {
		vCard.WriteString(fmt.Sprintf("CATEGORIES:%s\n", strings.Join(contact.Categories, ",")))
	}

	for _, url := range contact.URLs {
		vCard.WriteString(fmt.Sprintf("URL:%s\n", url))
	}

	vCard.WriteString("END:VCARD")

	return vCard.String()
}

func (cies *ContactImportExportService) exportToCSV(contacts []VCardContact, options ExportOptions) ([]byte, error) {
	var csv strings.Builder

	// CSV header
	csv.WriteString("Full Name,First Name,Last Name,Email,Phone,Organization,Title,Notes,Categories\n")

	for _, contact := range contacts {
		email := ""
		if len(contact.Email) > 0 {
			email = contact.Email[0]
		}

		phone := ""
		if len(contact.Phone) > 0 {
			phone = contact.Phone[0].Number
		}

		categories := strings.Join(contact.Categories, ";")

		// Escape commas and quotes in fields
		fields := []string{
			strings.ReplaceAll(contact.FullName, ",", ";"),
			strings.ReplaceAll(contact.FirstName, ",", ";"),
			strings.ReplaceAll(contact.LastName, ",", ";"),
			strings.ReplaceAll(email, ",", ";"),
			strings.ReplaceAll(phone, ",", ";"),
			strings.ReplaceAll(contact.Organization, ",", ";"),
			strings.ReplaceAll(contact.Title, ",", ";"),
			strings.ReplaceAll(contact.Notes, ",", ";"),
			categories,
		}

		csv.WriteString(strings.Join(fields, ",") + "\n")
	}

	return []byte(csv.String()), nil
}

func (cies *ContactImportExportService) exportToJSON(contacts []VCardContact, options ExportOptions) ([]byte, error) {
	// Simple JSON export
	result := map[string]interface{}{
		"contacts":    contacts,
		"exported_at": time.Now(),
		"format":      "json",
	}
	_ = options
	return json.Marshal(result)
}

func (cies *ContactImportExportService) buildFullName(firstName, lastName *string) string {
	if firstName != nil && lastName != nil {
		return *firstName + " " + *lastName
	} else if firstName != nil {
		return *firstName
	} else if lastName != nil {
		return *lastName
	}
	return ""
}
