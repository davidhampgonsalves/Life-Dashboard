extern crate chrono;
extern crate reqwest;
extern crate imageproc;
extern crate image;
extern crate openssl_probe;
extern crate serde;
extern crate serde_json;
extern crate rusttype;
extern crate rand;

#[macro_use]
extern crate serde_derive;

use rand::seq::SliceRandom;
use std::fs::{self, DirEntry};
use chrono::prelude::*;
use std::fs::File;
use std::path::Path;
use std::io::BufWriter;
use image::{Luma, LumaA, GrayAlphaImage};
use imageproc::rect::Rect;
use imageproc::drawing::{
    draw_filled_rect_mut,
    draw_text_mut,
};
use rusttype::{Font, Scale, point};
use chrono::{DateTime, FixedOffset};

#[derive(Debug, Deserialize)]
struct Event {
    faIcon: String,
    title: String,
    description: Option<String>,
    start: Option<DateTime<FixedOffset>>,
    end: Option<DateTime<FixedOffset>>,
}

#[derive(Debug, Deserialize)]
struct Weather {
    emoji: String,
    temperatureHigh: i32,
    temperatureLow: i32,
    description: String,
    weekDescription: String,
}

#[derive(Debug, Deserialize)]
struct Surf {
    maxRating: u8,
    fadedRating: u8,
    period: u8,
    height: f32,
}

#[derive(Debug, Deserialize)]
struct Finance {
    todayTotalDebits: f32,
    yesterdayTotalDebits: f32,
}

#[derive(Debug, Deserialize)]
struct Data {
    surf: Surf,
    weather: Weather,
    events: Vec<Event>,
    finance: Finance,
    now: DateTime<Utc>,
}

const LINE_PADDING:u32 = 0;
const PARAGRAPH_PADDING:u32 = 20;
const MARGIN:u32 = 20;
const WIDTH:u32 = 600;
const HEIGHT:u32 = 800;

fn calculate_glyph_width(font: &Font, scale: Scale, text: &str) -> u32 {
    let glyphs: Vec<_> = font
        .layout(text, scale, point(0.0, 0.0))
        .collect();

    let max_x = glyphs
        .into_iter()
        .map(|g| g.pixel_bounding_box().unwrap_or(rusttype::Rect { min: rusttype::Point { x: 0, y: 0 }, max: rusttype::Point { x: 0, y: 0 } }).max.x )
        .max()
        .unwrap();

    max_x as u32
}

fn draw_text_block(image: &mut GrayAlphaImage, color: LumaA<u8>, font: &Font, scale: Scale, text: &str, width: u32, x: u32, y: u32) -> u32 {
    let mut lines: Vec<String> = vec!["".to_string()];

    for word in text.split(" ") {
        let line_width = calculate_glyph_width(font, scale, &format!("{} {}", lines.last().unwrap(), word));

        if line_width > width {
            lines.push("".to_string());
        }

        let len = lines.len()-1;
        lines[len] = format!("{}{} ", lines.last().unwrap(), word);
    }

    let v_metrics = font.v_metrics(scale);
    let height =(v_metrics.ascent - v_metrics.descent).ceil() as u32 + LINE_PADDING;

    for (i, line) in lines.iter().enumerate() {
        draw_text_mut(image, color, x, y + (i as u32 * height), scale, &font, &line.trim_right());
    }

    height * lines.len() as u32
}

fn format(date: DateTime<FixedOffset>) -> String {
    if date.minute() == 0 {
        date.format("%-H").to_string()
    } else {
        date.format("%-H:%M").to_string()
    }
}

trait ParagraghDrawer {
    fn paragraph(&mut self, text: &str);
}

struct Draw {
    offset: u32,
    image: GrayAlphaImage,
}

impl ParagraghDrawer for Draw {
    fn paragraph(&mut self, text: &str) {
        let scale = Scale::uniform(40.0);
        let font = Font::from_bytes(include_bytes!("../fonts/OpenSans-Regular-AndroidEmoji.ttf") as &[u8]).expect("Error constructing Font");
        let color = LumaA([0u8, 100]);

        self.offset += draw_text_block(&mut self.image, color, &font, scale, text, WIDTH - (2 * MARGIN), MARGIN, self.offset) + PARAGRAPH_PADDING;
    }
}

fn main() {
    openssl_probe::init_ssl_cert_env_vars();

    let font = Font::from_bytes(include_bytes!("../fonts/OpenSans-Regular-AndroidEmoji.ttf") as &[u8]).expect("Error constructing Font");

    let background_color = LumaA([255u8, 255u8]);
    let color = LumaA([0u8, 255u8]);

    let mut image = GrayAlphaImage::new(600, 800);
    draw_filled_rect_mut(&mut image, Rect::at(0, 0).of_size(600, 800), background_color);

    let scale = Scale::uniform(38.0);
    let mut response = match reqwest::get("https://blakwkb41l.execute-api.us-east-1.amazonaws.com/dev/summary") {
        Ok(res) => res,
        Err(e) => {
            println!("error: {:?}", e);
            draw_text_mut(&mut image, color, 220, 260, scale, &font, &"Error pulling data.".to_string());

            let path = Path::new("image.png");
            let file = File::create(path).unwrap();

            let fout = &mut BufWriter::new(file);
            let mut encoder = image::png::PNGEncoder::new(fout);
            encoder.encode(&image, 600, 800, image::Gray(8));

            return;
        },
    };

    let data = match response.json::<Data>() {
        Ok(data) => data,
        Err(e) => {
            println!("error: {:?}", e);
            draw_text_mut(&mut image, color, 190, 260, scale, &font, &"Error transforming data.".to_string());

            let path = Path::new("image.png");
            let file = File::create(path).unwrap();

            let fout = &mut BufWriter::new(file);
            let mut encoder = image::png::PNGEncoder::new(fout);
            encoder.encode(&image, 600, 800, image::Gray(8));
            return;
        },
    };

    // Draw date
    let large_scale = Scale::uniform(50.0);
    let date_str = &data.now.format("%b %e").to_string();
    let date_margin = WIDTH - MARGIN - calculate_glyph_width(&font, large_scale, date_str);
    draw_text_mut(&mut image, color, date_margin, MARGIN, large_scale, &font, date_str);

    let mut draw = Draw {
        offset: 80,
        image: image,
    };

    for event in data.events {
        if let Some(start) = event.start {
            draw.paragraph(&format!("📅 {} 🕘{} - {}.", &event.title, format(start), format(event.end.unwrap())));
        } else {
            draw.paragraph(&format!("📅 {}", &event.title));
        }
    }

    if data.surf.maxRating > 0 {
        draw.paragraph(&format!("🌊 {}{} {} ft at {} secs.", "▪".repeat(data.surf.maxRating as usize), "▫".repeat(5 - data.surf.maxRating as usize), data.surf.height, data.surf.period).to_string());
    }

    if data.finance.todayTotalDebits + data.finance.yesterdayTotalDebits == 0 as f32 {
        draw.paragraph(&format!("💲 👏 Zero spent recently. 👏"));
    } else if data.finance.todayTotalDebits > 0 as f32 {
        draw.paragraph(&format!("💲 {} today, {} yesterday. 🏭", data.finance.todayTotalDebits, data.finance.yesterdayTotalDebits));
    } else {
        draw.paragraph(&format!("💲 {} yesterday. 🏭", data.finance.yesterdayTotalDebits));
    }

    draw.paragraph(&format!("{} {} - {}°C. {}", &data.weather.emoji, data.weather.temperatureLow, data.weather.temperatureHigh, data.weather.description));
    draw.paragraph(&data.weather.weekDescription);

    // pokemon
    let remaining_height = HEIGHT - draw.offset - (2 * MARGIN);
    let max_width = (WIDTH / 2) - (2 * MARGIN);
    if remaining_height > 100 {
        let paths: Vec<_> = fs::read_dir(&Path::new("pokemon/front/")).unwrap().map(|maybe_path| maybe_path.unwrap().path()).collect();
        let path = paths.choose(&mut rand::thread_rng()).unwrap();
        let file_name = path.to_str().unwrap().replace("pokemon/front/", "");

        let mut front_img = image::open(&Path::new(&format!("pokemon/front/{}", file_name))).ok().expect("Opening front image failed");
        let mut back_img = image::open(&Path::new(&format!("pokemon/back/{}", file_name))).ok().expect("Opening back image failed");

        front_img = front_img.resize(max_width, remaining_height, image::imageops::FilterType::Nearest);
        back_img = back_img.resize(max_width, remaining_height, image::imageops::FilterType::Nearest);

        image::imageops::overlay(&mut draw.image, &back_img.to_luma_alpha(), MARGIN, draw.offset + MARGIN);
        image::imageops::overlay(&mut draw.image, &front_img.to_luma_alpha(), (WIDTH / 2) + MARGIN, draw.offset + MARGIN);
    }

    let small_scale = Scale::uniform(15.0);
    draw_text_mut(&mut draw.image, color, WIDTH - MARGIN - MARGIN, HEIGHT - MARGIN, small_scale, &font, &data.now.format("%H:%M").to_string());

    let path = Path::new("image.png");
    let file = File::create(path).unwrap();

    let fout = &mut BufWriter::new(file);
    let encoder = image::png::PNGEncoder::new(fout);
    encoder.encode(&draw.image, 600, 800, image::GrayA(8));
}
