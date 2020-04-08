MATCH (:Account)-[tx:Paid {id: $id}]->(:Account)
DELETE (tx)
