begin transaction isolation level serializable;

create table if not exists user_account
(
    id               serial primary key,
    email            varchar(321) unique not null,
    passhash         varchar(255),
    first_name       varchar(64),
    last_name        varchar(64),
    is_active        bool default false,
    activation_token varchar(255)
);
create index if not exists idx_user_account_email on user_account using hash (email);

create table if not exists device
(
    id         serial primary key,
    name       varchar(50) not null,
    created_at timestamp default (now() at time zone 'utc'),
    user_id    int         not null references user_account
);

create table if not exists track_input
(
    id          bigserial primary key,
    lat         float     not null,
    lng         float     not null,
    recorded_at timestamp not null,
    created_at  timestamp default (now() at time zone 'utc'),
    device_id   int       not null references device
);

commit;
