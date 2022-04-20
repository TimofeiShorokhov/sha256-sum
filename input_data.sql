create table shasum
(
    id        bigserial primary key,
    file      text                                                 not null,
    checksum  text                                                 not null,
    file_path text,
    algorithm text                                                 not null,
    deleted   boolean default false                                not null,
    constraint shasum_unique unique (file_path, algorithm)
);

INSERT INTO shasum (file, file_path, checksum, algorithm) VALUES
                                                                    ('1.txt','1/1.txt','123','sha256'),
                                                                    ('1.txt','1/1.txt','1234','sha512'),
                                                                    ('2.txt','1/2.txt','123','sha256'),
                                                                    ('3.txt','1/3.txt','123','sha256'),
                                                                    ('3.txt','1/3.txt','1234','sha512');