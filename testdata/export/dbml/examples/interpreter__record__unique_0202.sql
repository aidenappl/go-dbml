Table "user_profiles" {
  "user_id" int
  "profile_type" varchar
  "data" text

  Indexes {
    (user_id, profile_type) [unique]
  }
}

Records user_profiles(user_id, profile_type, data) {
  1, 'work', 'Software Engineer'
  1, 'personal', 'Loves hiking'
  2, 'work', 'Designer'
}
