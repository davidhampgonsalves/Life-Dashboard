'use strict';

const R = require('ramda');
const rp = require('request-promise-native');
const fs = require('fs');
const moment = require('moment-timezone');

const filterByToday = R.filter(f => {
  const timestamp = moment.unix(f.timestamp).tz('America/Halifax')
  return timestamp.isSame(moment().tz('America/Halifax'), "day");
});

module.exports.fetchForecast = async () => {
  try {
    const key = R.trim(fs.readFileSync('creds/magicseaweed.txt', "utf8"));

    const body = await rp({ uri: `http://magicseaweed.com/api/${key}/forecast/?spot_id=787` });
    const w = JSON.parse(body);
    const maxStars = R.pipe(
			filterByToday,
			R.pluck('solidRating'),
			ratings => Math.max(...ratings))(w);

    return { maxRating: maxStars };
  } catch (e) {
		console.error(e);
    return {};
  }
};
