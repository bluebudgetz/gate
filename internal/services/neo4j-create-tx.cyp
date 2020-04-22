MATCH (src {id: $srcId}), (dst {id: $dstId})
MERGE (src)-[:Paid {
  id:        $id,
  amount:    $amount,
  issuedOn:  $issuedOn,
  comment:   $comment,
  createdOn: datetime(),
  origin:    'qa'}
]->(dst)
