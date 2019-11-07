#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate diesel;
#[macro_use] extern crate rocket;
#[macro_use] extern crate rocket_contrib;
#[macro_use] extern crate serde_derive;

extern crate argonautica;
extern crate chrono;
extern crate dotenv;
extern crate serde;
extern crate serde_json;

use dotenv::dotenv;
use rocket_contrib::serve::StaticFiles;
use rocket_contrib::templates::Template;

mod util;

#[database("database")]
pub struct DatabaseConnection(diesel::SqliteConnection);

fn main() {
    dotenv().ok();

    rocket::ignite()
        .mount("/", routes![])
        .mount("/static", StaticFiles::from("./static"))
        .register(catchers![])
        .attach(DatabaseConnection::fairing())
        .attach(Template::fairing())
        .launch();
}
