/*

create database db_thinkgo;

create table t_user(
id int(11) primary key auto_increment,
username varchar(255) not null,
password varchar(255) not null,
join_time date not null
);

create table userdetail(
id int(11) primary key auto_increment,
uid int(11) not null,
intro varchar(255),
position varchar(255)
);

*/

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db, err := sql.Open("mysql", "root:@/db_thinkgo?charset=utf8")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("insert into t_user(username, password, join_time) values(?, ?, ?)")
	checkErr(err)

	res, err := stmt.Exec("huangfugui", "qwe123", "2017-10-07")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Printf("insert回填id：%v\n", id)

	//更新数据
	stmt, err = db.Prepare("update t_user set password = ? where id = ?")
	checkErr(err)

	res, err = stmt.Exec("geek2017", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Printf("影响行数：%v\n", affect)

	//查询数据
	rows, err := db.Query("select * from t_user")
	checkErr(err)

	for rows.Next() {
		var id int
		var username string
		var password string
		var joinTime string
		err = rows.Scan(&id, &username, &password, &joinTime)
		checkErr(err)
		fmt.Printf("id:%v, username:%v, password:%v, joinTime:%v\n", id, username, password, joinTime)
	}

	//删除数据
	stmt, err = db.Prepare("delete from t_user where id = ?")
	checkErr(err)

	res, err = stmt.Exec(4)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Printf("影响行数：%v\n", affect)

	db.Close()
}

func checkErr(err error){
	if err != nil {
		panic(err)
	}
}
