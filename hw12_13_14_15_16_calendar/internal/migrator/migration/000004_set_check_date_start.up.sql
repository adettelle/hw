alter table event
add constraint check_date_start
check (date_end > date_start);