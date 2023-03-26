use rocket_sync_db_pools::{database, diesel};

#[database("diesel")]
pub struct DB(diesel::SqliteConnection);
