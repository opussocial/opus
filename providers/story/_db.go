package story


GetStoryBySlugMySQLQuery = `
SELECT * FROM stories 
WHERE slug = ?
AND deleted_at IS NULL;
`

GetStorySettingsByIDMySQLQuery = `
SELECT s.* FROM settings s
WHERE s.story_id = (SELECT id FROM stories WHERE id = ?')
ORDER BY s.scope, s.name;
`

GetStorySettingsBySlugMySQLQuery = `
SELECT s.* FROM settings s
WHERE s.story_id = (SELECT id FROM stories WHERE slug = ?)
ORDER BY s.scope, s.name;
`

GetDefinitionSchemasMySQLQuery = `
SELECT sch.* 
FROM schemas sch
JOIN definitions_schemas ds ON sch.id = ds.schema_id
JOIN definitions d ON ds.definition_id = d.id
WHERE d.id = ?
ORDER BY sch.name;
`

GetStoryDefinitionHierarchy = `
SELECT 
    parent.name as parent_name,
    parent.slug as parent_slug,
    child.name as child_name,
    child.slug as child_slug
FROM definition_hierarchy dh
JOIN definitions parent ON dh.parent_id = parent.id
JOIN definitions child ON dh.child_id = child.id
WHERE parent.story_id = (SELECT id FROM stories WHERE slug = 'your-story-slug')
   OR child.story_id = (SELECT id FROM stories WHERE slug = 'your-story-slug')
ORDER BY parent.name, child.name;
`

GetDefinitionStatus = `
SELECT 
    d.name as definition_name,
    ds.name as status_name,
    ds.slug as status_slug
FROM definition_status ds
JOIN definitions d ON ds.definition_id = d.id
WHERE d.id = ?
ORDER BY d.name, ds.name;
`

GetDefinitionStateMachine = `
SELECT 
    d.name as definition_name,
    current_status.name as current_status,
    next_status.name as next_status
FROM definition_state_machine sm
JOIN definition_status current_status ON sm.status_id = current_status.id
JOIN definition_status next_status ON sm.next_id = next_status.id
JOIN definitions d ON current_status.definition_id = d.id
WHERE d.id = ?
ORDER BY d.name, current_status.name;
`

GetStoryGraphRelationships = `
SELECT 
    gr.*,
    source_def.name as source_definition,
    target_def.name as target_definition
FROM graph_relationships gr
JOIN definitions source_def ON gr.source_id = source_def.id
JOIN definitions target_def ON gr.target_id = target_def.id
WHERE gr.story_id = (SELECT id FROM stories WHERE slug = 'your-story-slug')
ORDER BY gr.name;
`
const (
	// -- Get the main story
	IndexStoriesMySQLQuery = `
	SELECT * FROM stories 
WHERE slug = 'your-story-slug' 
AND deleted_at IS NULL;
`
	// -- Get all settings for the story
	IndexSettingsByStoryIDMySQLQuery = `
SELECT s.* FROM settings s
WHERE s.story_id = ?
ORDER BY s.scope, s.name;
`
	// -- Get all definitions for the story
	IndexDefinitionsByStoryIDMySQLQuery = `
SELECT d.* FROM definitions d
WHERE d.story_id = ?
ORDER BY d.scope, d.name;
`
	// -- Get schemas linked to story's definitions
	IndexDefinitionSchemasMySQLQuery = `
SELECT DISTINCT sch.* 
FROM schemas sch
JOIN definitions_schemas ds ON sch.id = ds.schema_id
JOIN definitions d ON ds.definition_id = d.id
WHERE d.story_id = ?
ORDER BY sch.name;
`
	// -- Get hierarchy relationships for story's definitions
	IndexDefinitionHierarchyMySQLQuery = `
SELECT 
    parent.name as parent_name,
    parent.slug as parent_slug,
    child.name as child_name,
    child.slug as child_slug
FROM definition_hierarchy dh
JOIN definitions parent ON dh.parent_id = parent.id
JOIN definitions child ON dh.child_id = child.id
WHERE parent.story_id = ?
   OR child.story_id = ?
ORDER BY parent.name, child.name;
`
	//-- Get all statuses for story's definitions
	IndexDefinitionStatusMySQLQuery = `
SELECT 
    d.name as definition_name,
    ds.name as status_name,
    ds.slug as status_slug
FROM definition_status ds
JOIN definitions d ON ds.definition_id = d.id
WHERE d.story_id = ?
ORDER BY d.name, ds.name;
`

	//-- Get state machine transitions for story's definitions
	IndexDefinitionStateMachineMySQLQuery = `
SELECT 
    d.name as definition_name,
    current_status.name as current_status,
    next_status.name as next_status
FROM definition_state_machine sm
JOIN definition_status current_status ON sm.status_id = current_status.id
JOIN definition_status next_status ON sm.next_id = next_status.id
JOIN definitions d ON current_status.definition_id = d.id
WHERE d.story_id = ?
ORDER BY d.name, current_status.name;
`
	//-- Get all roles for the story
	IndexStoryRolesMySQLQuery = `
SELECT r.* FROM roles r
WHERE r.story_id = ?
ORDER BY r.name;
`
	//-- Get permissions for each role in the story
	IndexRolePermissionsMySQLQuery = `
SELECT 
    r.name as role_name,
    r.slug as role_slug,
    p.name as permission_name,
    p.slug as permission_slug
FROM roles r
JOIN roles_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.story_id = ?
ORDER BY r.name, p.name
`
	//-- Get graph relationships for the story
	IndexStoryGraphRelationshipsMySQLQuery = `
SELECT \
    gr.*,
    source_def.name as source_definition,
    target_def.name as target_definition
FROM graph_relationships gr
JOIN definitions source_def ON gr.source_id = source_def.id
JOIN definitions target_def ON gr.target_id = target_def.id
WHERE gr.story_id = ?
ORDER BY gr.name
`
	// -- Get profile relationships for the story
	MySQLQuery = `
SELECT * FROM profile_relationships
WHERE user_id = ?
ORDER BY created_at
LIMIT 1
`
