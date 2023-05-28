CREATE TABLE users
(
    id              UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    gmail           VARCHAR(255) not null UNIQUE,
    username        VARCHAR(30) UNIQUE,
    nickname        VARCHAR(30)  not null DEFAULT '',
    role            VARCHAR(30)  not null DEFAULT 'default',
    is_registered   BOOLEAN      not null DEFAULT false
);

CREATE TABLE user_appid
(
    user_id UUID        not null PRIMARY KEY,
    app_id  UUID UNIQUE not null,
    constraint user_fk foreign key (user_id) references users (id)
);

