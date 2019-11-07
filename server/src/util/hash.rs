use argonautica::{Hasher, Verifier, Error};
use std::env::var;

pub fn generate_hash(password: &String) -> Result<String, Error> {
    let secret = var("HASH_SECRET").expect("Environment variable HASH_SECRET must be defined");

    Hasher::default().with_secret_key(secret).with_password(password).hash()
}

pub fn verify_hash(hash: &String, password: &String) -> Result<bool, Error> {
    let secret = var("HASH_SECRET").expect("Environment variable HASH_SECRET must be defined");

    Verifier::default().with_secret_key(secret).with_hash(hash).with_password(password).verify()
}
