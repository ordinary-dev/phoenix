#[macro_use] extern crate rocket;

/// Get the list of links.
#[get("/")]
fn index() -> &'static str {
    "Links."
}

/// Login page.
#[get("/login")]
fn login() -> &'static str {
    "Please introduce yourself."
}

/// Admin page.
#[get("/admin")]
fn admin() -> &'static str {
    "Admin page."
}

/// Create new group.
#[post("/groups")]
fn create_group() -> &'static str {
    "New group was created."
}

/// Edit group.
#[put("/groups/<group_id>")]
fn edit_group(group_id: u64) -> &'static str {
    "Group was saved."
}

/// Delete group.
#[delete("/groups/<group_id>")]
fn delete_group(group_id: u64) -> &'static str {
    "Group was deleted."
}

/// Create new link.
#[post("/links")]
fn create_link() -> &'static str {
    "New link was created."
}

/// Edit link.
#[put("/links/<link_id>")]
fn edit_link(link_id: u64) -> &'static str {
    "Link was saved."
}

/// Delete link.
#[delete("/links/<link_id>")]
fn delete_link(link_id: u64) -> &'static str {
    "Link was deleted."
}

#[launch]
fn rocket() -> _ {
    rocket::build()
        .mount("/", routes![index, login, admin])
        .mount("/", routes![create_group, edit_group, delete_group])
        .mount("/", routes![create_link, edit_link, delete_link])
}
