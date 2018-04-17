extern crate ical;

use chrono::prelude::*;
use std::error::Error;
use std::io::BufReader;
use event::Event;
use reqwest;

pub fn fetch_events(name:&str, url:&str) -> Result<Vec<Event>, String> {
    println!("fetching {} ical", name);
    let res = reqwest::get(url).unwrap();
    let reader = ical::IcalParser::new(BufReader::new(res));

    for line in reader {
        // https://github.com/Peltoche/ical-rs/blob/5a737134167829f0f7eb62920bfca4dd9ee3fbf8/src/parser/ical/component.rs
        println!("{:?}", line.unwrap().events.len());
    }

    return Ok(vec![Event {
        icon: "derby".to_string(),
        description: "BODY".to_string(),
        title: "".to_string(),
        start: Utc::now(),
        end: Utc::now(),
    }]);
}
