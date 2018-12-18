'use strict';

const R = require('ramda');
const rp = require('request-promise-native');
const fs = require('fs');
const moment = require('moment-timezone');
const cheerio = require('cheerio');

module.exports.fetchMostRecentEvent = async () => {
  try {
    const $ = await rp({ uri: "https://www.hrce.ca/about-our-schools/parents/school-cancellations", headers: { 'User-Agent': 'Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36' }, transform: cheerio.load });
    const title = $('#page-title').text().trim();

    if(!title.match(/close/i))
      return null;

    const start = moment.utc().startOf('day');
    const end = moment.utc().add(1, "days").startOf('day');

    return {
      faIcon: "",
      title,
			start,
			end,
		};
  } catch (e) {
		console.error(e);
    return null;
  }
};
