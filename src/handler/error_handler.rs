use jsonwebtoken::errors::Error as JwtError;
use rocket::http::ContentType;
use rocket::http::Status;
use rocket::request::Request;
use rocket::response::{self, Responder, Response};
use rocket::serde::json::{serde_json, serde_json::Error as SerdeError};
use serde::Serialize;
use std::env::VarError;
use std::io::Cursor;

#[derive(Serialize)]
pub struct ErrorResponse {
    pub message: String,
    pub status: u16,
}

//#[derive(Error)]
#[derive(Debug, Clone)]
pub enum CustomError {
    //#[resp("{0}")]
    Internal(String),

    //#[resp("{0}")]
    NotFound(String),

    //#[resp("{0}")]
    BadRequest(String),

    //#[resp("{0}")]
    Unauthorized(String),
}

impl CustomError {
    fn get_http_status(&self) -> (Status, std::string::String) {
        match self {
            CustomError::Internal(msg) => (Status::InternalServerError, msg.clone()),
            CustomError::NotFound(msg) => (Status::NotFound, msg.clone()),
            CustomError::Unauthorized(msg) => (Status::Unauthorized, msg.clone()),
            _ => (Status::BadRequest, "Unkown error".to_string()),
        }
    }
}

impl std::fmt::Display for CustomError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        let (_, msg) = self.get_http_status();
        write!(fmt, "{}", msg)
    }
}

impl From<SerdeError> for CustomError {
    fn from(err: SerdeError) -> Self {
        CustomError::Internal(format!("Json Serde: {}", err.to_string()))
    }
}

impl From<VarError> for CustomError {
    fn from(err: VarError) -> Self {
        CustomError::Internal(format!("Env Var: {}", err.to_string()))
    }
}

impl From<JwtError> for CustomError {
    fn from(error: JwtError) -> CustomError {
        match error {
            error => CustomError::Unauthorized(format!("JWT error {}", error)),
        }
    }
}

impl<'r> Responder<'r, 'static> for CustomError {
    fn respond_to(self, _: &'r Request<'_>) -> response::Result<'static> {
        // serialize struct into json string
        let err_response = serde_json::to_string(&ErrorResponse {
            message: format!("{}", self.to_string()),
            status: self.get_http_status().0.code,
        })
        .unwrap();

        let (status, _) = self.get_http_status();

        Response::build()
            .status(status)
            .header(ContentType::JSON)
            .sized_body(err_response.len(), Cursor::new(err_response))
            .ok()
    }
}
