CREATE TABLE short_url_record_0 (
    unique_id BIGINT UNSIGNED NOT NULL PRIMARY KEY COMMENT '生成的唯一id',
    short_url VARCHAR(1024) NOT NULL COMMENT '短url',
    long_url VARCHAR(1024) NOT NULL COMMENT '长url',
    view_cnt INT NOT NULL default 0 COMMENT '访问次数',
    ctime DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    mtime DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) ENGINE=Innodb DEFAULT CHARSET='utf8' COMMENT '短链记录表';

CREATE TABLE long_url_record_0 (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    long_url_md5 INT UNSIGNED NOT NULL COMMENT '长url的md5',
    long_url VARCHAR(1024) NOT NULL COMMENT '长url',
    unique_id BIGINT UNSIGNED NOT NULL COMMENT '生成的唯一id',
    ctime DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    mtime DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    INDEX idx_long_url_md5 (long_url_md5)
) ENGINE=Innodb DEFAULT CHARSET='utf8' COMMENT '长链记录表';