/// Create new group.
#[post("/groups")]
pub fn create_group() -> &'static str {
    "New group was created."
}

/// Edit group.
#[put("/groups/<group_id>")]
pub fn edit_group(group_id: u64) -> &'static str {
    "Group was saved."
}

/// Delete group.
#[delete("/groups/<group_id>")]
pub fn delete_group(group_id: u64) -> &'static str {
    "Group was deleted."
}
