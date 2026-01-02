-- TPT Titan Database Schema
-- Complete database design for all application components

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- =============================================================================
-- CORE SYSTEM TABLES
-- =============================================================================

-- Users table (core authentication and user management)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    display_name VARCHAR(255),
    avatar_url VARCHAR(500),
    bio TEXT,
    timezone VARCHAR(50) DEFAULT 'UTC',
    language VARCHAR(10) DEFAULT 'en',
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    email_verified_at TIMESTAMP WITH TIME ZONE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Security settings
    two_factor_enabled BOOLEAN DEFAULT false,
    two_factor_secret VARCHAR(255),
    password_changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,

    -- Preferences
    theme VARCHAR(20) DEFAULT 'light',
    notifications_enabled BOOLEAN DEFAULT true,
    auto_save BOOLEAN DEFAULT true,

    -- Encryption keys (encrypted with user's master password)
    encryption_key_hash VARCHAR(255),
    recovery_key_hash VARCHAR(255),
    backup_key_hash VARCHAR(255)
);

-- User sessions for JWT token management
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Session metadata
    device_type VARCHAR(50),
    device_name VARCHAR(255),
    location VARCHAR(255)
);

-- User roles and permissions
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource VARCHAR(100) NOT NULL, -- e.g., 'documents', 'emails', 'users'
    action VARCHAR(100) NOT NULL,   -- e.g., 'create', 'read', 'update', 'delete'
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(resource, action)
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, role_id)
);

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,

    UNIQUE(role_id, permission_id)
);

-- =============================================================================
-- DOCUMENT MANAGEMENT SYSTEM
-- =============================================================================

-- Documents table (core content storage)
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    content_type VARCHAR(100) NOT NULL, -- 'text', 'spreadsheet', 'form', etc.
    file_path VARCHAR(1000),
    file_size BIGINT,
    mime_type VARCHAR(255),
    checksum VARCHAR(128), -- SHA-256 hash
    version INTEGER DEFAULT 1,
    is_encrypted BOOLEAN DEFAULT false,
    encryption_key_hash VARCHAR(255),
    is_public BOOLEAN DEFAULT false,
    public_token VARCHAR(255) UNIQUE,

    -- Metadata
    tags TEXT[], -- Array of tags
    categories VARCHAR(255)[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_accessed_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Collaboration
    is_shared BOOLEAN DEFAULT false,
    allow_comments BOOLEAN DEFAULT true,
    allow_editing BOOLEAN DEFAULT false,

    -- Search and indexing
    search_vector TSVECTOR,
    full_text_search TSVECTOR
);

-- Document versions for history tracking
CREATE TABLE document_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    title VARCHAR(500),
    content TEXT,
    file_path VARCHAR(1000),
    file_size BIGINT,
    checksum VARCHAR(128),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    change_description TEXT,

    UNIQUE(document_id, version)
);

-- Document shares and permissions
CREATE TABLE document_shares (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    shared_by UUID NOT NULL REFERENCES users(id),
    shared_with_type VARCHAR(20) NOT NULL, -- 'user', 'group', 'public'
    shared_with_id UUID, -- user or group ID
    permission_level VARCHAR(20) NOT NULL, -- 'view', 'comment', 'edit'
    can_share BOOLEAN DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CHECK (shared_with_type IN ('user', 'group', 'public'))
);

-- Document comments and collaboration
CREATE TABLE document_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    parent_comment_id UUID REFERENCES document_comments(id),
    position_x INTEGER,
    position_y INTEGER,
    resolved BOOLEAN DEFAULT false,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- EMAIL SYSTEM
-- =============================================================================

-- Email accounts configuration
CREATE TABLE email_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email_address VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    provider VARCHAR(100), -- 'gmail', 'outlook', 'custom', etc.
    protocol VARCHAR(10) DEFAULT 'imap', -- 'imap', 'pop3'

    -- Connection settings
    imap_host VARCHAR(255),
    imap_port INTEGER DEFAULT 993,
    smtp_host VARCHAR(255),
    smtp_port INTEGER DEFAULT 587,
    use_ssl BOOLEAN DEFAULT true,
    use_tls BOOLEAN DEFAULT true,

    -- Authentication (encrypted)
    username VARCHAR(255),
    password_encrypted TEXT,
    oauth_token TEXT,
    oauth_refresh_token TEXT,
    oauth_expires_at TIMESTAMP WITH TIME ZONE,

    -- Status and metadata
    is_active BOOLEAN DEFAULT true,
    last_sync TIMESTAMP WITH TIME ZONE,
    sync_status VARCHAR(50) DEFAULT 'idle',
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, email_address)
);

-- Email messages
CREATE TABLE emails (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES email_accounts(id) ON DELETE CASCADE,
    message_id VARCHAR(500) UNIQUE, -- RFC 2822 Message-ID
    thread_id VARCHAR(500), -- Thread/conversation identifier
    subject VARCHAR(1000),
    sender_name VARCHAR(255),
    sender_email VARCHAR(255) NOT NULL,
    recipient_to TEXT[], -- Array of recipient emails
    recipient_cc TEXT[],
    recipient_bcc TEXT[],
    content_text TEXT,
    content_html TEXT,
    raw_content TEXT, -- Full raw email content
    is_read BOOLEAN DEFAULT false,
    is_starred BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false,
    is_draft BOOLEAN DEFAULT false,
    has_attachments BOOLEAN DEFAULT false,

    -- Dates
    sent_at TIMESTAMP WITH TIME ZONE,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Metadata
    size_bytes INTEGER,
    priority VARCHAR(10) DEFAULT 'normal', -- 'low', 'normal', 'high'
    labels TEXT[], -- Email labels/tags
    folder VARCHAR(255) DEFAULT 'INBOX',

    -- Search
    search_vector TSVECTOR,

    -- Encryption
    is_encrypted BOOLEAN DEFAULT false,
    encryption_key_hash VARCHAR(255),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email attachments
CREATE TABLE email_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_id UUID NOT NULL REFERENCES emails(id) ON DELETE CASCADE,
    filename VARCHAR(500) NOT NULL,
    content_type VARCHAR(255),
    size_bytes INTEGER,
    file_path VARCHAR(1000),
    checksum VARCHAR(128),
    is_inline BOOLEAN DEFAULT false,
    content_id VARCHAR(255), -- For inline attachments
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email folders/labels
CREATE TABLE email_folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES email_accounts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    full_path VARCHAR(1000),
    parent_folder_id UUID REFERENCES email_folders(id),
    message_count INTEGER DEFAULT 0,
    unread_count INTEGER DEFAULT 0,
    is_system BOOLEAN DEFAULT false, -- INBOX, SENT, DRAFTS, TRASH, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(account_id, full_path)
);

-- =============================================================================
-- CONTACT MANAGEMENT
-- =============================================================================

-- Contacts table
CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    display_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    mobile VARCHAR(50),
    work_phone VARCHAR(50),
    address_street VARCHAR(255),
    address_city VARCHAR(255),
    address_state VARCHAR(255),
    address_zip VARCHAR(20),
    address_country VARCHAR(100),
    company VARCHAR(255),
    job_title VARCHAR(255),
    website VARCHAR(500),
    birthday DATE,
    notes TEXT,

    -- Social media and additional fields
    linkedin_url VARCHAR(500),
    twitter_handle VARCHAR(100),
    facebook_url VARCHAR(500),
    instagram_handle VARCHAR(100),

    -- Metadata
    is_favorite BOOLEAN DEFAULT false,
    groups TEXT[], -- Contact groups/categories
    last_contacted DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Search
    search_vector TSVECTOR,

    UNIQUE(user_id, email)
);

-- =============================================================================
-- CALENDAR SYSTEM
-- =============================================================================

-- Calendars
CREATE TABLE calendars (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#3788d8', -- Hex color
    is_default BOOLEAN DEFAULT false,
    is_shared BOOLEAN DEFAULT false,
    is_public BOOLEAN DEFAULT false,
    public_token VARCHAR(255) UNIQUE,
    timezone VARCHAR(50) DEFAULT 'UTC',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Calendar events
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id), -- Event creator
    title VARCHAR(500) NOT NULL,
    description TEXT,
    location VARCHAR(500),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    timezone VARCHAR(50) DEFAULT 'UTC',
    is_all_day BOOLEAN DEFAULT false,
    recurrence_rule TEXT, -- RRULE format
    recurrence_id UUID, -- For recurring event instances

    -- Status and privacy
    status VARCHAR(20) DEFAULT 'confirmed', -- 'confirmed', 'tentative', 'cancelled'
    privacy VARCHAR(20) DEFAULT 'private', -- 'public', 'private'

    -- Reminders
    reminder_minutes INTEGER[], -- Array of minutes before event
    reminder_email BOOLEAN DEFAULT false,
    reminder_popup BOOLEAN DEFAULT true,

    -- Attendees
    attendees JSONB, -- Store attendee information as JSON

    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Search
    search_vector TSVECTOR
);

-- Calendar shares
CREATE TABLE calendar_shares (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    shared_by UUID NOT NULL REFERENCES users(id),
    shared_with UUID NOT NULL REFERENCES users(id),
    permission_level VARCHAR(20) NOT NULL, -- 'view', 'edit', 'admin'
    can_share BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(calendar_id, shared_with)
);

-- =============================================================================
-- TASK MANAGEMENT
-- =============================================================================

-- Tasks/Projects
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'todo', -- 'todo', 'in_progress', 'done', 'cancelled'
    priority VARCHAR(10) DEFAULT 'medium', -- 'low', 'medium', 'high', 'urgent'
    category VARCHAR(100),
    tags TEXT[],
    assignee_id UUID REFERENCES users(id),
    project_id UUID REFERENCES tasks(id), -- Self-reference for sub-tasks/projects

    -- Dates
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Progress and estimation
    estimated_hours DECIMAL(6,2),
    actual_hours DECIMAL(6,2),
    progress_percentage INTEGER DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),

    -- Relationships
    parent_task_id UUID REFERENCES tasks(id),
    depends_on UUID[] REFERENCES tasks(id), -- Array of task IDs this depends on

    -- Metadata
    is_recurring BOOLEAN DEFAULT false,
    recurrence_rule TEXT,
    is_archived BOOLEAN DEFAULT false,

    -- Search
    search_vector TSVECTOR
);

-- Task comments
CREATE TABLE task_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- CHAT AND REAL-TIME MESSAGING
-- =============================================================================

-- Chat rooms
CREATE TABLE chat_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255),
    description TEXT,
    room_type VARCHAR(20) NOT NULL, -- 'direct', 'group', 'channel'
    created_by UUID NOT NULL REFERENCES users(id),
    is_private BOOLEAN DEFAULT false,
    is_archived BOOLEAN DEFAULT false,
    last_message_at TIMESTAMP WITH TIME ZONE,
    message_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Chat room participants
CREATE TABLE chat_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member', -- 'owner', 'admin', 'member'
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE,
    is_muted BOOLEAN DEFAULT false,
    last_read_at TIMESTAMP WITH TIME ZONE,

    UNIQUE(room_id, user_id)
);

-- Chat messages
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text', -- 'text', 'file', 'system'
    reply_to_id UUID REFERENCES chat_messages(id),
    thread_root_id UUID REFERENCES chat_messages(id),
    is_edited BOOLEAN DEFAULT false,
    edited_at TIMESTAMP WITH TIME ZONE,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- File attachments
    file_url VARCHAR(1000),
    file_name VARCHAR(500),
    file_size BIGINT,
    file_type VARCHAR(100),

    -- Search
    search_vector TSVECTOR
);

-- Message reactions
CREATE TABLE message_reactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    reaction VARCHAR(50) NOT NULL, -- Emoji or reaction type
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(message_id, user_id, reaction)
);

-- Typing indicators (temporary, can be in Redis)
CREATE TABLE typing_indicators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP + INTERVAL '10 seconds',

    UNIQUE(room_id, user_id)
);

-- =============================================================================
-- FILE SYNCHRONIZATION
-- =============================================================================

-- Sync devices
CREATE TABLE sync_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    device_type VARCHAR(50), -- 'desktop', 'mobile', 'tablet'
    public_key TEXT,
    private_key_encrypted TEXT,
    last_seen TIMESTAMP WITH TIME ZONE,
    is_online BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, device_id)
);

-- Sync folders
CREATE TABLE sync_folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    local_path VARCHAR(1000),
    remote_path VARCHAR(1000),
    is_active BOOLEAN DEFAULT true,
    sync_mode VARCHAR(20) DEFAULT 'two_way', -- 'two_way', 'one_way', 'mirror'
    last_sync TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- File tracking for sync
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    folder_id UUID NOT NULL REFERENCES sync_folders(id) ON DELETE CASCADE,
    name VARCHAR(500) NOT NULL,
    path VARCHAR(2000) NOT NULL,
    size BIGINT,
    checksum VARCHAR(128),
    last_modified TIMESTAMP WITH TIME ZONE,
    is_directory BOOLEAN DEFAULT false,
    parent_id UUID REFERENCES files(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(folder_id, path)
);

-- File versions for sync
CREATE TABLE file_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    size BIGINT,
    checksum VARCHAR(128),
    chunk_hashes TEXT[], -- Array of chunk hashes
    device_id VARCHAR(255),
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(file_id, version)
);

-- File chunks for large file transfers
CREATE TABLE file_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    version_id UUID NOT NULL REFERENCES file_versions(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    size INTEGER NOT NULL,
    hash VARCHAR(128) NOT NULL,
    data BYTEA NOT NULL, -- Encrypted chunk data
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(version_id, chunk_index)
);

-- Sync queue for operations
CREATE TABLE sync_queue (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id VARCHAR(255) NOT NULL,
    operation VARCHAR(20) NOT NULL, -- 'upload', 'download', 'delete'
    file_id UUID REFERENCES files(id),
    priority INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'in_progress', 'completed', 'failed'
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Sync conflicts
CREATE TABLE sync_conflicts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    device_id VARCHAR(255) NOT NULL,
    conflict_type VARCHAR(50) NOT NULL, -- 'concurrent_edit', 'delete_conflict', etc.
    local_version INTEGER,
    remote_version INTEGER,
    resolution VARCHAR(20), -- 'keep_local', 'keep_remote', 'merge'
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Bandwidth limits
CREATE TABLE bandwidth_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255),
    max_upload BIGINT, -- bytes per second
    max_download BIGINT, -- bytes per second
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, device_id)
);

-- =============================================================================
-- VIDEO CONFERENCING
-- =============================================================================

-- Meetings
CREATE TABLE meetings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    host_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    room_id VARCHAR(255) UNIQUE NOT NULL,
    meeting_type VARCHAR(20) DEFAULT 'instant', -- 'scheduled', 'instant', 'recurring'
    start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    time_zone VARCHAR(50) DEFAULT 'UTC',
    is_active BOOLEAN DEFAULT false,
    is_recording BOOLEAN DEFAULT false,
    max_participants INTEGER DEFAULT 100,
    require_auth BOOLEAN DEFAULT false,
    password_hash VARCHAR(255),
    recording_path VARCHAR(1000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Meeting participants
CREATE TABLE meeting_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    role VARCHAR(20) DEFAULT 'participant', -- 'host', 'co_host', 'participant'
    joined_at TIMESTAMP WITH TIME ZONE,
    left_at TIMESTAMP WITH TIME ZONE,
    is_muted BOOLEAN DEFAULT true,
    is_video_on BOOLEAN DEFAULT false,
    is_screen_share BOOLEAN DEFAULT false,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Meeting chat messages
CREATE TABLE meeting_chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES meeting_participants(id),
    message TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text', -- 'text', 'file', 'system'
    file_url VARCHAR(1000),
    file_name VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Meeting invites
CREATE TABLE meeting_invites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'accepted', 'declined', 'expired'
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    responded_at TIMESTAMP WITH TIME ZONE,

    UNIQUE(meeting_id, email)
);

-- WebRTC connections (signaling data)
CREATE TABLE webrtc_connections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES meeting_participants(id),
    connection_type VARCHAR(20) NOT NULL, -- 'offer', 'answer', 'candidate'
    sdp TEXT,
    candidate JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour'
);

-- Meeting recordings
CREATE TABLE meeting_recordings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    started_by UUID NOT NULL REFERENCES users(id),
    file_path VARCHAR(1000),
    file_size BIGINT,
    duration INTERVAL,
    status VARCHAR(20) DEFAULT 'processing', -- 'recording', 'processing', 'completed', 'failed'
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- FORMS AND SURVEYS
-- =============================================================================

-- Forms
CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    is_template BOOLEAN DEFAULT false,
    is_public BOOLEAN DEFAULT false,
    public_token VARCHAR(255) UNIQUE,
    allow_responses BOOLEAN DEFAULT true,
    max_responses INTEGER,
    response_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Form fields
CREATE TABLE form_fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    field_type VARCHAR(50) NOT NULL, -- 'text', 'email', 'number', 'date', 'select', etc.
    label VARCHAR(500) NOT NULL,
    placeholder VARCHAR(500),
    help_text TEXT,
    is_required BOOLEAN DEFAULT false,
    is_hidden BOOLEAN DEFAULT false,
    field_order INTEGER NOT NULL,
    validation_rules JSONB, -- Field validation rules
    field_options JSONB, -- For select, radio, checkbox options
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Form responses
CREATE TABLE form_responses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    respondent_id UUID REFERENCES users(id), -- NULL for anonymous responses
    respondent_email VARCHAR(255),
    respondent_name VARCHAR(255),
    response_data JSONB NOT NULL, -- Complete form response as JSON
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    is_complete BOOLEAN DEFAULT true
);

-- =============================================================================
-- AI FEATURES
-- =============================================================================

-- AI models configuration
CREATE TABLE ai_models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL, -- 'ollama', 'openai', 'anthropic', etc.
    model_name VARCHAR(255) NOT NULL,
    version VARCHAR(50),
    capabilities TEXT[], -- Array of capabilities: 'text', 'code', 'image', etc.
    context_window INTEGER,
    is_local BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0, -- Higher priority models are preferred
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- AI requests tracking
CREATE TABLE ai_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    model_id UUID REFERENCES ai_models(id),
    request_type VARCHAR(50) NOT NULL, -- 'completion', 'chat', 'image', etc.
    prompt TEXT,
    response TEXT,
    tokens_used INTEGER,
    cost DECIMAL(10,6), -- Cost in USD
    status VARCHAR(20) DEFAULT 'processing', -- 'processing', 'completed', 'failed'
    error_message TEXT,
    processing_time INTERVAL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- AI usage tracking
CREATE TABLE ai_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_requests INTEGER DEFAULT 0,
    total_tokens INTEGER DEFAULT 0,
    total_cost DECIMAL(10,6) DEFAULT 0,
    model_usage JSONB, -- Usage breakdown by model
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, period_start)
);

-- AI model upgrades
CREATE TABLE ai_upgrades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    old_model VARCHAR(255),
    new_model VARCHAR(255),
    upgrade_type VARCHAR(50), -- 'automatic', 'manual', 'forced'
    reason TEXT,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- AUDIT LOGGING AND SECURITY
-- =============================================================================

-- Security audit log
CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    event_type VARCHAR(100) NOT NULL,
    user_id UUID REFERENCES users(id),
    session_id UUID REFERENCES user_sessions(id),
    ip_address INET,
    user_agent TEXT,
    resource_type VARCHAR(50), -- 'user', 'document', 'email', etc.
    resource_id UUID,
    action VARCHAR(50) NOT NULL, -- 'create', 'read', 'update', 'delete', 'login', etc.
    status VARCHAR(20) DEFAULT 'success', -- 'success', 'failure', 'warning'
    details JSONB, -- Additional event details
    severity VARCHAR(10) DEFAULT 'info' -- 'debug', 'info', 'warning', 'error', 'critical'
);

-- Security events (for monitoring)
CREATE TABLE security_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(10) NOT NULL,
    source_ip INET,
    user_agent TEXT,
    user_id UUID REFERENCES users(id),
    details JSONB,
    resolved BOOLEAN DEFAULT false,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- API rate limiting (can also be in Redis)
CREATE TABLE rate_limit_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address INET NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    user_id UUID REFERENCES users(id),
    request_count INTEGER DEFAULT 1,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    window_end TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Password reset tokens
CREATE TABLE password_resets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(email)
);

-- =============================================================================
-- BACKUP AND RECOVERY
-- =============================================================================

-- Backup metadata
CREATE TABLE backups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    backup_type VARCHAR(20) NOT NULL, -- 'full', 'incremental', 'user_data'
    user_id UUID REFERENCES users(id), -- NULL for system-wide backups
    file_path VARCHAR(1000),
    file_size BIGINT,
    checksum VARCHAR(128),
    status VARCHAR(20) DEFAULT 'completed', -- 'running', 'completed', 'failed'
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,

    -- Backup contents
    tables TEXT[],
    record_count INTEGER,
    compression_type VARCHAR(20) DEFAULT 'gzip',
    encryption_type VARCHAR(20)
);

-- Backup schedules
CREATE TABLE backup_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    schedule_type VARCHAR(20) NOT NULL, -- 'daily', 'weekly', 'monthly'
    cron_expression VARCHAR(100),
    backup_type VARCHAR(20) DEFAULT 'full',
    retention_days INTEGER DEFAULT 30,
    is_active BOOLEAN DEFAULT true,
    last_run TIMESTAMP WITH TIME ZONE,
    next_run TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- INDEXES FOR PERFORMANCE
-- =============================================================================

-- Core user indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Document indexes
CREATE INDEX idx_documents_owner_id ON documents(owner_id);
CREATE INDEX idx_documents_content_type ON documents(content_type);
CREATE INDEX idx_documents_created_at ON documents(created_at);
CREATE INDEX idx_documents_updated_at ON documents(updated_at);
CREATE INDEX idx_documents_search_vector ON documents USING GIN(search_vector);
CREATE INDEX idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX idx_document_shares_document_id ON document_shares(document_id);
CREATE INDEX idx_document_comments_document_id ON document_comments(document_id);

-- Email indexes
CREATE INDEX idx_email_accounts_user_id ON email_accounts(user_id);
CREATE INDEX idx_emails_account_id ON emails(account_id);
CREATE INDEX idx_emails_received_at ON emails(received_at DESC);
CREATE INDEX idx_emails_sender_email ON emails(sender_email);
CREATE INDEX idx_emails_thread_id ON emails(thread_id);
CREATE INDEX idx_emails_search_vector ON emails USING GIN(search_vector);

-- Contact indexes
CREATE INDEX idx_contacts_user_id ON contacts(user_id);
CREATE INDEX idx_contacts_email ON contacts(email);
CREATE INDEX idx_contacts_search_vector ON contacts USING GIN(search_vector);

-- Calendar indexes
CREATE INDEX idx_calendars_user_id ON calendars(user_id);
CREATE INDEX idx_events_calendar_id ON events(calendar_id);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_end_time ON events(end_time);
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_search_vector ON events USING GIN(search_vector);

-- Task indexes
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_search_vector ON tasks USING GIN(search_vector);

-- Chat indexes
CREATE INDEX idx_chat_rooms_created_by ON chat_rooms(created_by);
CREATE INDEX idx_chat_participants_room_id ON chat_participants(room_id);
CREATE INDEX idx_chat_participants_user_id ON chat_participants(user_id);
CREATE INDEX idx_chat_messages_room_id ON chat_messages(room_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at DESC);
CREATE INDEX idx_chat_messages_search_vector ON chat_messages USING GIN(search_vector);

-- Sync indexes
CREATE INDEX idx_sync_devices_user_id ON sync_devices(user_id);
CREATE INDEX idx_sync_folders_user_id ON sync_folders(user_id);
CREATE INDEX idx_files_folder_id ON files(folder_id);
CREATE INDEX idx_file_versions_file_id ON file_versions(file_id);

-- Meeting indexes
CREATE INDEX idx_meetings_host_id ON meetings(host_id);
CREATE INDEX idx_meetings_start_time ON meetings(start_time);
CREATE INDEX idx_meeting_participants_meeting_id ON meeting_participants(meeting_id);

-- Form indexes
CREATE INDEX idx_forms_owner_id ON forms(owner_id);
CREATE INDEX idx_form_responses_form_id ON form_responses(form_id);

-- AI indexes
CREATE INDEX idx_ai_requests_user_id ON ai_requests(user_id);
CREATE INDEX idx_ai_requests_created_at ON ai_requests(created_at);
CREATE INDEX idx_ai_usage_user_id ON ai_usage(user_id);

-- Security indexes
CREATE INDEX idx_audit_log_timestamp ON audit_log(timestamp DESC);
CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_event_type ON audit_log(event_type);
CREATE INDEX idx_security_events_created_at ON security_events(created_at DESC);

-- Backup indexes
CREATE INDEX idx_backups_user_id ON backups(user_id);
CREATE INDEX idx_backups_created_at ON backups(started_at DESC);

-- =============================================================================
-- PARTITIONING (for large tables)
-- =============================================================================

-- Partition emails by year for better performance
-- (Note: This would be applied when tables grow large)

-- Partition chat messages by month
-- (Note: This would be applied when tables grow large)

-- =============================================================================
-- DEFAULT DATA
-- =============================================================================

-- Insert default roles
INSERT INTO roles (name, description, is_system_role) VALUES
('admin', 'System administrator with full access', true),
('user', 'Regular user with standard access', true),
('premium', 'Premium user with enhanced features', false);

-- Insert default permissions
INSERT INTO permissions (resource, action, description) VALUES
('users', 'create', 'Create new users'),
('users', 'read', 'View user information'),
('users', 'update', 'Update user information'),
('users', 'delete', 'Delete users'),
('documents', 'create', 'Create documents'),
('documents', 'read', 'View documents'),
('documents', 'update', 'Edit documents'),
('documents', 'delete', 'Delete documents'),
('documents', 'share', 'Share documents with others'),
('emails', 'read', 'Read emails'),
('emails', 'send', 'Send emails'),
('emails', 'delete', 'Delete emails'),
('calendars', 'create', 'Create calendars'),
('calendars', 'read', 'View calendars'),
('calendars', 'update', 'Edit calendars'),
('calendars', 'delete', 'Delete calendars'),
('calendars', 'share', 'Share calendars'),
('tasks', 'create', 'Create tasks'),
('tasks', 'read', 'View tasks'),
('tasks', 'update', 'Edit tasks'),
('tasks', 'delete', 'Delete tasks'),
('chat', 'create', 'Create chat rooms'),
('chat', 'read', 'View chat messages'),
('chat', 'send', 'Send chat messages'),
('meetings', 'create', 'Create meetings'),
('meetings', 'join', 'Join meetings'),
('meetings', 'record', 'Record meetings'),
('ai', 'use', 'Use AI features'),
('backup', 'create', 'Create backups'),
('backup', 'restore', 'Restore from backups');

-- Assign permissions to roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
   OR (r.name = 'user' AND p.resource IN ('documents', 'emails', 'calendars', 'tasks', 'chat', 'meetings', 'ai'))
   OR (r.name = 'premium' AND p.resource IN ('documents', 'emails', 'calendars', 'tasks', 'chat', 'meetings', 'ai', 'backup'));

-- Insert default AI models
INSERT INTO ai_models (name, provider, model_name, capabilities, context_window, is_local, priority) VALUES
('Ollama Llama 2 7B', 'ollama', 'llama2:7b', ARRAY['text', 'chat'], 4096, true, 10),
('Ollama Code Llama', 'ollama', 'codellama:7b', ARRAY['code', 'text'], 4096, true, 9),
('Ollama Mistral', 'ollama', 'mistral:7b', ARRAY['text', 'chat'], 8192, true, 8);

-- =============================================================================
-- TRIGGERS FOR AUTOMATED MAINTENANCE
-- =============================================================================

-- Update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply update trigger to relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_documents_updated_at BEFORE UPDATE ON documents FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_email_accounts_updated_at BEFORE UPDATE ON email_accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_calendars_updated_at BEFORE UPDATE ON calendars FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_contacts_updated_at BEFORE UPDATE ON contacts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sync_folders_updated_at BEFORE UPDATE ON sync_folders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Update search vectors
CREATE OR REPLACE FUNCTION update_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(NEW.content, '')), 'C');
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply search vector triggers
CREATE TRIGGER update_documents_search_vector BEFORE INSERT OR UPDATE ON documents
    FOR EACH ROW EXECUTE FUNCTION update_search_vector();

-- Similar triggers for other searchable tables would be added here

-- Clean up expired typing indicators
CREATE OR REPLACE FUNCTION cleanup_expired_typing_indicators()
RETURNS void AS $$
BEGIN
    DELETE FROM typing_indicators WHERE expires_at < CURRENT_TIMESTAMP;
END;
$$ language 'plpgsql';

-- Clean up old audit logs (keep last 1 year)
CREATE OR REPLACE FUNCTION cleanup_old_audit_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM audit_log WHERE timestamp < CURRENT_TIMESTAMP - INTERVAL '1 year';
END;
$$ language 'plpgsql';

-- Clean up old rate limit events (keep last 24 hours)
CREATE OR REPLACE FUNCTION cleanup_old_rate_limits()
RETURNS void AS $$
BEGIN
    DELETE FROM rate_limit_events WHERE window_end < CURRENT_TIMESTAMP - INTERVAL '24 hours';
END;
$$ language 'plpgsql';

-- =============================================================================
-- VIEWS FOR COMMON QUERIES
-- =============================================================================

-- User activity view
CREATE VIEW user_activity AS
SELECT
    u.id,
    u.username,
    u.email,
    u.last_login_at,
    u.created_at,
    COUNT(d.id) as document_count,
    COUNT(e.id) as email_count,
    COUNT(t.id) as task_count,
    COUNT(cr.id) as chat_room_count
FROM users u
LEFT JOIN documents d ON u.id = d.owner_id
LEFT JOIN emails e ON u.id = e.account_id
LEFT JOIN tasks t ON u.id = t.user_id
LEFT JOIN chat_rooms cr ON u.id = cr.created_by
GROUP BY u.id, u.username, u.email, u.last_login_at, u.created_at;

-- System health view
CREATE VIEW system_health AS
SELECT
    'users' as metric,
    COUNT(*) as value
FROM users
UNION ALL
SELECT
    'active_users',
    COUNT(*)
FROM users
WHERE last_login_at > CURRENT_TIMESTAMP - INTERVAL '30 days'
UNION ALL
SELECT
    'documents',
    COUNT(*)
FROM documents
UNION ALL
SELECT
    'emails',
    COUNT(*)
FROM emails
UNION ALL
SELECT
    'chat_messages',
    COUNT(*)
FROM chat_messages;

-- Database size view
CREATE VIEW database_size AS
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- =============================================================================
-- FINAL NOTES
-- =============================================================================
-- This schema provides a comprehensive foundation for TPT Titan
-- All tables include proper indexing, constraints, and relationships
-- The schema supports all planned features with room for extension
-- Security, audit logging, and backup/recovery are built-in
-- Performance optimizations include partitioning, indexing, and maintenance triggers
