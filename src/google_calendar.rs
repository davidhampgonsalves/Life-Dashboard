use hyper;
use hyper_tls;
use hyper_tls::HttpsConnector;
use hyper::Client;
use calendar3::Channel;
use calendar3::{Result, Error};
use std::default::Default;
use calendar3::CalendarHub;
use chrono::prelude::*;
//use std::error::Error;
use std::result;
use event::Event;
use oauth;
//use reqwest;

pub fn fetch_events(name: &str) -> result::Result<Vec<Event>, String> {
    //println!("fetching {} calendar", name);
    //let body = match reqwest::get("https://www.rust-lang.org") {
        //Ok(mut res) => res.text(),
        //Err(e) => return Err(e.to_string()),
    //}.unwrap();
    //println!("{:?}", body);

    //let body = reqwest::get("https://www.rust-lang.org") {
        //Ok(body) => Ok(body),
        //Err(err) => Err(::error::from(err)),
    //};

    //println!("body = {:?}", text);

    ////print!("{:?}", json);


    let client_secret = oauth::service_account_key_from_file(&"creds/jwt.keys.json".to_string())
        .unwrap();
    let client = hyper::Client::with_connector(HttpsConnector::new(hyper_rustls::TlsClient::new()));
    let mut access = oauth::ServiceAccountAccess::new(client_secret, client);

    use oauth::GetToken;
    println!("{:?}",
             access.token(&vec!["https://www.googleapis.com/auth/pubsub"]).unwrap());


    let client = Client::builder().build::<_, hyper::Body>(https);
    let hub = CalendarHub::new(client, auth);
    let req = Channel::default();
    let result = hub.events().watch(req, "calendarId")
        .updated_min("eirmod")
        .time_zone("elitr")
        .time_min("amet")
        .time_max("no")
        .sync_token("labore")
        .single_events(true)
        .show_hidden_invitations(true)
        .show_deleted(true)
        .add_shared_extended_property("aliquyam")
        .q("accusam")
        .add_private_extended_property("Lorem")
        .page_token("sea")
        .order_by("et")
        .max_results(-70)
        .max_attendees(-21)
        .i_cal_uid("eirmod")
        .always_include_email(false)
        .doit();

    match result {
        Err(e) => match e {
            // The Error enum provides details about what exactly happened.
            // You can also just use its `Debug`, `Display` or `Error` traits
            Error::HttpError(_)
                |Error::MissingAPIKey
                |Error::MissingToken(_)
                |Error::Cancelled
                |Error::UploadSizeLimitExceeded(_, _)
                |Error::Failure(_)
                |Error::BadRequest(_)
                |Error::FieldClash(_)
                |Error::JsonDecodeError(_, _) => println!("{}", e),
        },
        Ok(res) => println!("Success: {:?}", res),
    }

    return Ok(vec![Event {
        icon: "google calendar".to_string(),
        description: "body".to_string(),
        title: "".to_string(),
        start: Utc::now(),
        end: Utc::now(),
    }]);
}
