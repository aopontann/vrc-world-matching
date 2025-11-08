package vrc_world_matching

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type User struct {
	ID string `db:"first_name"`
}

type World struct {
	ID        string    `db:"id" csv:"id"`
	Name      string    `db:"name" csv:"name"`
	Thumbnail string    `db:"thumbnail" csv:"thumbnail"`
	CreatedAt time.Time `db:"created_at" csv:"created_at"`
	UpdatedAt time.Time `db:"updated_at" csv:"updated_at"`
}

type WantGo struct {
	UserID    string    `db:"user_id" csv:"user_id"`
	WorldID   string    `db:"world_id" csv:"world_id"`
	CreatedAt time.Time `db:"created_at" csv:"created_at"`
	UpdatedAt time.Time `db:"updated_at" csv:"updated_at"`
}

type Tables struct {
	World  []World
	WantGo []WantGo
}

func init() {
	var err error
	db, err = sqlx.Connect("mysql", "root:@(localhost:4000)/test")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

var (
	insertWorldsQuery = "INSERT INTO worlds (id, name, thumbnail, created_at, updated_at) VALUES (:id, :name, :thumbnail, :created_at, :updated_at)"
	insertWantGoQuery = "INSERT INTO want_go (user_id, world_id, created_at, updated_at) VALUES (:user_id, :world_id, :created_at, :updated_at)"
)

// SetUp 引数の構造体のデータをDBのテーブルに登録
// 実装する上での注意：外部キーを付与する場合、インサートするテーブルの順番に気を付ける
func SetUp(tables Tables) error {
	if len(tables.World) != 0 {
		if _, err := db.NamedExec(insertWorldsQuery, tables.World); err != nil {
			return err
		}
	}
	if len(tables.WantGo) != 0 {
		if _, err := db.NamedExec(insertWantGoQuery, tables.WantGo); err != nil {
			return err
		}
	}
	return nil
}

func CleanUp() error {
	_, _ = db.Exec("TRUNCATE TABLE join_ban")
	_, _ = db.Exec("TRUNCATE TABLE join_members")
	_, _ = db.Exec("TRUNCATE TABLE messages")
	_, _ = db.Exec("TRUNCATE TABLE recruits")
	_, _ = db.Exec("TRUNCATE TABLE users")
	_, _ = db.Exec("TRUNCATE TABLE want_go")
	_, _ = db.Exec("TRUNCATE TABLE worlds")
	return nil
}
