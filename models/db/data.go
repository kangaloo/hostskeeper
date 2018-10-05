package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
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

func GetHostsIsInit(db *sqlx.DB, init bool) ([]Host, error) {
	hosts, err := getHostsWithQuery(db, isInit, init)

	if err != nil {
		return nil, err
	}

	err = getFiles(db, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func scanHosts(rows *sql.Rows) ([]Host, error) {
	defer rows.Close()

	hosts := make([]Host, 0, 256)

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

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	json.Indent(buf, hosts, "", "\t")

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

func getHostsWithQuery(db *sqlx.DB, query string, args ...interface{}) ([]Host, error) {
	rows, err := getRowsWithQuery(db, query, args...)

	if err != nil {
		return nil, err
	}

	return scanHosts(rows)
}

func getRowsWithQuery(db *sqlx.DB, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func scanRows(rows *sqlx.Rows, dest []interface{}) error {
	return nil
}

func getSpecId(s *Spec) (int, error) {
	if s == nil {
		return 0, errors.New("nil spec")
	}

	specs, err := getSpecs(DB, s.Cpu, s.Mem, s.Disk)

	if err != nil {
		return 0, err
	}

	if len(specs) > 1 {
		return 0, errors.New("too many specs")
	}

	if len(specs) == 0 {
		// 不存在时创建一个并返回规格id
		id, err := addSpec(DB, s.Cpu, s.Mem, s.Disk)
		if err != nil {
			return 0, err
		}

		return int(id), nil
	}

	// 存在这种规格时，直接返回规格id
	return specs[0].ID, nil
}

func getSpecs(db *sqlx.DB, c, m, d int) ([]Spec, error) {

	if c == 0 || m == 0 || d == 0 {
		return nil, errors.New("wrong spec")
	}

	query := "select id, cpu, mem, disk from specs where cpu = ? and mem = ? and disk = ?"
	rows, err := getRowsWithQuery(db, query, c, m, d)

	if err != nil {
		return nil, err
	}

	return scanSpecs(rows)
}

func scanSpecs(rows *sql.Rows) ([]Spec, error) {
	defer rows.Close()

	var specs []Spec

	for rows.Next() {
		s := &Spec{}
		err := rows.Scan(&(s.ID), &(s.Cpu), &(s.Mem), &(s.Disk))
		if err != nil {
			return nil, err
		}

		specs = append(specs, *s)
	}

	return specs, nil
}

func addSpec(db *sqlx.DB, c, m, d int) (int64, error) {
	res, err := db.Exec("insert into specs(cpu, mem, disk) VALUE (?, ?, ?)", c, m, d)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func insertHost(db *sqlx.DB, h *Host) (int, error) {

	res, err := db.Exec("insert into hosts(ip, host_name, spec_id) VALUE (?, ?, ?)", h.IP, h.HostName, h.SpecID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func AddHost(h *Host) (int, error) {

	specId, err := getSpecId(&(h.Spec))

	if err != nil {
		return 0, err
	}

	h.SpecID = specId

	return insertHost(DB, h)
}
