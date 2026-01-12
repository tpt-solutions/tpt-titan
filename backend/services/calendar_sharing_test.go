package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalendarSharingPermissions tests basic permission validation
func TestCalendarSharingPermissions(t *testing.T) {
	sharingService := &CalendarSharingService{}

	// Test valid permissions
	validPermissions := []string{"view", "edit", "admin", "delegate"}
	for _, perm := range validPermissions {
		isValid := sharingService.ValidatePermission(perm)
		assert.True(t, isValid, "Permission %s should be valid", perm)
	}

	// Test invalid permissions
	invalidPermissions := []string{"", "read", "write", "owner", "invalid"}
	for _, perm := range invalidPermissions {
		isValid := sharingService.ValidatePermission(perm)
		assert.False(t, isValid, "Permission %s should be invalid", perm)
	}
}

// TestCalendarSharingHierarchy tests permission hierarchy
func TestCalendarSharingHierarchy(t *testing.T) {
	sharingService := &CalendarSharingService{}

	tests := []struct {
		userPermission string
		requiredPermission string
		hasAccess bool
	}{
		{"view", "view", true},
		{"edit", "view", true},
		{"edit", "edit", true},
		{"admin", "view", true},
		{"admin", "edit", true},
		{"admin", "admin", true},
		{"delegate", "view", true},
		{"delegate", "edit", true},
		{"delegate", "admin", true},
		{"delegate", "delegate", true},
		{"view", "edit", false},
		{"view", "admin", false},
		{"edit", "admin", false},
	}

	for _, tt := range tests {
		hasAccess := sharingService.HasPermission(tt.userPermission, tt.requiredPermission)
		assert.Equal(t, tt.hasAccess, hasAccess,
			"User with %s permission should %s have %s access",
			tt.userPermission,
			map[bool]string{true: "", false: "not"}[tt.hasAccess],
			tt.requiredPermission)
	}
}

// TestCalendarShareCreation tests creating calendar shares
func TestCalendarShareCreation(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()
	shareeID := uuid.New()

	// Test creating a share with view permission
	share, err := sharingService.CreateShare(calendarID, ownerID, shareeID, "view", "Test message")
	assert.NoError(t, err)
	assert.NotNil(t, share)
	assert.Equal(t, calendarID, share.CalendarID)
	assert.Equal(t, ownerID, share.OwnerID)
	assert.Equal(t, shareeID, share.ShareeID)
	assert.Equal(t, "view", share.Permission)
	assert.Equal(t, "Test message", share.Message)
	assert.True(t, share.AcceptedAt.IsZero()) // Not accepted yet
	assert.True(t, share.RevokedAt.IsZero()) // Not revoked

	// Test creating a share with invalid permission
	_, err = sharingService.CreateShare(calendarID, ownerID, shareeID, "invalid", "")
	assert.Error(t, err)
}

// TestCalendarShareAcceptance tests accepting calendar shares
func TestCalendarShareAcceptance(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()
	shareeID := uuid.New()

	// Create a share
	share, err := sharingService.CreateShare(calendarID, ownerID, shareeID, "edit", "")
	require.NoError(t, err)

	// Accept the share
	acceptedShare, err := sharingService.AcceptShare(share.ID, shareeID)
	assert.NoError(t, err)
	assert.NotNil(t, acceptedShare)
	assert.False(t, acceptedShare.AcceptedAt.IsZero())
	assert.True(t, acceptedShare.RevokedAt.IsZero())

	// Test accepting non-existent share
	_, err = sharingService.AcceptShare(uuid.New(), shareeID)
	assert.Error(t, err)

	// Test accepting share with wrong user
	wrongUserID := uuid.New()
	_, err = sharingService.AcceptShare(share.ID, wrongUserID)
	assert.Error(t, err)
}

// TestCalendarShareRevocation tests revoking calendar shares
func TestCalendarShareRevocation(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()
	shareeID := uuid.New()

	// Create and accept a share
	share, err := sharingService.CreateShare(calendarID, ownerID, shareeID, "admin", "")
	require.NoError(t, err)

	acceptedShare, err := sharingService.AcceptShare(share.ID, shareeID)
	require.NoError(t, err)

	// Revoke the share
	revokedShare, err := sharingService.RevokeShare(acceptedShare.ID, ownerID)
	assert.NoError(t, err)
	assert.NotNil(t, revokedShare)
	assert.False(t, revokedShare.RevokedAt.IsZero())

	// Test revoking non-existent share
	_, err = sharingService.RevokeShare(uuid.New(), ownerID)
	assert.Error(t, err)

	// Test revoking share with wrong user (not owner)
	wrongUserID := uuid.New()
	_, err = sharingService.RevokeShare(acceptedShare.ID, wrongUserID)
	assert.Error(t, err)
}

// TestCalendarShareListing tests listing calendar shares
func TestCalendarShareListing(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()

	// Create multiple shares for the same calendar
	shareeIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	shares := make([]*CalendarShare, 0)

	for i, shareeID := range shareeIDs {
		permissions := []string{"view", "edit", "admin"}
		share, err := sharingService.CreateShare(calendarID, ownerID, shareeID, permissions[i], "")
		require.NoError(t, err)
		shares = append(shares, share)
	}

	// List all shares for the calendar
	calendarShares, err := sharingService.ListSharesForCalendar(calendarID, ownerID)
	assert.NoError(t, err)
	assert.Len(t, calendarShares, 3)

	// List shares for a specific user
	userShares, err := sharingService.ListSharesForUser(shareeIDs[0])
	assert.NoError(t, err)
	assert.Len(t, userShares, 1)
	assert.Equal(t, shareeIDs[0], userShares[0].ShareeID)

	// Test listing shares for calendar with wrong owner
	_, err = sharingService.ListSharesForCalendar(calendarID, uuid.New())
	assert.Error(t, err)
}

// TestCalendarShareDelegation tests delegate permissions
func TestCalendarShareDelegation(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()
	delegateID := uuid.New()
	regularUserID := uuid.New()

	// Create delegate share
	delegateShare, err := sharingService.CreateShare(calendarID, ownerID, delegateID, "delegate", "")
	require.NoError(t, err)

	_, err = sharingService.AcceptShare(delegateShare.ID, delegateID)
	require.NoError(t, err)

	// Create regular edit share
	regularShare, err := sharingService.CreateShare(calendarID, ownerID, regularUserID, "edit", "")
	require.NoError(t, err)

	_, err = sharingService.AcceptShare(regularShare.ID, regularUserID)
	require.NoError(t, err)

	// Test delegate can manage shares
	canDelegateManage, err := sharingService.CanUserManageShares(calendarID, delegateID)
	assert.NoError(t, err)
	assert.True(t, canDelegateManage)

	// Test regular user cannot manage shares
	canRegularManage, err := sharingService.CanUserManageShares(calendarID, regularUserID)
	assert.NoError(t, err)
	assert.False(t, canRegularManage)

	// Test owner can manage shares
	canOwnerManage, err := sharingService.CanUserManageShares(calendarID, ownerID)
	assert.NoError(t, err)
	assert.True(t, canOwnerManage)
}

// TestCalendarShareExpiration tests share expiration
func TestCalendarShareExpiration(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()
	shareeID := uuid.New()

	// Create share with expiration
	expiresAt := time.Now().Add(24 * time.Hour)
	share, err := sharingService.CreateShareWithExpiration(calendarID, ownerID, shareeID, "view", "", expiresAt)
	require.NoError(t, err)

	// Accept the share
	_, err = sharingService.AcceptShare(share.ID, shareeID)
	require.NoError(t, err)

	// Test share is not expired yet
	isExpired, err := sharingService.IsShareExpired(share.ID)
	assert.NoError(t, err)
	assert.False(t, isExpired)

	// Create expired share
	pastExpiration := time.Now().Add(-24 * time.Hour)
	expiredShare, err := sharingService.CreateShareWithExpiration(calendarID, ownerID, uuid.New(), "view", "", pastExpiration)
	require.NoError(t, err)

	// Test expired share
	isExpired, err = sharingService.IsShareExpired(expiredShare.ID)
	assert.NoError(t, err)
	assert.True(t, isExpired)
}

// TestCalendarShareNotifications tests notification system
func TestCalendarShareNotifications(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()
	shareeID := uuid.New()

	// Create and accept a share
	share, err := sharingService.CreateShare(calendarID, ownerID, shareeID, "edit", "Please review this calendar")
	require.NoError(t, err)

	_, err = sharingService.AcceptShare(share.ID, shareeID)
	require.NoError(t, err)

	// Test notification creation
	notification, err := sharingService.CreateShareNotification(share.ID)
	assert.NoError(t, err)
	assert.NotNil(t, notification)
	assert.Equal(t, shareeID, notification.UserID)
	assert.Equal(t, "calendar_share", notification.Type)
	assert.Contains(t, notification.Message, "shared a calendar")

	// Test notification dismissal
	err = sharingService.DismissNotification(notification.ID, shareeID)
	assert.NoError(t, err)

	// Test dismissing non-existent notification
	err = sharingService.DismissNotification(uuid.New(), shareeID)
	assert.Error(t, err)

	// Test dismissing notification with wrong user
	err = sharingService.DismissNotification(notification.ID, uuid.New())
	assert.Error(t, err)
}

// TestCalendarShareAuditLogging tests audit trail
func TestCalendarShareAuditLogging(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()
	shareeID := uuid.New()

	// Create share
	share, err := sharingService.CreateShare(calendarID, ownerID, shareeID, "admin", "")
	require.NoError(t, err)

	// Accept share
	_, err = sharingService.AcceptShare(share.ID, shareeID)
	require.NoError(t, err)

	// Revoke share
	_, err = sharingService.RevokeShare(share.ID, ownerID)
	require.NoError(t, err)

	// Test audit log retrieval
	auditLogs, err := sharingService.GetAuditLogsForCalendar(calendarID, ownerID)
	assert.NoError(t, err)
	assert.Len(t, auditLogs, 3) // create, accept, revoke

	// Verify log entries
	expectedActions := []string{"created", "accepted", "revoked"}
	for i, log := range auditLogs {
		assert.Equal(t, expectedActions[i], log.Action)
		assert.Equal(t, calendarID, log.CalendarID)
		assert.NotZero(t, log.Timestamp)
	}

	// Test audit log access control
	_, err = sharingService.GetAuditLogsForCalendar(calendarID, uuid.New())
	assert.Error(t, err) // Wrong user should not access logs
}

// TestCalendarShareBulkOperations tests bulk share management
func TestCalendarShareBulkOperations(t *testing.T) {
	sharingService := &CalendarSharingService{}

	calendarID := uuid.New()
	ownerID := uuid.New()

	// Create multiple shares
	shareeIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New()}
	shareIDs := make([]uuid.UUID, 0)

	for _, shareeID := range shareeIDs {
		share, err := sharingService.CreateShare(calendarID, ownerID, shareeID, "view", "")
		require.NoError(t, err)
		shareIDs = append(shareIDs, share.ID)
	}

	// Bulk revoke shares
	err := sharingService.BulkRevokeShares(shareIDs[:2], ownerID)
	assert.NoError(t, err)

	// Check revocation status
	for i, shareID := range shareIDs {
		share, err := sharingService.GetShare(shareID)
		require.NoError(t, err)

		if i < 2 {
			assert.False(t, share.RevokedAt.IsZero(), "First two shares should be revoked")
		} else {
			assert.True(t, share.RevokedAt.IsZero(), "Last two shares should not be revoked")
		}
	}
}

// TestCalendarShareErrorHandling tests error conditions
func TestCalendarShareErrorHandling(t *testing.T) {
	sharingService := &CalendarSharingService{}

	// Test creating share with empty calendar ID
	_, err := sharingService.CreateShare(uuid.Nil, uuid.New(), uuid.New(), "view", "")
	assert.Error(t, err)

	// Test creating share with empty owner ID
	_, err = sharingService.CreateShare(uuid.New(), uuid.Nil, uuid.New(), "view", "")
	assert.Error(t, err)

	// Test creating share with empty sharee ID
	_, err = sharingService.CreateShare(uuid.New(), uuid.New(), uuid.Nil, "view", "")
	assert.Error(t, err)

	// Test self-sharing (owner sharing with themselves)
	ownerID := uuid.New()
	_, err = sharingService.CreateShare(uuid.New(), ownerID, ownerID, "view", "")
	assert.Error(t, err)

	// Test accepting already accepted share
	share, err := sharingService.CreateShare(uuid.New(), uuid.New(), uuid.New(), "view", "")
	require.NoError(t, err)

	_, err = sharingService.AcceptShare(share.ID, share.ShareeID)
	require.NoError(t, err)

	// Try to accept again
	_, err = sharingService.AcceptShare(share.ID, share.ShareeID)
	assert.Error(t, err)
}
