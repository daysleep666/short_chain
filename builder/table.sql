CREATE TABLE short_url_record_0 (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    long_url VARCHAR(1024) NOT NULL COMMENT '长url',
    short_url VARCHAR(1024) NOT NULL COMMENT '短链',
    unique_id BIGINT NOT NULL COMMENT '生成的唯一id',
    ctime DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    mtime DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) ENGINE=Innodb DEFAULT CHARSET='utf8' COMMENT '短链记录表';

ALTER TABLE short_url_record_0 ADD UNIQUE uni_long_url (long_url);
ALTER TABLE short_url_record_0 ADD UNIQUE uni_short_url (short_url);