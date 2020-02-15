create table if not exists device
(
    id         serial primary key,
    name       varchar(50) not null,
    created_at timestamp default (now() at time zone 'utc')
);

create table if not exists track_input
(
    id          serial primary key,
    lat         float8    not null,
    lng         float8    not null,
    recorded_at timestamp not null,
    created_at  timestamp default (now() at time zone 'utc'),
    device_id   int       not null references device
);
