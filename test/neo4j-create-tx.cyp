MATCH (src {id: $srcId}), (dst {id: $dstId})
MERGE (src)-[:Paid {id: $id, amount: $amount, issuedOn: $issuedOn, comment: $comment}]->(dst)

/*
MERGE (cmp)-[:Paid {id: 'salary_2019' + month + '09', amount: 10000, issuedOn: datetime({year: 2019, month: month, day: 9}), comment: 'Salary'}]->(bnkAcc)
MERGE (carLn)-[:Paid {id: 'carloan_2019' + month + '01', amount: 250, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Loan installment'}]->(bank)
MERGE (rnvLn)-[:Paid {id: 'revenueloan_2019' + month + '01', amount: 350, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Loan installment'}]->(bank)
MERGE (homeMrtg)-[:Paid {id: 'mortgages_home_2019' + month + '01', amount: 2000, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Mortgage installment'}]->(bank)
MERGE (offcMrtg)-[:Paid {id: 'mortgages_office_2019' + month + '01', amount: 4000, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Mortgage installment'}]->(bank)
MERGE (lifeInsr)-[:Paid {id: 'insurances_life_2019' + month + '01', amount: 185, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Insurance installment'}]->(aig)
MERGE (hlthInsr)-[:Paid {id: 'insurances_health_2019' + month + '01', amount: 200, issuedOn: datetime({year: 2019, month: month, day: 1}), comment: 'Insurance installment'}]->(aig)
*/