use rocket::fairing::AdHoc;
use rocket_db_pools::diesel::PgPool;
use rocket_db_pools::Database;

#[derive(Database)]
#[database("gambler")]
pub struct Db(PgPool);

pub fn stage() -> AdHoc {
    AdHoc::on_ignite("Database Migrations", |rocket| async {
        rocket.attach(Db::init())
    })
}
