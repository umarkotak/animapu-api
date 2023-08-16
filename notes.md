<!-- upload backup -->
scp -i "~/.ssh/default.pem" .env ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app

<!-- create backup -->
pg_dump -d dexpense_development > dexpense.sql

<!-- download backup -->
scp -i "~/.ssh/default.pem" ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/dexpense.sql .

<!-- create db -->
psql -d postgres
create database dexpense_development

<!-- restore backup -->
pg_restore -d dexpense_development dexpense.sql
psql -d dexpense_development < dexpense.sql
