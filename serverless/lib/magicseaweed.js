"use strict";

const rp = require("request-promise-native");
const cheerio = require("cheerio");

module.exports.fetchForecast = async () => {
  try {
    const $ = await rp({
      uri: "https://www.ndbc.noaa.gov/station_page.php?station=44258",
      headers: {
        "User-Agent":
          "Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
      },
      transform: cheerio.load,
    });

    return {
      height: $($(".currentobs td")[8]).text(),
      period: $($(".currentobs td")[10]).text(),
    };
  } catch (e) {
    console.error(e);
    return {};
  }
};
