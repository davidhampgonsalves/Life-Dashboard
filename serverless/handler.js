'use strict';

const R = require('ramda');
const moment = require('moment-timezone');

const google = require('./lib/googleCalendar');
const ical = require('./lib/ical');
const forecastIO = require('./lib/forecastIO');
const magicseaweed = require('./lib/magicseaweed');
const mobileFoodMarket = require('./lib/mobileFoodMarket');

module.exports.hello = async (event, context, callback) => {
  let events = await Promise.all([
    google.fetchCalendarEvents('davidhampgonsalves@gmail.com'),
    google.fetchCalendarEvents('limbl4hcvdmioc654k3g72pbeg@group.calendar.google.com', ""),
    ical.fetchCalendarEvents("https://recollect.a.ssl.fastly.net/api/places/D23C8C62-A1B4-11E6-8E02-82F09D80A4F0/services/330/events.en.ics", ""),
    mobileFoodMarket.fetchMostRecentEvent(),
  ]);

  const startOfDay = moment.utc().startOf('day');
  const endOfDay = moment.utc().add(1, "days").startOf('day');

  console.log(events);
  events = R.pipe(
		R.flatten,
		R.reject(R.isNil),
		R.filter(e => e.start.isSame(startOfDay, "day")),
    R.sort((a, b) => {
      if(!a.start || !b.start) return Number.MAX_SAFE_INTEGER;
      return a.valueOf() - b.valueOf();
    }),
    R.map(e => {
      if(e.start.isSame(startOfDay) && e.end.isSame(endOfDay))
        return R.omit(["start", "end"], e);

      e.start.tz('America/Halifax');
      e.end.tz('America/Halifax');

      return e;
    }))(events);

  const weather = await forecastIO.fetchForecast();
  const surf = await magicseaweed.fetchForecast();

  const response = {
    statusCode: 200,
    body: JSON.stringify({ events, weather, surf }),
  };

  callback(null, response);
};
