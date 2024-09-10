use crate::database::Db;
use rocket::form::Form;
use rocket_db_pools::Connection;

pub fn auth_routes() -> Vec<rocket::Route> {
    routes![login]
}

#[derive(FromForm)]
struct Login<'r> {
    r#username: &'r str,
    r#password: &'r str,
}

#[post("/login", data = "<login>")]
fn login(db: Connection<Db>, login: Form<Login<'_>>) {
    format!("Hello, {}!", login.r#username);
}
