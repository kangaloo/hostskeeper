package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
)

func GetAllHost(db *sqlx.DB) ([]Host, error) {
	hosts, err := getHostsWithQuery(db, allHost)
	if err != nil {
		return nil, err
	}

	err = getFiles(db, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func GetHostsIsInit(db *sqlx.DB, init bool) ([]Host, error)  {
	hosts, err :=  getHostsWithQuery(db, isInit, init)

	if err != nil {
		return nil, err
	}

	err = getFiles(db, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func scanHosts(rows *sql.Rows) ([]Host, error)  {
	defer rows.Close()

	hosts := make([]Host, 0 , 256)

	for rows.Next() {
		h := &Host{}
		err := rows.Scan(&(h.ID), &(h.IP), &(h.HostName), &(h.IsInit), &(h.Spec.Cpu), &(h.Spec.Mem), &(h.Spec.Disk))

		if err != nil {
			return nil, err
		}

		hosts = append(hosts, *h)
	}

	return hosts, nil
}

func ConvertMap(src []Host) map[string]Host {

	dest := make(map[string]Host, 256)

	for _, v := range src {
		dest[v.IP] = v
	}

	return dest
}

func ConvertJson(src map[string]Host) ([]byte, error) {

	hosts, err := json.Marshal(src)

	dest := bytes.NewBuffer(make([]byte, 0, 1024))

	json.Indent(dest, hosts, "", "\t")

	if err != nil {
		return nil, err
	}

	return dest.Bytes(), nil
}

func GetHostsByIP(db *sqlx.DB, ip string) ([]Host, error) {
	hosts, err := getHostsWithQuery(db, queryHostWithIP, ip)

	if err != nil {
		return nil, err
	}

	err = getFiles(db, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

// 根据hostID从files和versions表联合查询的结果中获取主机文件列表
func getFilesById(db *sqlx.DB, id int) ([]File, error) {
	rows, err := db.Query(queryFileWithIP, id)
	if err != nil {
		return nil, err
	}

	return scanFiles(rows)
}

func scanFiles(rows *sql.Rows) ([]File, error) {
	defer rows.Close()

	var files []File

	for rows.Next() {
		f := &File{}
		err := rows.Scan(&(f.ID), &(f.Name), &(f.Path), &(f.Version), &(f.PresentVersion.Version))

		if err != nil {
			return nil, err
		}

		files = append(files, *f)
	}

	return files, nil
}

func getFiles(db *sqlx.DB, hosts *[]Host) error {
	for i, h := range *hosts {
		files, err := getFilesById(db, h.ID)
		if err != nil {
			return err
		}
		(*hosts)[i].Files = files
	}
	return nil
}

func getHostsWithQuery(db *sqlx.DB, query string, args... interface{}) ([]Host, error) {
	rows, err := getRowsWithQuery(db, query, args...)

	if err != nil {
		return nil, err
	}

	return scanHosts(rows)
}

func getRowsWithQuery(db *sqlx.DB, query string, args... interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func scanRows(rows *sqlx.Rows, dest []interface{}) error {
	return nil
}

