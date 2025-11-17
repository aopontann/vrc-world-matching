package vrc_world_matching

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Icon      string    `db:"icon"`
	Bio       string    `db:"bio"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type World struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Thumbnail string    `db:"thumbnail"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Recruit struct {
	ID        string    `db:"id"`
	Content   string    `db:"content"`
	Closed    bool      `db:"closed"`
	UserID    string    `db:"user_id"`
	WorldID   string    `db:"world_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type WantGo struct {
	UserID    string    `db:"user_id"`
	WorldID   string    `db:"world_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type JoinMember struct {
	RecruitID string    `db:"recruit_id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Message struct {
	ID        string    `db:"id"`
	RecruitID string    `db:"recruit_id"`
	UserID    string    `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type JoinBan struct {
	RecruitID string    `db:"recruit_id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Tables struct {
	World      []World
	User       []User
	Result     []Recruit
	WantGo     []WantGo
	JoinMember []JoinMember
	Message    []Message
	JoinBan    []JoinBan
}

func init() {
	var err error
	db, err = sqlx.Connect("mysql", "root:@(localhost:4000)/test?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

var (
	insertWorldsQuery     = "INSERT INTO worlds (id, name, thumbnail, created_at, updated_at) VALUES (:id, :name, :thumbnail, :created_at, :updated_at)"
	insertUsersQuery      = "INSERT INTO users (id, name, icon, bio, created_at, updated_at) VALUES (:id, :name, :icon, :bio, :created_at, :updated_at)"
	insertRecruitsQuery   = "INSERT INTO recruits (id, content, closed, user_id, world_id, created_at, updated_at) VALUES (:id, :content, :closed, :user_id, :world_id, :created_at, :updated_at)"
	insertWantGoQuery     = "INSERT INTO want_go (user_id, world_id, created_at, updated_at) VALUES (:user_id, :world_id, :created_at, :updated_at)"
	insertJoinMemberQuery = "INSERT INTO join_members (recruit_id, user_id, created_at, updated_at) VALUES (:recruit_id, :user_id, :created_at, :updated_at)"
	insertMessageQuery    = "INSERT INTO messages (id, recruit_id, user_id, content, created_at, updated_at) VALUES (:id, :recruit_id, :user_id, :content, :created_at, :updated_at)"
	insertJoinBanQuery    = "INSERT INTO join_ban (recruit_id, user_id, created_at, updated_at) VALUES (:recruit_id, :user_id, :created_at, :updated_at)"
)

// SetUp 引数の構造体のデータをDBのテーブルに登録
// 実装する上での注意：外部キーを付与する場合、インサートするテーブルの順番に気を付ける
func SetUp(tables Tables) error {
	if len(tables.World) != 0 {
		if _, err := db.NamedExec(insertWorldsQuery, tables.World); err != nil {
			return err
		}
	}
	if len(tables.User) != 0 {
		if _, err := db.NamedExec(insertUsersQuery, tables.User); err != nil {
			return err
		}
	}
	if len(tables.Result) != 0 {
		if _, err := db.NamedExec(insertRecruitsQuery, tables.Result); err != nil {
			return err
		}
	}
	if len(tables.WantGo) != 0 {
		if _, err := db.NamedExec(insertWantGoQuery, tables.WantGo); err != nil {
			return err
		}
	}
	if len(tables.JoinMember) != 0 {
		if _, err := db.NamedExec(insertJoinMemberQuery, tables.JoinMember); err != nil {
			return err
		}
	}
	if len(tables.Message) != 0 {
		if _, err := db.NamedExec(insertMessageQuery, tables.Message); err != nil {
			return err
		}
	}
	if len(tables.JoinBan) != 0 {
		if _, err := db.NamedExec(insertJoinBanQuery, tables.JoinBan); err != nil {
			return err
		}
	}
	return nil
}

func CleanUp() error {
	// 外部キー制約でTRUNCATEができないため、一時的に無効化
	_, err := db.Exec("SET foreign_key_checks = 0")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE join_ban")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE join_members")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE messages")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE recruits")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE want_go")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE worlds")
	if err != nil {
		return err
	}
	// 一時的に無効化した外部キー制約を有効化
	_, err = db.Exec("SET foreign_key_checks = 1")
	if err != nil {
		return err
	}
	return nil
}
