begin transaction isolation level serializable;

alter table device
    add column alerting_enabled     bool default false,
    add column lat                  float     null,
    add column lng                  float     null,
    add column last_alert_timestamp timestamp null;

commit;
