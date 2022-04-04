package repository

import "time"

var todoSchema = `
CREATE TABLE IF NOT EXISTS todo (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT DEFAULT "",
  priority INT DEFAULT 1,
  duedate TIMESTAMP DEFAULT 0,
  completed TINYINT(1) DEFAULT 0,
  completiondate TIMESTAMP DEFAULT 0
);`

// database model
type Todo struct {
	Id             int
	Name           string
	Description    string
	Priority       int
	Duedate        time.Time
	Completed      int
	Completiondate time.Time
}

var customerSchema = `
CREATE TABLE IF NOT EXISTS customer (
  id int(11) NOT NULL AUTO_INCREMENT,
  username varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  email varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  password_salt varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
	password_hash varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4`

// database model
type Customer struct {
	Id            int
	Username      string
	Email         string
	Password_hash string
	Password_salt string
}

var serverSchema = `
CREATE TABLE IF NOT EXISTS server (
  id int(11) NOT NULL AUTO_INCREMENT,
  username varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  email varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  password_salt varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
	password_hash varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4`

// database model
type Server struct {
	Id            int
	Username      string
	Email         string
	Password_hash string
	Password_salt string
}
