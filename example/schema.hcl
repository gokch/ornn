table "user" {
  schema = schema.test
  column "seq" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = true
  }
  column "id" {
    null    = false
    type    = varchar(50)
    default = ""
  }
  column "ord" {
    null = true
    type = bigint
  }
  column "name" {
    null    = true
    type    = varchar(50)
    default = ""
  }
  column "pw" {
    null = true
    type = varbinary(50)
  }
  primary_key {
    columns = [column.seq]
  }
}
schema "test" {
  charset = "utf8mb4"
  collate = "utf8mb4_general_ci"
}
