#[macro_use] extern crate rocket;

mod groups;
mod links;
mod login;
mod admin;
mod db;
mod schema;

#[launch]
fn rocket() -> _ {
    rocket::build()
        .attach(db::DB::fairing())
        .mount("/", routes![login::login, admin::admin])
        .mount("/", routes![groups::create_group, groups::edit_group, groups::delete_group])
        .mount("/", routes![links::get_links, links::create_link, links::edit_link, links::delete_link])
}
