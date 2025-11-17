package vrc_world_matching

import "time"

type WorldWithWantGo struct {
	ID           string `db:"id"`
	Name         string `db:"name"`
	Thumbnail    string `db:"thumbnail"`
	WantGoCount  int    `db:"want_go_count"`
	RecruitCount int    `db:"recruit_count"`
}

type RecruitDetail struct {
	ID             string    `db:"id"`
	Content        string    `db:"content"`
	UserID         string    `db:"user_id"`
	Username       string    `db:"user_name"`
	UserIcon       string    `db:"user_icon"`
	WorldID        string    `db:"world_id"`
	WorldName      string    `db:"world_name"`
	WorldThumbnail string    `db:"world_thumbnail"`
	JoiningCount   int       `db:"joining_count"`
	MessagesCount  int       `db:"messages_count"`
	CreatedAt      time.Time `db:"created_at"`
}

// ListWorld 行きたいワールドリストを取得
// 自身が行きたいワールドのみ
// いきたい人数が多い順
// ワールド名でフィルター
// 募集開始日時順
// ワールドIDリスト
func ListWorld() ([]WorldWithWantGo, error) {
	worlds := []WorldWithWantGo{}
	query := `
SELECT
    w.id,
    w.name,
    w.thumbnail,
    IFNULL(wg.want_go_count, 0) AS want_go_count,
    IFNULL(r.recruit_count, 0) AS recruit_count
FROM test.worlds AS w
LEFT JOIN (
    SELECT world_id, COUNT(*) AS want_go_count FROM test.want_go GROUP BY world_id
) AS wg
ON w.id = wg.world_id
LEFT JOIN (
    SELECT
        world_id,
        COUNT(*) AS recruit_count,
        MAX(created_at) AS latest_recruit_start_time
    FROM test.recruits
    WHERE closed = 0
    GROUP BY world_id
) AS r
ON w.id = r.world_id
ORDER BY r.latest_recruit_start_time DESC
`
	err := db.Select(&worlds, query)
	return worlds, err
}

// GetWorldInfo ワールド情報を取得
func GetWorldInfo(worldID string) (WorldWithWantGo, error) {
	query := `
SELECT
    w.id,
    w.name,
    w.thumbnail,
    IFNULL(wg.want_go_count, 0) AS want_go_count,
    IFNULL(r.recruit_count, 0) AS recruit_count
FROM worlds AS w
LEFT JOIN (
    SELECT
        world_id,
        COUNT(*) AS want_go_count
	FROM want_go
	WHERE world_id = :world_id
) AS wg
ON w.id = wg.world_id
LEFT JOIN (
	SELECT
		world_id,
		COUNT(*) AS recruit_count,
		MAX(created_at) AS latest_recruit_start_time
	FROM recruits
	WHERE world_id = :world_id
    AND closed = 0
) AS r
ON w.id = r.world_id
WHERE w.id = :world_id
`
	m := map[string]interface{}{"world_id": worldID}
	rows, err := db.NamedQuery(query, m)
	if err != nil {
		return WorldWithWantGo{}, err
	}

	// 1行のみ取得するSQLなので、最初のループのみ実施
	for rows.Next() {
		var w WorldWithWantGo
		err := rows.StructScan(&w)
		if err != nil {
			return WorldWithWantGo{}, err
		}
		return w, nil
	}

	// 指定したIDのワールドが存在しない場合
	return WorldWithWantGo{}, NotFoundError
}

// RegisterWorld ワールド情報を登録
// 既に登録されている場合はなにもしない
// ワールドIDが不正かチェックをしている
func RegisterWorld(worldID string) error {
	// 既に登録されているワールドかチェック
	var c int
	err := db.Get(&c, "SELECT count(*) FROM test.worlds WHERE id = ?", worldID)
	if err != nil {
		return err
	}
	if c == 1 {
		return nil
	}
	// ワールドIDが正しいものかVRChat APIを使用してチェック
	w, err := GetWorldInfoFromVRChatAPI(worldID)
	if err != nil {
		return err
	}

	_, err = db.NamedExec(`INSERT INTO test.worlds (id, name, thumbnail) VALUES (:id, :name, :thumbnail)`, w)
	return err
}

// RegisterWantGoWorld いきたいワールドを登録
// 既に行きたいワールドとして登録済みの場合 AlreadyRegisteredError を返す
func RegisterWantGoWorld(worldID, userID string) error {
	// 既にいきたいワールドに登録されているかチェック
	var c int
	err := db.Get(&c, "SELECT count(*) FROM test.want_go WHERE user_id = ? AND world_id = ?", userID, worldID)
	if err != nil {
		return err
	}
	if c == 1 {
		return AlreadyRegisteredError
	}

	_, err = db.Exec(`INSERT INTO test.want_go (user_id, world_id) VALUES (?, ?)`, userID, worldID)
	return err
}

// UnregisterWantGoWorld 行きたいワールドを登録解除
// 行きたいワールドとして登録していない、もしくは不正なワールドIDを指定した場合 NotRegisteredError を返す
func UnregisterWantGoWorld(worldID string, userID string) error {
	// 登録解除対象のワールドを行きたいワールドに登録しているかチェック
	var c int
	err := db.Get(&c, "SELECT count(*) FROM test.want_go WHERE user_id = ? AND world_id = ?", userID, worldID)
	if err != nil {
		return err
	}
	if c == 0 {
		return NotRegisteredError
	}

	_, err = db.Exec("DELETE FROM test.want_go WHERE user_id = ? AND world_id = ?", userID, worldID)
	return err
}

// ListRecruitInWorld ワールドの募集リストを取得
// 引数：orderColumn には "created_at" or "joining_count" を指定
// 上記以外の文字を指定した場合 InvalidArgumentError を返す
// 引数：asc 昇順にする場合は TRUE を指定
func ListRecruitInWorld(worldID, orderColumn string, asc bool) ([]RecruitDetail, error) {
	// テーブルに存在しないワールドIDを指定した場合は NotFoundError を返す
	// 指定されたワールドIDがテーブルには存在するが、募集が一つもない場合は空のスライスを返す
	var c int
	err := db.Get(&c, "SELECT count(*) FROM test.worlds WHERE id = ?", worldID)
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return nil, NotFoundError
	}

	var q string
	if orderColumn == "created_at" {
		if asc {
			q = `SELECT * FROM view_recruits_in_world WHERE world_id = ? ORDER BY created_at`
		} else {
			q = `SELECT * FROM view_recruits_in_world WHERE world_id = ? ORDER BY created_at DESC`
		}
	} else if orderColumn == "joining_count" {
		if asc {
			q = `SELECT * FROM view_recruits_in_world WHERE world_id = ? ORDER BY joining_count`
		} else {
			q = `SELECT * FROM view_recruits_in_world WHERE world_id = ? ORDER BY joining_count DESC`
		}
	} else {
		return nil, InvalidArgumentError
	}
	var rd []RecruitDetail
	err = db.Select(&rd, q, worldID)
	return rd, err
}
