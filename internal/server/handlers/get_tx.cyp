MATCH (src:Account)-[tx:Paid {id: $id}]->(dst:Account)
RETURN tx, src, dst
