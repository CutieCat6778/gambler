#[macro_use]
extern crate rocket;
extern crate rocket_sync_db_pools;

mod database;
mod handler;
mod interface;
mod models;
mod schema;

#[launch]
fn rocket() -> _ {
    rocket::build()
        .mount("/", interface::get_routes())
        .attach(database::stage())
}
