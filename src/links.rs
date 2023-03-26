use diesel::prelude::*;
use crate::db::DB;

#[derive(Queryable)]
#[derive(Debug)]
pub struct Link {
    pub id: i32,
    pub name: String,
    pub href: String,
    pub group_id: i32,
}

impl Link {
    pub fn get_all(connection: &mut SqliteConnection) -> Vec<Link> {
        use crate::schema::links::dsl::*;
        return links
            .load::<Link>(connection)
            .expect("Error loading links")
    }
}

/// Get the list of links.
#[get("/")]
pub async fn get_links(db: DB) -> String {
    let links: Vec<Link> = db.run(move |conn| {
        Link::get_all(conn)
    }).await;
    
    format!("Links: {:#?}", links)
}

/// Create new link.
#[post("/links")]
pub fn create_link() -> &'static str {
    "New link was created."
}

/// Edit link.
#[put("/links/<link_id>")]
pub fn edit_link(link_id: u64) -> &'static str {
    "Link was saved."
}

/// Delete link.
#[delete("/links/<link_id>")]
pub fn delete_link(link_id: u64) -> &'static str {
    "Link was deleted."
}
