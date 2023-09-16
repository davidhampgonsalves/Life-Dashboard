"use strict";

const rp = require("request-promise-native");
const cheerio = require("cheerio");

function observationToUnicode(observationType) {
  const ot = observationType.toLowerCase();
  if (ot.match(/snow|freezing|ice|squalls/)) return "❄️";
  if (ot.match(/rain|mist|precipitation|drizzle|thunder/)) return "🌂";
  if (ot.match(/fog|haze/)) return "🌫";
  if (ot.match(/clear/)) return "☀️";
  if (ot.match(/partly cloud/)) return "⛅";
  if (ot.match(/cloud/)) return "☁";
  return "🦞";
}
module.exports.fetchForecast = async () => {
  try {
    const $ = await rp({
      uri: "https://weather.gc.ca/city/pages/ns-19_metric_e.html",
      headers: {
        "User-Agent":
          "Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
      },
      transform: cheerio.load,
    });
    const observationType = $("img[data-v-33b01f9c]").first().attr("alt");
    console.log(observationType, observationToUnicode(observationType));
    let description = $(".pdg-tp-0").first().children("td").last().text();
    description = description.split(".").slice(0, 2).join(".");

    return {
      emoji: observationToUnicode(observationType),
      description,
      temperatureHigh: $(".mrgn-lft-sm[title=High]").first().text(),
      temperatureLow: $(".mrgn-lft-sm[title=Low]").first().text(),
    };
  } catch (e) {
    console.error(e);
    return {};
  }
};
