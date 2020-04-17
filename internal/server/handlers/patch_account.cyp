MERGE (n:Account {id: $id})
  ON CREATE SET name = coalesce($name, n.name), createdOn = datetime()
  ON MATCH SET name = coalesce($name, n.name), updatedOn = datetime()
WITH n
CALL apoc.do.when(
$parentId <> null,
'MERGE (p:Account {id:parentId}) MERGE (a)-[:childOf]->(p) RETURN a',
'MATCH (a)-[r:childOf]->() DELETE (r) RETURN a',
{a: n, parentId: $parentId}
)
YIELD value
RETURN value AS account
