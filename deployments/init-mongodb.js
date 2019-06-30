db = db || {accounts: {}, transactions: {}, users: {}};

function account(_id, createdOn, name, parentId) {
    return {_id, createdOn, updatedOn: null, name, parentId};
}

function tx(date, source, target, amount, comments) {
    return {
        origin: "Initialization",
        createdOn: date,
        updatedOn: null,
        issuedOn: date,
        source: source,
        target: target,
        amount: amount,
        comments: comments
    };
}

// ACCOUNTS
const BIG_COMPANY = ObjectId();
const ACME_BANK = ObjectId();
const AIG = ObjectId();
const BANK_ACCOUNT = ObjectId();
const LOANS = ObjectId();
const CAR_LOAN = ObjectId();
const MORTGAGES = ObjectId();
const HOME_MORTGAGE = ObjectId();
const OFFICE_MORTGAGE = ObjectId();
const INSURANCES = ObjectId();
const LIFE_INSURANCE = ObjectId();
const HEALTH_INSURANCE = ObjectId();
db.accounts.insert([
    account(BIG_COMPANY, new Date("2019-01-01T00:00:00Z"), "Big Company", null),
    account(ACME_BANK, new Date("2019-01-01T00:00:00Z"), "A.C.M.E Bank", null),
    account(AIG, new Date("2019-01-01T00:00:00Z"), "A.I.G", null),
    account(BANK_ACCOUNT, new Date("2019-01-01T00:00:00Z"), "Bank Account", null),
    account(LOANS, new Date("2019-02-01T00:00:00Z"), "Loans", BANK_ACCOUNT),
    account(CAR_LOAN, new Date("2019-02-01T00:00:00Z"), "Car Loan", LOANS),
    account(MORTGAGES, new Date("2019-03-01T00:00:00Z"), "Mortgages", BANK_ACCOUNT),
    account(HOME_MORTGAGE, new Date("2019-03-01T00:00:00Z"), "Home Mortgage", MORTGAGES),
    account(OFFICE_MORTGAGE, new Date("2019-05-01T00:00:00Z"), "Big Company", MORTGAGES),
    account(INSURANCES, new Date("2019-02-01T00:00:00Z"), "Big Company", BANK_ACCOUNT),
    account(LIFE_INSURANCE, new Date("2019-02-01T00:00:00Z"), "Big Company", INSURANCES),
    account(HEALTH_INSURANCE, new Date("2019-04-01T00:00:00Z"), "Big Company", INSURANCES),
]);

// TRANSACTIONS
db.transactions.insert([
    tx(new Date("2019-02-01T00:00:00Z"), BIG_COMPANY, BANK_ACCOUNT, 24700, "Salary"),
    tx(new Date("2019-03-01T00:00:00Z"), BIG_COMPANY, BANK_ACCOUNT, 24600, "Salary"),
    tx(new Date("2019-04-01T00:00:00Z"), BIG_COMPANY, BANK_ACCOUNT, 24800, "Salary"),
    tx(new Date("2019-05-01T00:00:00Z"), BIG_COMPANY, BANK_ACCOUNT, 24400, "Salary"),
    tx(new Date("2019-06-01T00:00:00Z"), BIG_COMPANY, BANK_ACCOUNT, 24200, "Salary"),
    tx(new Date("2019-02-01T00:00:00Z"), CAR_LOAN, ACME_BANK, 1289, "Car loan"),
    tx(new Date("2019-03-01T00:00:00Z"), CAR_LOAN, ACME_BANK, 1293, "Car loan"),
    tx(new Date("2019-04-01T00:00:00Z"), CAR_LOAN, ACME_BANK, 1290, "Car loan"),
    tx(new Date("2019-05-01T00:00:00Z"), CAR_LOAN, ACME_BANK, 1295, "Car loan"),
    tx(new Date("2019-06-01T00:00:00Z"), CAR_LOAN, ACME_BANK, 1285, "Car loan"),
    tx(new Date("2019-03-01T00:00:00Z"), HOME_MORTGAGE, ACME_BANK, 5980, "Home mortgage"),
    tx(new Date("2019-04-01T00:00:00Z"), HOME_MORTGAGE, ACME_BANK, 5975, "Home mortgage"),
    tx(new Date("2019-05-01T00:00:00Z"), HOME_MORTGAGE, ACME_BANK, 5985, "Home mortgage"),
    tx(new Date("2019-06-01T00:00:00Z"), HOME_MORTGAGE, ACME_BANK, 5993, "Home mortgage"),
    tx(new Date("2019-05-01T00:00:00Z"), OFFICE_MORTGAGE, ACME_BANK, 1390, "Office mortgage"),
    tx(new Date("2019-06-01T00:00:00Z"), OFFICE_MORTGAGE, ACME_BANK, 1390, "Office mortgage"),
    tx(new Date("2019-02-01T00:00:00Z"), LIFE_INSURANCE, AIG, 199, "Life insurance"),
    tx(new Date("2019-03-01T00:00:00Z"), LIFE_INSURANCE, AIG, 199, "Life insurance"),
    tx(new Date("2019-04-01T00:00:00Z"), LIFE_INSURANCE, AIG, 199, "Life insurance"),
    tx(new Date("2019-05-01T00:00:00Z"), LIFE_INSURANCE, AIG, 199, "Life insurance"),
    tx(new Date("2019-06-01T00:00:00Z"), LIFE_INSURANCE, AIG, 199, "Life insurance"),
    tx(new Date("2019-04-01T00:00:00Z"), HEALTH_INSURANCE, AIG, 89, "Life insurance"),
    tx(new Date("2019-05-01T00:00:00Z"), HEALTH_INSURANCE, AIG, 92, "Life insurance"),
    tx(new Date("2019-06-01T00:00:00Z"), HEALTH_INSURANCE, AIG, 92, "Life insurance"),
]);

// USERS
db.users.insert([
    {_id: "arikkfir", password: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92", name: "Arik Kfir"} // pwd:123456
]);
