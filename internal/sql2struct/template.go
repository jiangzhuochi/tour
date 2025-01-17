package sql2struct

import (
	"fmt"
	"os"
	"text/template"

	"github.com/go-programming-tour-book/tour/internal/word"
)

const strcutTpl = `
// {{.TableComment}}
type {{.TableName | ToCamelCase}} struct {
{{range .Columns -}} 
    {{$length := len .Comment -}}
    {{if gt $length 0}}    // {{.Comment}}{{else}}    // {{.Name}}{{end -}}
    {{$length := len .ColumnKey -}}
    {{if gt $length 0}}    {{.ColumnKey}}{{end -}}
    {{if eq .IsNullable "YES"}}    Nullable{{end}}
    {{$typeLen := len .Type -}} 
    {{if gt $typeLen 0}}{{.Name|ToCamelCase}} {{.Type}} {{.Tag}}{{else}}{{.Name}}{{end}}
{{end}}}

func (model {{.TableName|ToCamelCase}}) TableName() string {
    return "{{.TableName}}"
}

`

type StructTemplate struct {
	strcutTpl string
}

type StructColumn struct {
	Name       string
	Type       string
	ColumnKey  string
	IsNullable string
	Tag        string
	Comment    string
}

type StructTemplateDB struct {
	TableName    string
	TableComment string
	Columns      []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{strcutTpl: strcutTpl}
}

func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {
		tag := fmt.Sprintf("`json:%q`", column.ColumnName)
		tplColumns = append(tplColumns, &StructColumn{
			Name:       column.ColumnName,
			Type:       DBTypeToStructType[column.DataType],
			ColumnKey:  column.ColumnKey,
			IsNullable: column.IsNullable,
			Tag:        tag,
			Comment:    column.ColumnComment,
		})
	}

	return tplColumns
}

func (t *StructTemplate) Generate(
	tableName string,
	tableComment string,
	tplColumns []*StructColumn,
) error {

	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase": word.UnderscoreToUpperCamelCase,
	}).Parse(t.strcutTpl))

	tplDB := StructTemplateDB{
		TableName: tableName,
		TableComment: "Automatically generate from table: " +
			tableName + " " + tableComment,
		Columns: tplColumns,
	}
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}

	return nil
}
