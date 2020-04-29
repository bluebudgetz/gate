// Fetch all accounts, along with a list of direct & indirect children for each node
MATCH (a:Account)<-[:childOf*0..]-(cn:Account)
WITH a, [c IN collect(cn) WHERE c.id <> a.id | c.id] AS children

// Fetch & summarize all outgoing payments for each node
OPTIONAL MATCH p = (a)<-[:childOf*0..]-(:Account)-[r:Paid]->(:Account)
WITH a, children, sum(r.amount) AS outgoing

// Fetch & summarize all incoming payments for each node
OPTIONAL MATCH p = (a)<-[:childOf*0..]-(:Account)<-[r:Paid]-(:Account)
WITH a, children, sum(r.amount) AS incoming, outgoing

// Pull everything together
RETURN a as account, children, toFloat(incoming) AS incoming, toFloat(outgoing) AS outgoing, toFloat(incoming - outgoing) AS balance
  ORDER BY a.name
// TODO: support skip/limit