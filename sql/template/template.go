package template

import "fmt"

func genQuery_body_setArgs(_arrs_arg []string) (genArg string) {
	var items string
	items = genQuery_body_arg(_arrs_arg)
	if items != "" {
		items += "\n"
	}

	return fmt.Sprintf(`args := make([]interface{}, 0, %d)
args = append(args, I_to_arri(%s)...)
`,
		len(_arrs_arg),
		items,
	)
}

func genQuery_body_arg(args []string) (ret string) {
	for _, arg := range args {
		ret += fmt.Sprintf("\n\t%s,", arg)
	}
	return ret
}
