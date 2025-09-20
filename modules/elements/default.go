package elements

/*
actions to create 

INSERT INTO definitions (name, slug, app_id, active) VALUES
    ("Profile", "profile", 1),
    ("Page", "page", 1);

INSERT INTO definitions_status (name, slug, definition_id) VALUES
    ("Default", "default", 1),
    ("Draft", "draft", 2),
    ("Published", "published", 2);

INSERT INTO definitions_schemas (schema_id, definition_id) VALUES
    (2, 3),
    (3, 2);

INSERT INTO definitions_hierarchy (parent_id, child_id) VALUES
    (NULL, 2),
    (2, 2);

INSERT INTO definition_state_machine (status_id, next_id) VALUES
    (2, 3),
    (3, 2);

*/ 
