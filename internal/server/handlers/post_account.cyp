CREATE (account:Account {id: apoc.create.uuid(), name: $name, createdOn: datetime()})
WITH account
CALL apoc.do.when(
$parentId <> null,
'MERGE (parent:Account {id:parentId}) MERGE (child)-[:childOf]->(parent) RETURN child',
'RETURN child',
{child: account, parentId: $parentId}
)
YIELD value
RETURN value AS account
