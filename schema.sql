DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS worlds;
DROP TABLE IF EXISTS recruits;
DROP TABLE IF EXISTS want_go;
DROP TABLE IF EXISTS join_members;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS join_ban;

-- ユーザ管理テーブル
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    icon VARCHAR(200) DEFAULT '',
    bio TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ワールド管理テーブル
CREATE TABLE worlds (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    thumbnail VARCHAR(200) DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 募集管理テーブル
CREATE TABLE recruits (
    id VARCHAR(50) PRIMARY KEY, -- いらない？　ユーザIDとワールドIDで複合主キーにする？ それとも文字連結した値にする？
    content TEXT NOT NULL,
    closed BOOLEAN DEFAULT FALSE,
    user_id VARCHAR(50), -- 外部キー設定する？
    world_id VARCHAR(50), -- 外部キー設定する？
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- いきたいワールド管理テーブル
-- ユーザIDとワールドIDで複合主キーにする？
CREATE TABLE want_go (
    user_id VARCHAR(50), -- 外部キー設定する？
    world_id VARCHAR(50), -- 外部キー設定する？
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, world_id)
);

-- 募集に参加しているユーザ管理テーブル
CREATE TABLE join_members (
    recruit_id VARCHAR(50), -- recruitsのIDがいらない場合どうする？
    user_id VARCHAR(50),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (recruit_id, user_id)
);

-- 募集内のメッセージ管理テーブル
CREATE TABLE messages (
    id VARCHAR(50) PRIMARY KEY,
    recruit_id VARCHAR(50),
    user_id VARCHAR(50),
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 募集に参加できないようにBANする管理テーブル
CREATE TABLE join_ban (
    recruit_id VARCHAR(50),
    user_id VARCHAR(50),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (recruit_id, user_id)
);