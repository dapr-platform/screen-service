-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE o_big_screen(
                             id VARCHAR(32) NOT NULL,
                             created_by VARCHAR(32) NOT NULL,
                             created_time TIMESTAMP NOT NULL,
                             updated_by VARCHAR(32) NOT NULL,
                             updated_time TIMESTAMP NOT NULL,
                             name VARCHAR(255) NOT NULL,
                             state INTEGER NOT NULL,
                             index_image VARCHAR(255) NOT NULL,
                             remarks VARCHAR(255) NOT NULL,
                             content TEXT NOT NULL,
                             PRIMARY KEY (id)
);

COMMENT ON TABLE o_big_screen IS '大屏';
COMMENT ON COLUMN o_big_screen.id IS '唯一标识';
COMMENT ON COLUMN o_big_screen.created_by IS '创建人';
COMMENT ON COLUMN o_big_screen.created_time IS '创建时间';
COMMENT ON COLUMN o_big_screen.updated_by IS '更新人';
COMMENT ON COLUMN o_big_screen.updated_time IS '更新时间';
COMMENT ON COLUMN o_big_screen.name IS '名称';
COMMENT ON COLUMN o_big_screen.state IS '状态（-1未发布，1发布）';
COMMENT ON COLUMN o_big_screen.index_image IS '首页图片';
COMMENT ON COLUMN o_big_screen.remarks IS '备注';
COMMENT ON COLUMN o_big_screen.content IS '内容';


CREATE TABLE o_group(
                        id VARCHAR(32) NOT NULL,
                        created_by VARCHAR(32),
                        created_time TIMESTAMP,
                        updated_by VARCHAR(32),
                        updated_time TIMESTAMP,
                        name VARCHAR(255),
                        parent_id VARCHAR(255),
                        level INTEGER,
                        sort_num INTEGER,
                        type VARCHAR(255),
                        PRIMARY KEY (id)
);

COMMENT ON TABLE o_group IS '分组';
COMMENT ON COLUMN o_group.id IS '唯一标识';
COMMENT ON COLUMN o_group.created_by IS '创建人';
COMMENT ON COLUMN o_group.created_time IS '创建时间';
COMMENT ON COLUMN o_group.updated_by IS '更新人';
COMMENT ON COLUMN o_group.updated_time IS '更新时间';
COMMENT ON COLUMN o_group.name IS '名称';
COMMENT ON COLUMN o_group.parent_id IS '父id';
COMMENT ON COLUMN o_group.level IS '层次';
COMMENT ON COLUMN o_group.sort_num IS '排序';
COMMENT ON COLUMN o_group.type IS '分类';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS o_big_screen cascade ;
DROP TABLE IF EXISTS o_group cascade ;
-- +goose StatementEnd
