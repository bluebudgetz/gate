// Fetch account graph, starting with $id, and all its direct & indirect children for each node
MATCH (a:Account)<-[:childOf*0..]-(cn:Account)
WHERE a.id = $id
WITH a, [c IN collect(cn) WHERE c.id <> a.id | c.id] AS children

// Fetch & summarize all outgoing payments for each node
OPTIONAL MATCH p = (a)<-[:childOf*0..]-(:Account)-[r:Paid]->(:Account)
WITH a, children, sum(r.amount) AS outgoing

// Fetch & summarize all incoming payments for each node
OPTIONAL MATCH p = (a)<-[:childOf*0..]-(:Account)<-[r:Paid]-(:Account)
WITH a, children, sum(r.amount) AS incoming, outgoing

// Pull everything together
RETURN a as account, children, incoming, outgoing, incoming - outgoing AS balance
  ORDER BY a.name
