'use strict';

const R = require('ramda');
const moment = require('moment-timezone');
const { google } = require('googleapis');

module.exports.fetchCalendarEvents = async (calendarID, icon="") => {
  const client = await google.auth
    .getClient({
      keyFile: 'creds/jwt.keys.json',
      scopes: 'https://www.googleapis.com/auth/calendar.readonly',
    })
    .catch(err => {
      console.log(err);
    });

  const calendar = google.calendar({ version: 'v3', auth: client });

  const startOfDay = moment.utc().startOf('day');
  const endOfDay = moment.utc().add(1, "days").startOf('day');

  const events = await calendar.events
    .list({
      calendarId: calendarID,
      timeMin: startOfDay.toISOString(),
      timeMax: endOfDay.toISOString(),
      maxResults: 10,
      singleEvents: true,
      orderBy: 'startTime',
    })
    .catch(err => {
      console.log(err);
    });

  return R.map(
    e => {
      const json = {
        faIcon: icon,
        title: e.summary,
        description: e.description,
      };

      if(R.has("dateTime", e.start)) {
        json.start = moment(e.start.dateTime).utc();
        json.end = moment(e.end.dateTime).utc();
      } else {
        json.start = moment(e.start.date).utc();
        json.end = moment(e.end.date).utc();
      }

      return json;
    },
    events.data.items
  );
};
