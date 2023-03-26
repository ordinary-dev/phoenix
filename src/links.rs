/// Get the list of links.
#[get("/")]
pub fn get_links() -> &'static str {
    "Links."
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
