delete from "_migrations";
drop index if exists appointment_resource_idx;
drop index if exists patient_resource_idx;
drop index if exists practitioner_txid_idx;
drop table if exists migration_test;