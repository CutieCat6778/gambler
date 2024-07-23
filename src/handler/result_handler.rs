use rocket::http::ContentType;
use rocket::http::Status;
use rocket::request::Request;
use rocket::response::{self, Responder, Response};
use serde::Serialize;
use std::io::Cursor;

#[derive(Serialize)]
pub struct CustomResult {
    pub result: String,
}

impl std::fmt::Display for CustomResult {
    fn fmt(&self, fmt: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(fmt, "{}", self.result)
    }
}

impl<'r> Responder<'r, 'static> for CustomResult {
    fn respond_to(self, _: &'r Request<'_>) -> response::Result<'static> {
        // serialize struct into json string
        Response::build()
            .status(Status::Ok)
            .header(ContentType::JSON)
            .sized_body(self.result.len(), Cursor::new(self.result))
            .ok()
    }
}
