// @generated automatically by Diesel CLI.

diesel::table! {
    groups (id) {
        id -> Integer,
        name -> Text,
    }
}

diesel::table! {
    links (id) {
        id -> Integer,
        name -> Text,
        href -> Text,
        group_id -> Integer,
    }
}

diesel::table! {
    users (id) {
        id -> Integer,
        username -> Text,
        bcrypt -> Text,
    }
}

diesel::joinable!(links -> groups (group_id));

diesel::allow_tables_to_appear_in_same_query!(
    groups,
    links,
    users,
);
