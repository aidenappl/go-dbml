Enum "job_status" {
  "created" [note: 'Job created and pending']
  "running" [note: 'Waiting for warehouse to process']
  "done"
  "failed"
  "wait for validation" [note: 'Enum label that has white spaces']
}

Enum "order status" {
  "created" [note: 'Order created']
  "pending"
  "processing"
  "completed"
}

Table "jobs" {
  "id" integer [pk]
  "status" job_status [note: 'This is a column note']
}

Table "orders" {
  "id" int [unique, pk]
  "status" "order status"
  "created_at" varchar
}
