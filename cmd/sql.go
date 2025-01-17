package cmd

import (
	"log"

	"github.com/go-programming-tour-book/tour/internal/sql2struct"
	"github.com/spf13/cobra"
)

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql转换和处理",
	Long:  "sql转换和处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var sql2structCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql转换",
	Long:  "sql转换",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo := &sql2struct.DBInfo{
			DBType:   dbType,
			Host:     host,
			UserName: username,
			Password: password,
			Charset:  charset,
		}
		dbModel := sql2struct.NewDBModel(dbInfo)
		err := dbModel.Connect()
		if err != nil {
			log.Fatalf("dbModel.Connect err: %v", err)
		}
		columns, err := dbModel.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel.GetColumns err: %v", err)
		}
		tableComment, err := dbModel.GetTableComment(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel.GetTableComment err: %v", err)
		}

		template := sql2struct.NewStructTemplate()
		templateColumns := template.AssemblyColumns(columns)
		err = template.Generate(tableName, tableComment, templateColumns)
		if err != nil {
			log.Fatalf("template.Generate err: %v", err)
		}
	},
}

func init() {
	sqlCmd.AddCommand(sql2structCmd)
	sql2structCmd.Flags().StringVarP(&username, "username", "u", "root", "数据库账号")
	sql2structCmd.Flags().StringVarP(&password, "password", "p", "root", "数据库密码")
	sql2structCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1:3306", "数据库的HOST")
	sql2structCmd.Flags().StringVarP(&charset, "charset", "", "utf8mb4", "数据库的编码")
	sql2structCmd.Flags().StringVarP(&dbType, "type", "", "mysql", "数据库实例类型")
	sql2structCmd.Flags().StringVarP(&dbName, "db", "", "", "数据库名称")
	sql2structCmd.Flags().StringVarP(&tableName, "table", "", "", "表名称")
	sql2structCmd.Flags().SortFlags = false
}
