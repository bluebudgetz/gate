MATCH (src:Account)-[tx:Paid {id: $id}]->(dst:Account)
SET tx.issuedOn = coalesce($issuedOn, tx.issuedOn)
SET tx.origin = coalesce($origin, tx.origin)
SET tx.amount = coalesce($amount, tx.amount)
SET tx.comment = coalesce($comment, tx.comment)
RETURN src, tx, dst
