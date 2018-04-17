use chrono::prelude::*;

#[derive(Debug)]
pub struct Event {
    pub icon: String,
    pub title: String,
    pub description: String,
    pub start: DateTime<Utc>,
    pub end: DateTime<Utc>,
}
