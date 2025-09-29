package elements

const (
// ### Get All Root Elements (No Parent)
	IndexRootElementsMySQLQuery = `
SELECT * FROM elements
WHERE parent_id IS NULL
AND deleted_at IS NULL;
	`
	// 
	IndexElementChildrenMySQLQuery = `
SELECT * FROM elements WHERE parent_id = ?;
	`



    // TODO: will WITH RECURSIVE work in MySQL?
	// ### Get All Descendants of an Element (Recursive)
	IndexElementDescendantsMySQLQuery = `
WITH RECURSIVE element_tree AS (
    -- Anchor member: select the root element
    SELECT id, parent_id, name, slug, 1 as level
    FROM elements 
    WHERE id = ?
    
    UNION ALL
    
    -- Recursive member: select children
    SELECT e.id, e.parent_id, e.name, e.slug, et.level + 1
    FROM elements e
    INNER JOIN element_tree et ON e.parent_id = et.id
)
SELECT * FROM element_tree ORDER BY level, id;
	`
	//### Count Children for Each Element
	CountElementChildreMySQLQuery = `
SELECT 
    parent.id,
    parent.name,
    parent.slug,
    COUNT(child.id) as child_count
FROM elements parent
LEFT JOIN elements child ON parent.id = child.parent_id
-- WHERE parent.definition_id = ?
GROUP BY parent.id, parent.name, parent.slug
ORDER BY child_count DESC;

	`
	//### Find Orphaned Elements (Parent doesn't exist)
	IndexOrphanedElements = `
SELECT e.* 
FROM elements e
LEFT JOIN elements p ON e.parent_id = p.id
WHERE e.parent_id IS NOT NULL 
-- AND e.definition_id = ?
AND e.parent_id != 0 
AND p.id IS NULL;
	`
	//### Get Complete Subtree with Level Information
	query = `
WITH RECURSIVE subtree AS (
    SELECT 
        id, 
        parent_id, 
        name, 
        slug, 
        level,
        CAST(slug AS CHAR(1000)) as hierarchy_path
    FROM elements 
    WHERE id = [root_element_id]
    
    UNION ALL
    
    SELECT 
        e.id, 
        e.parent_id, 
        e.name, 
        e.slug, 
        e.definition_slug,
        s.level,
        CONCAT(s.hierarchy_path, ' > ', e.slug)
    FROM elements e
    INNER JOIN subtree s ON e.parent_id = s.id
)
SELECT * FROM subtree ORDER BY level, name;
	`

	//### Get Organizational Structure
	query = `

WITH RECURSIVE org_structure AS (
    -- Start from top-level organizational elements
    SELECT 
        id, 
        parent_id, 
        name, 
        slug, 
        definition_slug,
        level,
        CAST(name AS CHAR(255)) as position
    FROM elements 
    WHERE definition_slug = 'organization' 
    AND parent_id IS NULL
    
    UNION ALL
    
    SELECT 
        e.id, 
        e.parent_id, 
        e.name, 
        e.slug, 
        e.definition_slug,
        os.level,
        CONCAT(os.position, ' - ', e.name)
    FROM elements e
    INNER JOIN org_structure os ON e.parent_id = os.id
    WHERE os.level < 10 -- Limit depth
)
SELECT * FROM org_structure ORDER BY level, position;
	`
    