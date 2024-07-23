use crate::error_handler::CustomError;
use jsonwebtoken::{decode, errors::ErrorKind, DecodingKey, EncodingKey, Header, Validation};
use rocket::{
    http::Status,
    request::{self, FromRequest, Outcome},
    Request,
};
use serde::{Deserialize, Serialize};
use std::env;

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub sub: String,
    pub iat: usize,
    pub exp: usize,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct JwtGuard {
    pub sub: i32,
    pub iat: usize,
    pub exp: usize,
}

#[derive(Serialize, Deserialize)]
pub struct Jwt {
    pub access_token: String,
    pub refresh_token: String,
    pub id: i32,
}

impl Jwt {
    pub fn sign(id: i32) -> Result<Self, CustomError> {
        let mut jwt = Self {
            access_token: "".to_string(),
            refresh_token: "".to_string(),
            id,
        };
        jwt = jwt.encode()?;
        Ok(jwt)
    }
    fn encode(&mut self) -> Result<Self, CustomError> {
        let access_token_claims = Claims::from(self.id, 1);
        let access_token = jsonwebtoken::encode(
            &Header::default(),
            &access_token_claims,
            &EncodingKey::from_secret(env::var("JWT_SECRET")?.as_ref()),
        )?;
        let refresh_token_claims = Claims::from(self.id, 30);
        let refresh_token = jsonwebtoken::encode(
            &Header::default(),
            &refresh_token_claims,
            &EncodingKey::from_secret(env::var("JWT_SECRET")?.as_ref()),
        )?;
        Ok(Self {
            access_token,
            refresh_token,
            id: self.id,
        })
    }
    pub fn decode(token: String) -> Result<Claims, CustomError> {
        let token_data = decode::<Claims>(
            &token,
            &DecodingKey::from_secret(
                env::var("JWT_SECRET")
                    .expect("JWT Token is missing")
                    .as_ref(),
            ),
            &Validation::default(),
        )
        .map_err(|err| match *err.kind() {
            ErrorKind::InvalidToken => CustomError::Unauthorized("Token is invalid".to_string()),
            ErrorKind::InvalidIssuer => CustomError::BadRequest("Issuer is invalid".to_string()),
            _ => CustomError::Internal("Some other error occurred".to_string()),
        })?;
        Ok(token_data.claims)
    }
}

#[derive(Debug)]
pub enum TokenError {
    MissingKey,
    InvalidKey,
}

impl Claims {
    fn from(id: i32, days: usize) -> Self {
        let current_date = chrono::Local::now().timestamp() as usize;
        let exp_date = current_date + 1000 * 60 * 60 * 24 * days;
        Claims {
            sub: id.to_string(),
            iat: current_date,
            exp: exp_date,
        }
    }
}

macro_rules! impl_jwt_guard {
    ($T:ident) => {
        #[rocket::async_trait]
        impl<'r> FromRequest<'r> for $T {
            type Error = TokenError;
            async fn from_request(request: &'r Request<'_>) -> request::Outcome<Self, Self::Error> {
                match request.headers().get_one("Authorization") {
                    Some(s) => match Jwt::decode(s.to_string().replace("Bearer ", "")) {
                        Ok(claims) => match Self::from_claims(claims) {
                            Ok(c) => Outcome::Success(c),
                            Err(_) => {
                                return Outcome::Error((
                                    Status::Unauthorized,
                                    TokenError::InvalidKey,
                                ))
                            }
                        },
                        Err(_) => Outcome::Error((Status::Unauthorized, TokenError::InvalidKey)),
                    },
                    None => Outcome::Forward(Status::Unauthorized),
                }
            }
        }
    };
}

macro_rules! impl_jwt_guard_from_claims {
    ($T:ident) => {
        impl $T {
            fn from_claims(claim: Claims) -> Result<Self, CustomError> {
                let res = Self {
                    sub: claim.sub.parse().unwrap(),
                    iat: claim.iat,
                    exp: claim.exp,
                };
                Ok(res)
            }
        }
    };
}

impl_jwt_guard_from_claims!(JwtGuard);
impl_jwt_guard!(JwtGuard);
