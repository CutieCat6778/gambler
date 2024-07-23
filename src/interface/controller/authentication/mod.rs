use crate::{
    database::Db,
    handler::{error_handler::CustomError, result_handler::CustomResult},
    interface::service::authentication::jwt::JwtGuard,
};
use rocket::serde::json::Json;
use rocket_db_pools::Connection;
use serde::{Deserialize, Serialize};

pub fn auth_routes() -> Vec<rocket::Route> {
    routes![login]
}

#[derive(Serialize, Deserialize)]
struct Login {
    username: String,
    password: String,
}

#[post("/login", data = "<login>")]
fn login(_db: Connection<Db>, login: Json<Login>) {
    format!("Hello, {}!", login.r#username);
}

#[get("/logout")]
fn logout(_db: Connection<Db>, _claim: JwtGuard) -> Result<CustomResult, CustomError> {
    Ok(CustomResult {
        result: "Logout".to_string(),
    })
}
