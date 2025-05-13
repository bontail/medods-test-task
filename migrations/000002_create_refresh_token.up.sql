CREATE TABLE RefreshTokens
(
    id           bigserial PRIMARY KEY,
    user_guid    UUID REFERENCES Users (guid) ON DELETE NO ACTION NOT NULL,
    secret_value char(60)                                         NOT NULL,
    created_at   timestamp                                        NOT NULL,
    expires_at   timestamp                                        NOT NULL,
    blocked_at   timestamp,
    user_agent   varchar(1000)                                    NOT NULL, -- стандарт не устанавливается максимальное значение для UA, но нужно иметь хоть какие-нибудь ограничения.
    ip           inet                                             NOT NULL
);