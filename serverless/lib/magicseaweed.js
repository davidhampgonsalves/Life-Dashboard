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

    const forecasts = await Promise.all([
      rp({ uri: `http://magicseaweed.com/api/${key}/forecast/?spot_id=787` }),
      rp({ uri: `http://magicseaweed.com/api/${key}/forecast/?spot_id=787` }),
    ]);

    return R.pipe(
      R.map(JSON.parse),
      R.flatten,
			filterByToday,
			f => ({
          maxRating: Math.max(...R.pluck('solidRating', f)),
          fadedRating: Math.max(...R.pluck('fadedRating', f)),
          height: R.max(...R.map(R.path(["swell", "components", "combined", "height"]), f)),
          period: R.max(...R.map(R.path(["swell", "components", "combined", "period"]), f)),
        }))(forecasts);
  } catch (e) {
		console.error(e);
    return {};
  }
};
