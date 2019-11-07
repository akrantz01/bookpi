use rocket::response::{self, Responder, Response};
use rocket::request::Request;
use rocket::http::{ContentType, Status};
use rocket_contrib::json::JsonValue;

#[derive(Debug)]
pub struct ApiResponse {
    json: JsonValue,
    status: Status,
}

impl ApiResponse {
    pub fn success() -> ApiResponse {
        ApiResponse {
            status: Status::Ok,
            json: json!({"status": "success"})
        }
    }

    pub fn success_with_data(data: JsonValue) -> ApiResponse {
        ApiResponse {
            status: Status::Ok,
            json: json!({"status": "success", "data": data})
        }
    }

    pub fn error(message: String, status: Status) -> ApiResponse {
        ApiResponse {
            status: status,
            json: json!({"status": "error", "reason": message})
        }
    }
}

impl<'r> Responder<'r> for ApiResponse {
    fn respond_to(self, req: &Request) -> response::Result<'r> {
        Response::build_from(self.json.respond_to(&req).unwrap())
            .status(self.status)
            .header(ContentType::JSON)
            .ok()
    }
}
