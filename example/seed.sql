set search_path to schema_test;
select set_config('test.user', 'seed_from_sql', true);
insert into main (username,email) values ('From raw SQL','test');
