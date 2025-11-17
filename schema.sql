DROP TABLE IF EXISTS join_ban;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS join_members;
DROP TABLE IF EXISTS want_go;
DROP TABLE IF EXISTS recruits;
DROP TABLE IF EXISTS worlds;
DROP TABLE IF EXISTS users;

-- ユーザ管理テーブル
CREATE TABLE users (
    id VARCHAR(50),
    name VARCHAR(50) NOT NULL,
    icon VARCHAR(200) DEFAULT '',
    bio TEXT,
    created_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- ワールド管理テーブル
CREATE TABLE worlds (
    id VARCHAR(50),
    name VARCHAR(500) NOT NULL,
    thumbnail VARCHAR(200) DEFAULT '',
    created_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- 募集管理テーブル
CREATE TABLE recruits (
    id VARCHAR(50) PRIMARY KEY,
    content TEXT NOT NULL,
    closed BOOLEAN DEFAULT FALSE,
    user_id VARCHAR(50),
    world_id VARCHAR(50),
    created_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY fk_user_id(user_id) REFERENCES users(id),
    FOREIGN KEY fk_world_id(world_id) REFERENCES worlds(id)
);

-- いきたいワールド管理テーブル
-- ユーザIDとワールドIDで複合主キーにする？
CREATE TABLE want_go (
    user_id VARCHAR(50), -- 外部キー設定する？
    world_id VARCHAR(50), -- 外部キー設定する？
    created_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, world_id),
    FOREIGN KEY fk_user_id(user_id) REFERENCES users(id),
    FOREIGN KEY fk_world_id(world_id) REFERENCES worlds(id)
);

-- 募集に参加しているユーザ管理テーブル
CREATE TABLE join_members (
    recruit_id VARCHAR(50), -- recruitsのIDがいらない場合どうする？
    user_id VARCHAR(50),
    created_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (recruit_id, user_id),
    FOREIGN KEY fk_recruit_id(recruit_id) REFERENCES recruits(id),
    FOREIGN KEY fk_user_id(user_id) REFERENCES users(id)
);

-- 募集内のメッセージ管理テーブル
CREATE TABLE messages (
    id VARCHAR(50),
    recruit_id VARCHAR(50),
    user_id VARCHAR(50),
    content TEXT NOT NULL,
    created_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY fk_recruit_id(recruit_id) REFERENCES recruits(id)
);

-- 募集に参加できないようにBANする管理テーブル
CREATE TABLE join_ban (
    recruit_id VARCHAR(50),
    user_id VARCHAR(50),
    created_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (recruit_id, user_id),
    FOREIGN KEY fk_recruit_id(recruit_id) REFERENCES recruits(id),
    FOREIGN KEY fk_user_id(user_id) REFERENCES users(id)
);