MATCH (src:Account)-[tx:Paid {id: $id}]->(dst:Account)
SET tx.issuedOn = $issuedOn
SET tx.origin = $origin
SET tx.amount = $amount
SET tx.comment = $comment
RETURN src, tx, dst
