MATCH (src:Account)-[tx:Paid]->(dst:Account)
RETURN tx, src, dst
  ORDER BY tx.issuedOn DESC
  SKIP $skip
  LIMIT $limit
