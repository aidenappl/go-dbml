Table "employees" {
  "employee_id" number [pk, increment]
  "first_name" nvarchar2(50) [not null, default: 'Unknown']
  "last_name" nvarchar2(50) [not null]
  "middle_name" varchar2(50) [default: '']
  "department" varchar2(30) [unique, not null, default: 'General']
  "years_of_service" integer [not null, default: 0]
  "employee_rank" number(2) [unique, not null, default: 1]
  "hire_date" date [not null, default: `SYSDATE`]
  "last_review_date" timestamp [default: `current_timestamp`]
  "termination_date" date
  "birth_date" date [not null]
  "is_active" number(1) [not null, default: 1]
  "is_manager" number(1) [default: 0]
  "notes" clob
  "profile_data" nclob [default: '{}']
}

Table "projects" {
  "project_id" int [pk, increment]
  "project_code" varchar2(20) [unique, not null]
  "project_name" nvarchar2(100) [not null, default: 'Unnamed Project']
  "start_date" date [not null, default: `SYSDATE`]
  "end_date" date
  "status" varchar2(20) [not null, default: 'PLANNED']
  "priority" integer [unique, default: 999]
}

Table "constraint_orderings" {
  "field1" varchar2(50) [unique, not null, default: 'test1']
  "field2" varchar2(50) [unique, not null, default: 'test2']
  "field3" varchar2(50) [unique, not null, default: 'test3']
  "field4" varchar2(50) [unique, not null, default: 'test4']
  "field5" varchar2(50) [not null, default: 'test5']
  "field6" varchar2(50) [unique, default: 'test6']
  "field7" int [pk, not null, default: 1]
}
