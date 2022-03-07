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
