package dashboard

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

type FaultStatusCount struct {
	OpenCount int `json:"openCount"`
	ClosedCount int `json:"closedCount"`
}

type Repository interface {
	getCarCount(string)(int,error)
	getCarCountByProject(string)([]map[string]interface{},error)
	getFaultCountByType(string)([]map[string]interface{},error)
	getDashboardYears(string)([]map[string]interface{},error)
	getDashboardProjects(string)([]map[string]interface{},error)
	getDashboardTypes(string)([]map[string]interface{},error)
	getDashboardSpecs(string)([]map[string]interface{},error)
	getFaultCountByStatus(string)(*FaultStatusCount,error)
	getFaultList(string)([]map[string]interface{},error)
	query(sql string)([]map[string]interface{},error)
	closeFault(diagReport string,remark string)
}

type DefatultRepository struct {
	DB *sql.DB
}

func (repo *DefatultRepository)query(sql string)([]map[string]interface{},error){
	rows, err := repo.DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows)
}

func (repo *DefatultRepository)getCarCount(whereStr string)(int,error){
	row := repo.DB.QueryRow("select count(*) as count from diag_result "+whereStr)
    var count int = 0
	if err := row.Scan(&count); err != nil {
        log.Println("getCarCount error")
		log.Println(err)
        return 0,nil
    }
	return count, nil
}

func (repo *DefatultRepository)toMap(rows *sql.Rows)([]map[string]interface{},error){
	cols,_:=rows.Columns()
	columns:=make([]interface{},len(cols))
	colPointers:=make([]interface{},len(cols))
	for i,_:=range columns {
		colPointers[i] = &columns[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err:= rows.Scan(colPointers...)
		if err != nil {
			log.Println(err)
			return nil,nil
		}
		row:=make(map[string]interface{})
		for i,colName :=range cols {
			val:=colPointers[i].(*interface{})
			switch (*val).(type) {
			case []byte:
				row[colName]=string((*val).([]byte))
			default:
				row[colName]=*val
			} 
		}
		list=append(list,row)
	}
	return list,nil
}

func (repo *DefatultRepository)getFaultList(whereStr string)([]map[string]interface{},error){
	rows, err := repo.DB.Query("select * from  diag_result "+whereStr+" order by status asc,time desc limit 0,500")
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows)
}

func (repo *DefatultRepository)getCarCountByProject(whereStr string)([]map[string]interface{},error){
	rows, err := repo.DB.Query("select * from (select project_num as ProjectNum, count(*) as count from diag_result "+whereStr+" group by project_num) as t order by count desc limit 0,10")
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows) 
}

func (repo *DefatultRepository)getFaultCountByType(whereStr string)([]map[string]interface{},error){
	rows,err:= repo.DB.Query("select * from (select type,count(1) as count from diag_result "+whereStr+" group by type) as t order by count desc limit 0,10 ")
	if err!=nil {
		log.Println(err)
		return nil,nil
	}

	defer rows.Close()

	return repo.toMap(rows)
}

func (repo *DefatultRepository)getFaultCountByStatus(whereStr string)(*FaultStatusCount,error){
	var statusCount FaultStatusCount
	row:= repo.DB.QueryRow("SELECT count(if(status=0,true,null)) as openCount,count(if(status=1,true,null)) as closedCount FROM diag_result "+whereStr)
	if err := row.Scan(&statusCount.OpenCount, &statusCount.ClosedCount); err != nil {
        log.Println("getFaultCountByType error")
		log.Println(err)
    } 
	return &statusCount, nil
}

func (repo *DefatultRepository)closeFault(diagReport string,remark string){
	repo.DB.Exec("update DiagResult set Status='1',Remark=? where DiagReport = ?", remark,diagReport)
}

func (repo *DefatultRepository)Connect(server string,user string,password string,dbName string){
	// Capture connection properties.
    cfg := mysql.Config{
        User:   user,
        Passwd: password,
        Net:    "tcp",
        Addr:   server,
        DBName: dbName,
		AllowNativePasswords:true,
    }
    // Get a database handle.
    var err error
    repo.DB, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := repo.DB.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    log.Println("connect to mysql server "+server)
}

func (repo *DefatultRepository)getDashboardYears(whereStr string)([]map[string]interface{},error){
	rows,err:= repo.DB.Query("select * from (SELECT distinct(SUBSTR(time,1,4)) as year FROM `diag_result` "+whereStr+") as t order by year desc")
	if err!=nil {
		log.Println(err)
		return nil,nil
	}

	defer rows.Close()

	return repo.toMap(rows)
}

func (repo *DefatultRepository)getDashboardProjects(whereStr string)([]map[string]interface{},error){
	rows,err:= repo.DB.Query("select * from (SELECT distinct(project_num) as project FROM `diag_result` "+whereStr+") as t order by project asc")
	if err!=nil {
		log.Println(err)
		return nil,nil
	}

	defer rows.Close()

	return repo.toMap(rows)
}

func (repo *DefatultRepository)getDashboardTypes(whereStr string)([]map[string]interface{},error){
	rows,err:= repo.DB.Query("select * from (SELECT distinct(type) as type FROM `diag_result` "+whereStr+") as t order by type asc")
	if err!=nil {
		log.Println(err)
		return nil,nil
	}

	defer rows.Close()

	return repo.toMap(rows)
}

func (repo *DefatultRepository)getDashboardSpecs(whereStr string)([]map[string]interface{},error){
	rows,err:= repo.DB.Query("select * from (SELECT distinct(specifications) as spec FROM `diag_result` "+whereStr+") as t order by spec asc")
	if err!=nil {
		log.Println(err)
		return nil,nil
	}

	defer rows.Close()

	return repo.toMap(rows)
}