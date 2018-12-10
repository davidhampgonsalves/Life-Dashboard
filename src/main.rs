extern crate chrono;
extern crate reqwest;
extern crate imageproc;
extern crate image;
extern crate openssl_probe;
extern crate serde;
extern crate serde_json;
extern crate rusttype;

#[macro_use]
extern crate serde_derive;

use std::panic;
use chrono::prelude::*;
use serde_json::Error;
use std::sync::mpsc;
use std::fs::File;
use std::path::Path;
use std::io::BufWriter;
use std::thread;
use image::{Luma, GrayImage};
use imageproc::rect::Rect;
use imageproc::drawing::{
    draw_cross_mut,
    draw_line_segment_mut,
    draw_hollow_rect_mut,
    draw_filled_rect_mut,
    draw_hollow_circle_mut,
    draw_filled_circle_mut,
    draw_text_mut,
};
use rusttype::{Font, FontCollection, Scale, point};
use chrono::{DateTime, FixedOffset};

#[derive(Debug, Deserialize)]
struct Event {
    faIcon: String,
    title: String,
    description: Option<String>,
    start: Option<DateTime::<FixedOffset>>,
    end: Option<DateTime::<FixedOffset>>,
}

#[derive(Debug, Deserialize)]
struct Weather {
    faIcon: String,
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
    now: DateTime::<Utc>,
}

const LINE_PADDING:u32 = 1;
const PARAGRAPH_PADDING:u32 = 10;
const ICON_SIZE:u32 = 50;
const MARGIN:u32 = 20;
const WIDTH:u32 = 600;
const CONTENT_WIDTH:u32 = WIDTH - (2 * MARGIN);
const EVENT_TIME_WIDTH:u32 = 80;
const EVENT_CONTENT_MARGIN:u32 = MARGIN + ICON_SIZE + MARGIN;
const EVENT_CONTENT_WIDTH:u32 = WIDTH - EVENT_CONTENT_MARGIN - MARGIN;
const HEIGHT:u32 = 800;
const WEATHER_OFFSET:u32 = HEIGHT - 200;
const FINANCE_OFFSET:u32 = HEIGHT - 300;

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

fn draw_text_block(image: &mut GrayImage, color: Luma<u8>, font: &Font, scale: Scale, text: &str, width: u32, x: u32, y: u32) -> u32 {
    let mut lines: Vec<String> = vec!["".to_string()];

    for word in text.split(" ") {
        let line_width = calculate_glyph_width(font, scale, &format!("{} {}", lines.last().unwrap(), word));

        if(line_width > width) {
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

fn main() {
    openssl_probe::init_ssl_cert_env_vars();

    let serif_font = Font::from_bytes(include_bytes!("../fonts/Bookerly-Regular.ttf") as &[u8]).expect("Error constructing Font");
    let icon_font = Font::from_bytes(include_bytes!("../fonts/fa-solid-900.ttf") as &[u8]).expect("Error constructing Font");

    let background_color = Luma([255u8]);
    let text_color_1 = Luma([0u8]);
    let icon_color = Luma([100u8]);

    let mut image = GrayImage::new(600, 800);
    draw_filled_rect_mut(&mut image, Rect::at(0, 0).of_size(600, 800), background_color);

    let scale = Scale::uniform(30.0);
    let small_scale = Scale::uniform(15.0);
    let large_scale = Scale::uniform(50.0);
    let icon_scale = Scale::uniform(ICON_SIZE as f32);

    let mut response = match reqwest::get("https://blakwkb41l.execute-api.us-east-1.amazonaws.com/dev/summary") {
        Ok(res) => res,
        Err(e) => {
            println!("error: {:?}", e);
            draw_text_mut(&mut image, icon_color, (WIDTH / 2 - 20), 200, icon_scale, &icon_font, &"".to_string());
            draw_text_mut(&mut image, icon_color, 220, 260, scale, &serif_font, &"Error pulling data.".to_string());

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
            draw_text_mut(&mut image, icon_color, (WIDTH / 2 - 20), 200, icon_scale, &icon_font, &"".to_string());
            draw_text_mut(&mut image, icon_color, 190, 260, scale, &serif_font, &"Error transforming data.".to_string());

            let path = Path::new("image.png");
            let file = File::create(path).unwrap();

            let fout = &mut BufWriter::new(file);
            let mut encoder = image::png::PNGEncoder::new(fout);
            encoder.encode(&image, 600, 800, image::Gray(8));
            return;
        },
    };

    // Draw date
    draw_text_mut(&mut image, text_color_1, WIDTH - MARGIN - (large_scale.x * 3.0) as u32, MARGIN, large_scale, &serif_font, &data.now.format("%b %e").to_string());

    let mut offset = MARGIN + large_scale.x as u32 + LINE_PADDING;
    for event in data.events {
        let mut time_offset:u32 = 0;

        if let Some(start) = event.start {
            time_offset = EVENT_TIME_WIDTH;
            draw_text_mut(&mut image, text_color_1, MARGIN, offset, scale, &serif_font, &start.format("%H:%M").to_string());
            draw_text_mut(&mut image, text_color_1, MARGIN, offset + scale.y as u32 + LINE_PADDING, scale, &serif_font, &event.end.unwrap().format("%H:%M").to_string());
        }
        draw_text_mut(&mut image, icon_color, MARGIN + time_offset, offset, icon_scale, &icon_font, &event.faIcon);

        offset += draw_text_block(&mut image, text_color_1, &serif_font, scale, &event.title, EVENT_CONTENT_WIDTH - time_offset, EVENT_CONTENT_MARGIN + time_offset, offset);
        offset += draw_text_block(&mut image, text_color_1, &serif_font, scale, &event.description.unwrap_or("".to_string()), EVENT_CONTENT_WIDTH - time_offset, EVENT_CONTENT_MARGIN + time_offset, offset);

        offset += PARAGRAPH_PADDING;
    }

    // Draw surf
    draw_text_mut(&mut image, icon_color, MARGIN, offset, icon_scale, &icon_font, &"".to_string());
    offset += draw_text_block(&mut image, text_color_1, &serif_font, scale, &format!("{}ft at {} seconds,", data.surf.height, data.surf.period).to_string(), EVENT_CONTENT_WIDTH, EVENT_CONTENT_MARGIN, offset);
    offset += draw_text_block(&mut image, text_color_1, &serif_font, scale, &format!("{}-{} stars.", data.surf.maxRating, data.surf.fadedRating).to_string(), EVENT_CONTENT_WIDTH, EVENT_CONTENT_MARGIN, offset);

    // Draw finance
    draw_text_mut(&mut image, icon_color, MARGIN, FINANCE_OFFSET + LINE_PADDING, icon_scale, &icon_font, &"".to_string());
    let finance_offset = draw_text_block(&mut image, text_color_1, &serif_font, scale, &format!("${} spent today, ${} yesterday.", data.finance.todayTotalDebits, data.finance.yesterdayTotalDebits), EVENT_CONTENT_WIDTH, EVENT_CONTENT_MARGIN, FINANCE_OFFSET);

    // Draw weather
    draw_text_mut(&mut image, icon_color, MARGIN, WEATHER_OFFSET + LINE_PADDING, icon_scale, &icon_font, &data.weather.faIcon);
    let weather_offset = draw_text_block(&mut image, text_color_1, &serif_font, scale, &format!("{} - {}°C. {}", data.weather.temperatureLow, data.weather.temperatureHigh, data.weather.description), EVENT_CONTENT_WIDTH, EVENT_CONTENT_MARGIN, WEATHER_OFFSET);
    draw_text_block(&mut image, text_color_1, &serif_font, scale, &data.weather.weekDescription, EVENT_CONTENT_WIDTH, EVENT_CONTENT_MARGIN, WEATHER_OFFSET + weather_offset + PARAGRAPH_PADDING);

    // Draw render time
    draw_text_mut(&mut image, text_color_1, WIDTH - MARGIN - MARGIN, HEIGHT - MARGIN, small_scale, &serif_font, &data.now.format("%H:%M").to_string());

    let path = Path::new("image.png");
    let file = File::create(path).unwrap();

    let fout = &mut BufWriter::new(file);
    let mut encoder = image::png::PNGEncoder::new(fout);
    encoder.encode(&image, 600, 800, image::Gray(8));
}

