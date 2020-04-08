MERGE (account:Account {id: $id})
  ON CREATE SET name = $name, createdOn = datetime()
  ON MATCH SET name = $name, updatedOn = datetime()
WITH account
CALL apoc.do.when(
$parentId <> null,
'MERGE (parent:Account {id:parentId}) MERGE (child)-[:childOf]->(parent) RETURN child',
'MATCH (child)-[r:childOf]->() DELETE (r) RETURN child',
{child: account, parentId: $parentId}
)
YIELD value
RETURN value as account
