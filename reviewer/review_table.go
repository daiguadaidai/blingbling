package reviewer

import "fmt"

type ReviewTable struct {
	SchemaName string
	TableName string
	Alias string
	WhereColumnName map[string]bool
}

func (this *ReviewTable) ToString() string {
	if this.SchemaName != "" {
		return fmt.Sprintf("%v.%v", this.SchemaName, this.TableName)
	}

	return this.TableName
}

func (this *ReviewTable) ToLongString() string {
	var name string
	if this.SchemaName != "" {
		name = this.SchemaName
	}
	if this.TableName != "" {
		name = fmt.Sprintf("%v.%v", name, this.TableName)
	}
	if this.Alias != "" {
		name = fmt.Sprintf("%v.%v", name, this.Alias)
	}

	return name
}

func (this *ReviewTable) Equal(that *ReviewTable) bool {
	return this.ToString() == that.ToString()
}
