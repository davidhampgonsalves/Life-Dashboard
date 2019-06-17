'use strict';

const R = require('ramda');
const rp = require('request-promise-native');
const fs = require('fs');

// https://darksky.net/dev/docs#forecast-request
// clear-day, clear-night, rain, snow, sleet, wind, fog, cloudy, partly-cloudy-day, or partly-cloudy-night

const emojiToUnicode = {
  "clear-day": "☀️",
  "clear-night": "🌙",
  rain: "🌂",
  snow: "❄️",
  sleet: "❄️",
  fog: "🌫",
  wind: "🎏",
  cloudy: "☁",
  "partly-cloudy-day": "⛅",
  "partly-cloudy-night": "⛅",
};
module.exports.fetchForecast = async () => {
  try {
    const key = fs.readFileSync('creds/forecast.io.txt', "utf8");

    const body = await rp({ uri: `https://api.darksky.net/forecast/${R.trim(key)}/44.652,-63.601?units=auto&exclude=currently,minutely,alerts,flags` });
    const w = JSON.parse(body);

    return {
      emoji: emojiToUnicode[w.hourly.icon],
      temperatureHigh: Math.round(w.daily.data[0].temperatureHigh),
      temperatureLow: Math.round(w.daily.data[0].temperatureLow),
      description: w.hourly.summary,
      weekDescription: w.daily.summary,
    }
  } catch (e) {
    console.error(e);
    return {};
  }
};
