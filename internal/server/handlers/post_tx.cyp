MATCH (src:Account {id: $sourceAccountId}), (dst:Account {id: $targetAccountId})
CREATE(src)-[tx:Paid {id:       apoc.create.uuid(),
                      issuedOn: $issuedOn,
                      origin:   $origin,
                      amount:   $amount,
                      comment:  $comment}]->(dst)
RETURN src, tx, dst
