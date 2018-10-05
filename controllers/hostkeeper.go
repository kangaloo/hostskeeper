package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"hostskeeper/models/db"
	"hostskeeper/models/misc"
	"strings"
)

type AddHosts struct {
	beego.Controller
}

func (c *AddHosts) AddByIDs() {
	id := c.Ctx.Request.FormValue("ids")
	ids := strings.Split(id, ",")

	for _, id := range ids {
		if !misc.IsID(id) {
			c.Ctx.ResponseWriter.Write([]byte(id))
			c.Ctx.ResponseWriter.Write([]byte(" is not a valid id"))
			return
		}
	}

	c.Ctx.ResponseWriter.Write([]byte{})
	fmt.Println(ids)
	return
}

func (c *AddHosts) AddByIPs() {
	ip := c.GetString("ips")
	ips := strings.Split(ip, ",")

	for _, ip := range ips {
		if !misc.IsIP(ip) {
			s := fmt.Sprintf("%s is not a valid ip\n", ip)
			c.Ctx.ResponseWriter.Write([]byte(s))
			return
		}
	}

	return
}

type Host struct {
	beego.Controller
}

func (c *Host) ListAll() {
	hs, err := db.GetAllHost(db.DB)

	if err != nil {
		logs.Error(err)
		c.Ctx.ResponseWriter.Write([]byte("error"))
		return
	}

	hosts, err := db.ConvertJson(db.ConvertMap(hs))

	if err != nil {
		logs.Error(err)
		c.Ctx.ResponseWriter.Write([]byte("error"))
		return
	}

	c.Ctx.ResponseWriter.Write(hosts)
}

func (c *Host) ListInit() {

	i, err := c.GetBool("init")

	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return
	}

	hs, err := db.GetHostsIsInit(db.DB, i)

	if err != nil {
		logs.Error(err)
		c.Ctx.ResponseWriter.Write([]byte("error"))
		return
	}

	hosts, err := db.ConvertJson(db.ConvertMap(hs))

	if err != nil {
		logs.Error(err)
		c.Ctx.ResponseWriter.Write([]byte("error"))
		return
	}

	c.Ctx.ResponseWriter.Write(hosts)
}

func (c *Host) GetByIp() {
	ip := c.GetString("ip")
	// ip合法性检查
	h, err := db.GetHostsByIP(db.DB, ip)

	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return
	}

	host := bytes.NewBuffer(make([]byte, 0, 256))

	hs, err := json.Marshal(h)

	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return
	}

	err = json.Indent(host, hs, "", "\t")

	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return
	}

	c.Ctx.ResponseWriter.Write(host.Bytes())
	return
}

func (c *Host) Add() {
	h := scanHost(c)
	s, err := scanSpec(c)

	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return
	}

	h.Spec = *s

	id, err := db.AddHost(h)
	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return
	}

	logs.Info(id)
	c.Ctx.ResponseWriter.Write([]byte("\n"))
}

func scanSpec(c *Host) (*db.Spec, error) {
	cpu, err := c.GetInt("cpu")
	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return nil, err
	}

	mem, err := c.GetInt("mem")
	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return nil, err
	}

	disk, err := c.GetInt("disk")
	if err != nil {
		c.Ctx.ResponseWriter.Write([]byte("error"))
		logs.Error(err)
		return nil, err
	}

	spec := &db.Spec{
		Cpu:  cpu,
		Mem:  mem,
		Disk: disk,
	}

	return spec, nil
}

func scanHost(c *Host) *db.Host {

	ip := c.GetString("ip")
	hostName := c.GetString("hostname")

	host := &db.Host{
		IP:       ip,
		HostName: hostName,
	}

	return host
}
