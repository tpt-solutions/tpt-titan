// backend/services/calendar_sharing_test.go
// Run with: cd backend && go test ./services/... -run TestCalendar -v

package services

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

// ─── Permission constants ─────────────────────────────────────────────────────

func TestPermissionConstants_Defined(t *testing.T) {
	if PermissionView != "view" {
		t.Errorf("PermissionView = %q, want %q", PermissionView, "view")
	}
	if PermissionEdit != "edit" {
		t.Errorf("PermissionEdit = %q, want %q", PermissionEdit, "edit")
	}
	if PermissionAdmin != "admin" {
		t.Errorf("PermissionAdmin = %q, want %q", PermissionAdmin, "admin")
	}
	if PermissionDelegate != "delegate" {
		t.Errorf("PermissionDelegate = %q, want %q", PermissionDelegate, "delegate")
	}
}

// ─── permissionSufficient ─────────────────────────────────────────────────────
// Tests the unexported helper directly (same package).

func TestPermissionSufficient_ExactMatch(t *testing.T) {
	svc := &CalendarSharingService{}

	if !svc.permissionSufficient(PermissionView, PermissionView) {
		t.Error("view >= view should be true")
	}
	if !svc.permissionSufficient(PermissionEdit, PermissionEdit) {
		t.Error("edit >= edit should be true")
	}
	if !svc.permissionSufficient(PermissionAdmin, PermissionAdmin) {
		t.Error("admin >= admin should be true")
	}
	if !svc.permissionSufficient(PermissionDelegate, PermissionDelegate) {
		t.Error("delegate >= delegate should be true")
	}
}

func TestPermissionSufficient_HigherAllowsLower(t *testing.T) {
	svc := &CalendarSharingService{}

	cases := []struct {
		user     SharePermission
		required SharePermission
		want     bool
	}{
		// edit can do view
		{PermissionEdit, PermissionView, true},
		// admin can do view
		{PermissionAdmin, PermissionView, true},
		// admin can do edit
		{PermissionAdmin, PermissionEdit, true},
		// delegate can do view
		{PermissionDelegate, PermissionView, true},
		// delegate can do edit
		{PermissionDelegate, PermissionEdit, true},
		// delegate CANNOT do admin (delegate=3, admin=4)
		{PermissionDelegate, PermissionAdmin, false},
		// view cannot do edit
		{PermissionView, PermissionEdit, false},
		// view cannot do admin
		{PermissionView, PermissionAdmin, false},
		// edit cannot do admin
		{PermissionEdit, PermissionAdmin, false},
		// edit cannot do delegate
		{PermissionEdit, PermissionDelegate, false},
	}

	for _, tc := range cases {
		got := svc.permissionSufficient(tc.user, tc.required)
		if got != tc.want {
			t.Errorf("permissionSufficient(%q, %q) = %v, want %v",
				tc.user, tc.required, got, tc.want)
		}
	}
}

func TestPermissionSufficient_UnknownPermission(t *testing.T) {
	svc := &CalendarSharingService{}

	// Unknown permission levels should return false
	if svc.permissionSufficient("unknown", PermissionView) {
		t.Error("unknown user permission should return false")
	}
	if svc.permissionSufficient(PermissionView, "unknown") {
		t.Error("unknown required permission should return false")
	}
}

// ─── generateInviteToken ──────────────────────────────────────────────────────

func TestGenerateInviteToken_ReturnsValidUUID(t *testing.T) {
	svc := &CalendarSharingService{}

	token := svc.generateInviteToken()
	if token == "" {
		t.Fatal("expected non-empty invite token")
	}

	// Should be parseable as a UUID
	if _, err := uuid.Parse(token); err != nil {
		t.Errorf("token %q is not a valid UUID: %v", token, err)
	}
}

func TestGenerateInviteToken_UniqueEachCall(t *testing.T) {
	svc := &CalendarSharingService{}

	t1 := svc.generateInviteToken()
	t2 := svc.generateInviteToken()

	if t1 == t2 {
		t.Error("invite tokens should be unique across calls")
	}
}

// ─── CalendarShare struct ─────────────────────────────────────────────────────

func TestCalendarShare_ZeroValue(t *testing.T) {
	var share CalendarShare

	if share.ID != uuid.Nil {
		t.Error("zero value ID should be uuid.Nil")
	}
	if share.CanInviteOthers {
		t.Error("zero value CanInviteOthers should be false")
	}
	if share.Status != "" {
		t.Error("zero value Status should be empty")
	}
}

func TestCalendarShare_Construction(t *testing.T) {
	calendarID := uuid.New()
	ownerID := uuid.New()
	sharedWithID := uuid.New()

	share := CalendarShare{
		ID:              uuid.New(),
		CalendarID:      calendarID,
		OwnerID:         ownerID,
		SharedWithID:    sharedWithID,
		Permission:      PermissionEdit,
		CanInviteOthers: true,
		Status:          "pending",
	}

	if share.CalendarID != calendarID {
		t.Errorf("CalendarID mismatch")
	}
	if share.Permission != PermissionEdit {
		t.Errorf("Permission = %q, want %q", share.Permission, PermissionEdit)
	}
	if !share.CanInviteOthers {
		t.Error("CanInviteOthers should be true")
	}
}

// ─── SharingInvite struct ─────────────────────────────────────────────────────

func TestSharingInvite_Construction(t *testing.T) {
	invite := SharingInvite{
		ID:         uuid.New(),
		CalendarID: uuid.New(),
		Email:      "test@example.com",
		Permission: PermissionView,
		Status:     "pending",
	}

	if invite.Email != "test@example.com" {
		t.Errorf("Email = %q, want test@example.com", invite.Email)
	}
	if invite.Status != "pending" {
		t.Errorf("Status = %q, want pending", invite.Status)
	}
}

// ─── CalendarGroup struct ─────────────────────────────────────────────────────

func TestCalendarGroup_Members(t *testing.T) {
	members := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	group := CalendarGroup{
		ID:         uuid.New(),
		Name:       "Finance Team",
		Members:    members,
		Permission: PermissionEdit,
		IsPublic:   false,
	}

	if len(group.Members) != 3 {
		t.Errorf("expected 3 members, got %d", len(group.Members))
	}
	if group.IsPublic {
		t.Error("group should not be public")
	}
}

// ─── NewCalendarSharingService ────────────────────────────────────────────────

func TestNewCalendarSharingService_NilDB(t *testing.T) {
	// Service can be constructed with a nil db for unit tests
	// (it will panic on any DB call, but construction itself is fine)
	svc := NewCalendarSharingService(nil)
	if svc == nil {
		t.Fatal("expected non-nil CalendarSharingService")
	}
}

// ─── canUserShareCalendar via permission level ────────────────────────────────
// We test the pure permission logic separately since canUserShareCalendar needs DB.

func TestAdminAndDelegateCanShare(t *testing.T) {
	svc := &CalendarSharingService{}

	// Admin and Delegate level should be sufficient to share (level >= 3)
	if !svc.permissionSufficient(PermissionAdmin, PermissionDelegate) {
		// admin=4 >= delegate=3 → should be true
		// Actually: delegate=3, admin=4 so admin IS sufficient for delegate level
		t.Error("admin should be >= delegate")
	}
}

// ─── CalendarDomainSharing struct ────────────────────────────────────────────

func TestCalendarDomainSharing_Fields(t *testing.T) {
	ds := CalendarDomainSharing{
		ID:         uuid.New(),
		CalendarID: uuid.New(),
		Domain:     "company.com",
		Permission: PermissionView,
		AutoAccept: true,
		CreatedBy:  uuid.New(),
	}

	if ds.Domain != "company.com" {
		t.Errorf("Domain = %q, want company.com", ds.Domain)
	}
	if !ds.AutoAccept {
		t.Error("AutoAccept should be true")
	}
	if !strings.Contains(ds.Domain, ".") {
		t.Error("domain should contain a dot")
	}
}

// ─── CalendarACL struct ───────────────────────────────────────────────────────

func TestCalendarACL_Fields(t *testing.T) {
	acl := CalendarACL{
		UserID:     uuid.New(),
		CalendarID: uuid.New(),
		Permission: PermissionEdit,
		Inherited:  false,
		Source:     "direct",
	}

	if acl.Source != "direct" {
		t.Errorf("Source = %q, want direct", acl.Source)
	}
	if acl.Inherited {
		t.Error("Inherited should be false")
	}
}
