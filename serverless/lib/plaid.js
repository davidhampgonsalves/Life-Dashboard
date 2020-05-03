'use strict';

const R = require('ramda');
const fs = require('fs');
const plaid = require('plaid');
const moment = require('moment');

const totalDebitsByDay = (transactions, dateStr) => R.pipe(
  R.filter(R.propEq("date", dateStr)),
  R.pluck("amount"),
  R.filter(v => v > 0),
  R.sum)(transactions);

const round = f => Math.round(f * 100) / 100;

module.exports.fetchFinance = async () => {
  try {
    const { client_id, secret, public_key, access_token, account_id } = JSON.parse(fs.readFileSync('creds/plaid.json', "utf8"));

    // README - to create new account use the plaid-link.html and watch the console when you are done adding the account, then
    // use this code to generate an access_token which will not expire for future use.
    // then you can make a request to get transactions without accounts to determine the account ids
    //
    //const plaidClient = new plaid.Client(client_id, secret, public_key, plaid.environments.development, {version: '2018-05-22'});
    //let res = await plaidClient.exchangePublicToken("public-development-04ce89ad-bbf2-411c-9bba-c3ecc78216e8");
    //console.log("TOKEN", res.access_token);
    //const todayDateStr = moment().tz('America/Halifax').format("YYYY-MM-DD");
    //const yesterdayDateStr = moment().tz('America/Halifax').subtract(1, "days").format("YYYY-MM-DD");
    //res = await plaidClient.getTransactions(res.access_token, yesterdayDateStr, todayDateStr);
    //console.log(res);

    const plaidClient = new plaid.Client(client_id, secret, public_key, plaid.environments.development, {version: '2018-05-22'});

    const todayDateStr = moment().tz('America/Halifax').format("YYYY-MM-DD");
    const yesterdayDateStr = moment().tz('America/Halifax').subtract(1, "days").format("YYYY-MM-DD");
    const res = await plaidClient.getTransactions(access_token, yesterdayDateStr, todayDateStr, { account_ids: [account_id] });

    const { transactions } = res;
    const todayTotalDebits = totalDebitsByDay(transactions, todayDateStr);
    const yesterdayTotalDebits = totalDebitsByDay(transactions, yesterdayDateStr);

    return { todayTotalDebits: round(todayTotalDebits), yesterdayTotalDebits: round(yesterdayTotalDebits) };
  } catch (e) {
    console.error(e);
    return {};
  }
};
