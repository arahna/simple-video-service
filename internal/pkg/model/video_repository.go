package model

import (
	"database/sql"
	"fmt"
)

type videoRepository struct {
	DB *sql.DB
}

type VideoRepository interface {
	Find(uid string) (*Video, error)
	FindWithStatus(status VideoStatus, search string, offset, limit int) ([]*Video, error)
	FindOneWithStatus(status VideoStatus) (*Video, error)
	Save(v *Video) error
}

func NewVideoRepository(db *sql.DB) VideoRepository {
	return &videoRepository{db}
}

func (r *videoRepository) FindWithStatus(status VideoStatus, search string, offset, limit int) ([]*Video, error) {
	params := []interface{}{status}
	q := `SELECT
			id, uid, title, status, duration, file_name 
		  FROM
		  	video 
		  WHERE
		  	status = ?`
	if search != `` {
		q = q + ` AND name LIKE %?%`
		params = append(params, search)
	}
	if limit > 0 {
		q = fmt.Sprintf(`%s LIMIT %d`, q, limit)
	}
	if offset > 0 {
		q = fmt.Sprintf(`%s OFFSET %d`, q, offset)
	}
	rows, err := r.DB.Query(q, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*Video
	for rows.Next() {
		video, err := r.scanVideo(rows)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (r *videoRepository) FindOneWithStatus(status VideoStatus) (*Video, error) {
	q := `SELECT
			id, uid, title, status, duration, file_name 
		  FROM
		  	video 
		  WHERE
		  	status = ?
		  LIMIT 1`
	rows, err := r.DB.Query(q, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	video, err := r.scanVideo(rows)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (r *videoRepository) Find(uid string) (*Video, error) {
	rows, err := r.DB.Query(`SELECT id, uid, title, status, duration, file_name FROM video WHERE uid = ?`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	video, err := r.scanVideo(rows)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (r *videoRepository) Save(v *Video) error {
	if isNew := v.Id == 0; isNew {
		return r.insert(v)
	} else {
		return r.update(v)
	}
}

func (r *videoRepository) insert(v *Video) error {
	res, err := r.DB.Exec(`
		INSERT INTO
			video
		SET 
			uid = ?,
			title = ?,
			status = ?,
			duration = ?,
			file_name = ?
	`, v.Uid, v.Title, v.Status, v.Duration, v.FileName)

	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		v.Id = uint(id)
	}
	return err
}

func (r *videoRepository) update(v *Video) error {
	_, err := r.DB.Exec(`
		UPDATE
			video
		SET 
			title = ?,
			status = ?,
			duration = ?,
			file_name = ?
		WHERE
			id = ?	
	`, v.Title, v.Status, v.Duration, v.FileName, v.Id)

	return err
}

func (r *videoRepository) scanVideo(rows *sql.Rows) (*Video, error) {
	var video Video
	err := rows.Scan(&video.Id, &video.Uid, &video.Title, &video.Status, &video.Duration, &video.FileName)
	return &video, err
}
