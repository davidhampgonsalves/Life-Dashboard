"use strict";

const R = require("ramda");
const moment = require("moment-timezone");
const { google } = require("googleapis");

module.exports.fetchCalendarEvents = async (calendarID) => {
  const client = await google.auth
    .getClient({
      keyFile: "creds/jwt.keys.json",
      scopes: "https://www.googleapis.com/auth/calendar.readonly",
    })
    .catch(console.error);

  const calendar = google.calendar({ version: "v3", auth: client });

  const startOfDay = moment.utc().startOf("day");
  const endOfDay = moment.utc().add(1, "days").startOf("day");

  const events = await calendar.events
    .list({
      calendarId: calendarID,
      timeMin: startOfDay.toISOString(),
      timeMax: endOfDay.toISOString(),
      maxResults: 10,
      singleEvents: true,
      orderBy: "startTime",
    })
    .catch(console.error);

  return R.map((e) => {
    var json = {
      title: e.summary,
      description: e.description,
    };
    if (e.visibility === "private")
      json = { title: "Work 📎", description: "" };

    console.log(e.visibility);
    if (R.has("dateTime", e.start)) {
      json.start = moment(e.start.dateTime).utc();
      json.end = moment(e.end.dateTime).utc();
    } else {
      json.start = moment(e.start.date).utc();
      json.end = moment(e.end.date).utc();
    }

    return json;
  }, events.data.items);
};
