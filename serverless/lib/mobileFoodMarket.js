'use strict';

const R = require('ramda');
const rp = require('request-promise-native');
const fs = require('fs');
const moment = require('moment-timezone');
const cheerio = require('cheerio');

const title = "North End Mobile Food Market";
module.exports.fetchMostRecentEvent = async () => {
  try {
    const $ = await rp({ uri: "https://www.mobilefoodmarket.ca/calendar/", headers: { 'User-Agent': 'Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36' }, transform: cheerio.load });
		const dateString = $('.eventlist-title-link').parentsUntil('eventlist-column-info').find('.event-date').attr('datetime');

    return {
      faIcon: "",
      title,
			start: moment.tz(dateString, 'America/Halifax'),
			end: moment.tz(dateString, 'America/Halifax').add(1, "days"),
		};
  } catch (e) {
		console.error(e);
    return {};
  }
};
