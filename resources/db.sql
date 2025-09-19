-- drop database opus_test; create database opus_test; use opus_test;
-- AUTH MODULE

CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid_identifier CHAR(36) DEFAULT "",
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(64) NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    confirmed_at DATETIME DEFAULT NULL,
    last_login_at DATETIME DEFAULT NULL,
    
    deleted_at DATETIME DEFAULT NULL
);

CREATE TABLE tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED REFERENCES users(id),

    token VARCHAR(255) UNIQUE NOT NULL,
    
    -- scope defines token use
    -- 'auth', 'refresh', 'reset', 'confirm'
    scope VARCHAR(16) NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL
);

CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);


-- STORY MODULE

CREATE TABLE stories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid_identifier CHAR(36) DEFAULT "",
    user_id BIGINT UNSIGNED DEFAULT NULL REFERENCES users(id),

    name VARCHAR(64) NOT NULL,
    slug VARCHAR(64) UNIQUE NOT NULL,
    description VARCHAR(255) DEFAULT "",

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    published_at DATETIME DEFAULT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
);

CREATE TABLE settings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    -- scope defines if user_id and story_id are null or not
    -- 'user', 'story', 'user:story'
    scope VARCHAR(16) NOT NULL,

    user_id BIGINT UNSIGNED DEFAULT NULL REFERENCES users(id),
    story_id BIGINT UNSIGNED DEFAULT NULL REFERENCES stories(id),

    -- TODO: change name to key or add key field
    name VARCHAR(64) UNIQUE NOT NULL,
    value VARCHAR(255) DEFAULT NULL, -- this is the only mutable field

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO settings (scope, name, value) VALUES 
    ("system", "installed", NULL),
    ("system", "home", NULL),
    ("system", "email", NULL),
    ("system", "max_levels", NULL),
    ("system", "max_definitions", NULL),  
    ("system", "max_elements", NULL);

CREATE TABLE definitions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    story_id BIGINT UNSIGNED REFERENCES stories(id),

    name VARCHAR(64) NOT NULL,
    slug VARCHAR(64) NOT NULL,
    description VARCHAR(255) NOT NULL,

    icon VARCHAR(32) NOT NULL DEFAULT "",
    color CHAR(7) DEFAULT "#ffffff",

    has_interactions TINYINT DEFAULT 0,
    has_taxonomy TINYINT DEFAULT 0,
    belongs_to_graph TINYINT DEFAULT 0,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_story_slug(story_id, slug)
);

CREATE TABLE content_schemas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    schema_name VARCHAR(64) NOT NULL,
    UNIQUE KEY schema_name(schema_name)
);

INSERT INTO content_schemas (name, schema_name) VALUES
    ("Text", "schema_text"),
    ("Resource", "schema_resources"),
    ("Article", "schema_article"),
    ("Address", "schema_address"),
    ("Contacts", "schema_contacts"),
    ("Event", "schema_events"),
    ("Files", "schema_files");

CREATE TABLE definitions_schemas (
    schema_id BIGINT UNSIGNED REFERENCES content_schemas(id),
    definition_id BIGINT UNSIGNED REFERENCES definitions(id),
    PRIMARY KEY (schema_id, definition_id)
);

CREATE TABLE definition_hierarchy (
    parent_id BIGINT UNSIGNED REFERENCES definitions(id),
    child_id BIGINT UNSIGNED REFERENCES definitions(id),
    PRIMARY KEY (parent_id, child_id)
);

CREATE TABLE definition_status (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    definition_id BIGINT UNSIGNED REFERENCES definitions(id),

    name VARCHAR(64) NOT NULL,
    slug VARCHAR(64) NOT NULL,

    icon VARCHAR(32) NOT NULL DEFAULT "",
    color CHAR(7) DEFAULT "#ffffff",

    UNIQUE KEY unique_definition_status (definition_id, slug)
);

CREATE TABLE definition_state_machine (
    status_id BIGINT UNSIGNED REFERENCES definition_status(id),
    next_id BIGINT UNSIGNED REFERENCES definition_status(id),
    PRIMARY KEY (status_id, next_id)
);

CREATE TABLE state_machine_triggers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    next_id BIGINT UNSIGNED REFERENCES definition_state_machine(next_id),
    action VARCHAR(32)
);

CREATE TABLE roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    story_id BIGINT UNSIGNED REFERENCES stories(id),
    slug VARCHAR(64) NOT NULL,
    name VARCHAR(64) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_role (story_id, slug)
);

CREATE TABLE permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    -- should be a pattern of story:type:action OR a registry index
    action VARCHAR(64) UNIQUE NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE roles_permissions (
    role_id BIGINT UNSIGNED REFERENCES roles(id),
    permission_id BIGINT UNSIGNED REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE graph_relationships (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    slug VARCHAR(64) NOT NULL,
    name VARCHAR(64) NOT NULL,
    description VARCHAR(255) NOT NULL,

    source_id BIGINT UNSIGNED REFERENCES definitions(id),
    target_id BIGINT UNSIGNED REFERENCES definitions(id),

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_graph_relationship (slug)
);

CREATE TABLE graph_triggers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    graph_relationship_id BIGINT UNSIGNED REFERENCES graph_relationships(id),
    action VARCHAR(32),
    scope VARCHAR(8) -- source or target 
);

-- ELEMENT MODULE

CREATE TABLE elements (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    uuid_identifier CHAR(36) DEFAULT "",
    identifier VARCHAR(128) DEFAULT "",
    main_entity VARCHAR(64) DEFAULT "",
    same_as VARCHAR(128) DEFAULT "",
    url VARCHAR(128) DEFAULT "",

    slug VARCHAR(255) DEFAULT "",
    name VARCHAR(128) DEFAULT "",
    alternate_name VARCHAR(128) DEFAULT "",
    description VARCHAR(255) DEFAULT "",

    level TINYINT DEFAULT 0,
    active TINYINT(1) DEFAULT 0,
    locked TINYINT(1) DEFAULT 0,

    parent_id BIGINT UNSIGNED,
    created_by_profile_id BIGINT UNSIGNED,
    definition_id BIGINT UNSIGNED REFERENCES definitions(id),
    status_id BIGINT UNSIGNED REFERENCES definition_status(id),

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    published_at DATETIME DEFAULT NULL,
    expires_at DATETIME DEFAULT NULL,
    deleted_at DATETIME DEFAULT NULL
);

CREATE TABLE profile_relationships (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED REFERENCES users(id),
    profile_element_id BIGINT UNSIGNED REFERENCES elements(id),
    role_id BIGINT UNSIGNED REFERENCES roles(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    element_id BIGINT UNSIGNED DEFAULT NULL REFERENCES elements(id),
    created_by_profile_id BIGINT UNSIGNED REFERENCES elements(id),

    -- scope defines if comment and rating are null or not
    -- like, fav, comment, rating
    scope VARCHAR(255) NOT NULL DEFAULT "",

    comment VARCHAR(255) NOT NULL DEFAULT "",
    rating TINYINT NOT NULL DEFAULT 0,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE taxonomy (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- future proofing for internal use like campaign relationships in marketing module 
    -- element, definition slug or custom, like 'campaign'
    scope VARCHAR(255) NOT NULL DEFAULT "",
    
    term VARCHAR(128) NOT NULL DEFAULT "",
    slug VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE UNIQUE INDEX idx_taxonomy_slug ON taxonomy(slug);

CREATE TABLE element_taxonomy (
    element_id BIGINT UNSIGNED REFERENCES elements(id),
    taxonomy_id BIGINT UNSIGNED REFERENCES taxonomy(id),
    PRIMARY KEY (element_id, taxonomy_id)
);

CREATE TABLE graph (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    graph_relationship_id BIGINT UNSIGNED REFERENCES graph_relationships(id),
    element_source_id BIGINT UNSIGNED REFERENCES elements(id),
    element_target_id BIGINT UNSIGNED REFERENCES elements(id),
    UNIQUE KEY unique_graph_node (graph_relationship_id, element_source_id, element_target_id)
);

-- MKT MODULE

CREATE TABLE leads (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT DEFAULT NULL REFERENCES users(id),
    email VARCHAR(255) DEFAULT "",
    phone VARCHAR(255) DEFAULT "",
    region VARCHAR(255) DEFAULT "",
    source VARCHAR(64) DEFAULT "",
    cookie VARCHAR(255) NOT NULL,
    signup_date DATETIME DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_cookie ON leads(cookie);

CREATE TABLE operations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64), 
    slug VARCHAR(64), 
    -- scope defines operation type
    -- "email", "ad"
    scope VARCHAR(16),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    operation_id BIGINT REFERENCES operations(id),
    -- scope defines which user segment was targeted
    -- lead, user, profile
    scope VARCHAR(16),
    executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE visits (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    operation_id BIGINT REFERENCES operations(id),
    lead_id BIGINT REFERENCES leads(id),
    
    page_url VARCHAR(255) DEFAULT NULL,
    utm_source VARCHAR(64) DEFAULT NULL,
    utm_medium VARCHAR(64) DEFAULT NULL,
    utm_campaign VARCHAR(64) DEFAULT NULL,
    
    referrer_url VARCHAR(255) DEFAULT NULL,
    user_agent TEXT DEFAULT NULL,
    ip_address VARCHAR(32) DEFAULT NULL,

    visited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE responses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    lead_id BIGINT REFERENCES leads(id),
    operation_id BIGINT REFERENCES operations(id),
    scope VARCHAR(16), -- "email_open", "ad_click", "conversion"
    responded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payoffs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    operation_id BIGINT REFERENCES operations(id),
    points INTEGER DEFAULT 0,
    targets INTEGER DEFAULT 0,
    calculated_at DATETIME DEFAULT CURRENT_TIMESTAMP -- When the rating was calculated
);
