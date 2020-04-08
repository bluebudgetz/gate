MATCH (n) DETACH DELETE (n);
CREATE (cmp:Account {id: apoc.create.uuid(), name: 'Big Company', createdOn: datetime()})
CREATE (bank:Account {id: apoc.create.uuid(), name: 'Big Bank', createdOn: datetime()})
CREATE (aig:Account {id: apoc.create.uuid(), name: 'Insurance Company', createdOn: datetime()})
CREATE (bnkAcc:Account {id: apoc.create.uuid(), name: 'Checkings Account', createdOn: datetime()})
CREATE (loans:Account {id: apoc.create.uuid(), name: 'Loans', createdOn: datetime()})-[:childOf]->(bnkAcc)
CREATE (carLn:Account {id: apoc.create.uuid(), name: 'Car loan', createdOn: datetime()})-[:childOf]->(loans)
CREATE (rnvLn:Account {id: apoc.create.uuid(), name: 'Renovations loan', createdOn: datetime()})-[:childOf]->(loans)
CREATE (mrtg:Account {id: apoc.create.uuid(), name: 'Mortgages', createdOn: datetime()})-[:childOf]->(bnkAcc)
CREATE (homeMrtg:Account {id: apoc.create.uuid(), name: 'Home mortgage', createdOn: datetime()})-[:childOf]->(mrtg)
CREATE (offcMrtg:Account {id: apoc.create.uuid(), name: 'Office mortgage', createdOn: datetime()})-[:childOf]->(mrtg)
CREATE (insr:Account {id: apoc.create.uuid(), name: 'Insurances', createdOn: datetime()})-[:childOf]->(bnkAcc)
CREATE (lifeInsr:Account {id: apoc.create.uuid(), name: 'Life insurance', createdOn: datetime()})-[:childOf]->(insr)
CREATE (hlthInsr:Account {id: apoc.create.uuid(), name: 'Health insurance', createdOn: datetime()})-[:childOf]->(insr)
WITH *
UNWIND range(1, 12) AS month
CREATE (cmp)-[:Paid {id: apoc.create.uuid(), amount: 10000, issuedOn: datetime({year: 2019, month: month, day: 9}), comment: 'Salary'}]->(bnkAcc)
CREATE (carLn)-[:Paid {id: apoc.create.uuid(), amount: 250, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Loan installment'}]->(bank)
CREATE (rnvLn)-[:Paid {id: apoc.create.uuid(), amount: 350, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Loan installment'}]->(bank)
CREATE (homeMrtg)-[:Paid {id: apoc.create.uuid(), amount: 2000, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Mortgage installment'}]->(bank)
CREATE (offcMrtg)-[:Paid {id: apoc.create.uuid(), amount: 4000, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Mortgage installment'}]->(bank)
CREATE (lifeInsr)-[:Paid {id: apoc.create.uuid(), amount: 185, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Insurance installment'}]->(aig)
CREATE (hlthInsr)-[:Paid {id: apoc.create.uuid(), amount: 200, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Insurance installment'}]->(aig);

// Roots:
//MATCH p=(parent:Account)<-[c:childOf*0..]-(child:Account)
//  WHERE NOT (parent)-[:childOf]->()
//RETURN parent.name, collect(child.name);

//MATCH (n) DETACH DELETE (n);
//CREATE (c1:Account {id: 'Child001'})
//CREATE (c11:Account {id: 'Child011'})-[:childOf{comment:'c011->c001'}]->(c1)
//CREATE (c12:Account {id: 'Child012'})-[:childOf{comment:'c012->c001'}]->(c1)
//CREATE (c13:Account {id: 'Child013'})-[:childOf{comment:'c013->c001'}]->(c1)
//CREATE (c131:Account {id: 'Child131'})-[:childOf{comment:'c131->c001'}]->(c13)
//CREATE (c2:Account {id: 'Child002'})
//CREATE (c21:Account {id: 'Child021'})-[:childOf{comment:'c021->c002'}]->(c2)
//CREATE (c22:Account {id: 'Child022'})-[:childOf{comment:'c022->c002'}]->(c2)
//
//CREATE (c1)-[:Paid {amount: 1, comment:'c001-to-c002'}]->(c2)
//CREATE (c11)-[:Paid {amount: 11, comment:'c011-to-c021'}]->(c21)
//CREATE (c12)-[:Paid {amount: 12, comment:'c012-to-c022'}]->(c22)
//
//CREATE (c2)-[:Paid {amount: 2, comment:'c002-to-c001'}]->(c1)
//CREATE (c22)-[:Paid {amount: 22, comment:'c022-to-c012'}]->(c12)
