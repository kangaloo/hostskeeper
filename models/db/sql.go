package db

var allHost = `select h.id, h.ip, h.host_name, h.is_init, s.cpu, s.mem, s.disk from hosts h join specs s on h.spec_id = s.id`
var isInit = `select h.id, h.ip, h.host_name, h.is_init, s.cpu, s.mem, s.disk from hosts h join specs s on h.spec_id = s.id where is_init = ?`
var queryHostWithIP = `select h.id, h.ip, h.host_name, h.is_init, s.cpu, s.mem, s.disk from hosts h join specs s on h.spec_id = s.id where ip = ?`
var queryFileWithIP = `select f.id, f.name, f.path, f.version, v.version from files f join versions v on f.id = v.file_id where v.host_id = ?`
