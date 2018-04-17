'use strict';

const R = require('ramda');
const ical = require('ical');
const rp = require('request-promise-native');
const moment = require('moment-timezone');

module.exports.fetchCalendarEvents = async (uri, icon="") => {
  const icalStr = await rp({ uri });
  const data = R.values(ical.parseICS(icalStr));

  return R.reduce((events, e) => R.append({
    faIcon: icon,
    title: e.summary,
    start: moment.tz(e.start, 'America/Halifax').utc(),
    end: (e.end ? moment.tz(e.end, 'America/Halifax') : moment.tz(e.start, 'America/Halifax').add(1, "day")).utc(),
  }, events), [], data);
};
