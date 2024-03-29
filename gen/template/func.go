package template

import (
	"fmt"
	"strings"

	"github.com/gokch/ornn/gen/util"
)

func Select(args []string, tpls []string, query string, selectSingle bool, structName string, instanceName string, retName, retItemName, retItemType string) string {
	var bodyRetDeclare, bodyRetSet string
	if selectSingle == true {
		bodyRetSet = fmt.Sprintf("%s = scan\n\tbreak", retItemName)
	} else {
		bodyRetDeclare = fmt.Sprintf("\n%s = make(%s, 0, 100)", retItemName, retItemType)
		bodyRetSet = fmt.Sprintf("%s = append(%s, scan)", retItemName, retItemName)
	}
	return parseTemplate("func_select.template", map[string]interface{}{
		"arg":      genQuery_body_setArgs(args),
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"struct":   structName,
		"instance": instanceName,
		"body":     bodyRetDeclare,
		"scan":     retName,
		"retSet":   bodyRetSet,
		"ret":      retItemName,
	})
}

func Insert(args []string, tpls []string, query string, insertMulti bool, structName, instanceName string) string {
	var multiInsert, genArgs string
	if insertMulti == true { // multi insert
		queryVal := util.Util_ExportInsertQueryValues(query)
		query = strings.TrimSuffix(query, ";")
		query += "%s"
		genArgs = genQuery_body_multiInsertProc(args)
		multiInsert = genQuery_body_multiInsert(queryVal)
	} else { // insert
		genArgs = genQuery_body_setArgs(args)
	}
	return parseTemplate("func_insert.template", map[string]interface{}{
		"arg":      genArgs,
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"multi":    multiInsert,
		"struct":   structName,
		"instance": instanceName,
	})
}

func Update(args []string, tpls []string, query string, structName, instanceName string) string {
	return parseTemplate("func_update.template", map[string]interface{}{
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"arg":      genQuery_body_setArgs(args),
		"struct":   structName,
		"instance": instanceName,
	})
}

func Delete(args []string, query string, tpls []string, structName, instanceName string) string {
	return parseTemplate("func_delete.template", map[string]interface{}{
		"arg":      genQuery_body_setArgs(args),
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"struct":   structName,
		"instance": instanceName,
	})
}
