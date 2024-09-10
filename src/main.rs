#[macro_use]
extern crate rocket;
#[macro_use]
extern crate rocket_sync_db_pools;

mod authentication;
mod database;
mod models;
mod schema;

#[launch]
fn rocket() -> _ {
    let routes = vec![authentication::auth_routes()].concat();
    rocket::build().mount("/", routes).attach(database::stage())
}
