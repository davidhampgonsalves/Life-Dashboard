'use strict';

const R = require('ramda');
const rp = require('request-promise-native');
const fs = require('fs');
const moment = require('moment-timezone');
const cheerio = require('cheerio');

module.exports.fetchMostRecentEvent = async () => {
  try {
    const $ = await rp({ uri: "https://www.halifax.ca/transportation/winter-operations/service-updates", headers: { 'User-Agent': 'Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36' }, transform: cheerio.load });
    const title = $("#tablefield-paragraph-126041-field_table-0").find('.row_1.col_1.c-table__cell').text();

    if(title.match(/not be in effect/i))
      return null;

    const start = moment.utc().startOf('day');
    const end = moment.utc().add(1, "days").startOf('day');

    return {
      icon: "🚗",
      title,
			start,
			end,
		};
  } catch (e) {
		console.error(e);
    return null;
  }
};
