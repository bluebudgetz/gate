if (!db) {
    var db = db || {};
}
if (!db.accounts) db.accounts = {
    insert: function (o) {
        return {insertedId: 0};
    },
    insertOne: function (o) {
        return {insertedId: 0};
    }
};
if (!db.transactions) db.transactions = {
    insert: function (o) {
    }
};
var d = new Date();

// ACCOUNTS
var a = {createdOn: d, updatedOn: null};
var company = db.accounts.insertOne(Object.assign({}, a, {name: "Big Company", parentId: null})).insertedId;
var acmeBank = db.accounts.insertOne(Object.assign({}, a, {name: "A.C.M.E Bank", parentId: null})).insertedId;
var aig = db.accounts.insertOne(Object.assign({}, a, {name: "A.I.G", parentId: null})).insertedId;
var bankAccount = db.accounts.insertOne(Object.assign({}, a, {name: "Bank Account", parentId: null})).insertedId;
var loans = db.accounts.insertOne(Object.assign({}, a, {name: "Loans", parentId: bankAccount})).insertedId;
var carLoan = db.accounts.insertOne(Object.assign({}, a, {name: "Car Loan", parentId: loans})).insertedId;
var mortgages = db.accounts.insertOne(Object.assign({}, a, {name: "Mortgages", parentId: bankAccount})).insertedId;
var homeMortgage = db.accounts.insertOne(Object.assign({}, a, {name: "Home Mortgage", parentId: mortgages})).insertedId;
var officeMortgage = db.accounts.insertOne(Object.assign({}, a, {
    name: "Office Mortgage",
    parentId: mortgages
})).insertedId;
var insurances = db.accounts.insertOne(Object.assign({}, a, {name: "Insurances", parentId: bankAccount})).insertedId;
var lifeInsurance = db.accounts.insertOne(Object.assign({}, a, {
    name: "Life Insurance",
    parentId: insurances
})).insertedId;
var healthInsurance = db.accounts.insertOne(Object.assign({}, a, {
    name: "Health Insurance",
    parentId: insurances
})).insertedId;

// TRANSACTIONS
var t = {origin: "Initialization", createdOn: d, issuedOn: d, updatedOn: null};
db.transactions.insert([
    Object.assign({}, t, {source: company, target: bankAccount, amount: 4990, comment: "April Salary"}),
    Object.assign({}, t, {source: company, target: bankAccount, amount: 5000, comment: "May Salary"}),
    Object.assign({}, t, {source: carLoan, target: acmeBank, amount: 590, comment: "April car loan"}),
    Object.assign({}, t, {source: carLoan, target: acmeBank, amount: 589, comment: "May car loan"}),
    Object.assign({}, t, {source: homeMortgage, target: acmeBank, amount: 410, comment: "February home mortgage"}),
    Object.assign({}, t, {source: officeMortgage, target: acmeBank, amount: 890, comment: "June office mortgage"}),
    Object.assign({}, t, {source: lifeInsurance, target: aig, amount: 199, comment: "March life insurance"}),
    Object.assign({}, t, {source: healthInsurance, target: aig, amount: 98, comment: "January health insurance"})
]);
