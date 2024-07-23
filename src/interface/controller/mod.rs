pub mod authentication;

pub fn get_routes() -> Vec<rocket::Route> {
    let routes = vec![authentication::auth_routes()].concat();
    return routes;
}
