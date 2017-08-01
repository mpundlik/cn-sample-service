// File generated by github.com/ungerik/pkgreflect
package tap

import "reflect"

var Types = map[string]reflect.Type{
	"SwInterfaceTapDetails": reflect.TypeOf((*SwInterfaceTapDetails)(nil)).Elem(),
	"SwInterfaceTapDump":    reflect.TypeOf((*SwInterfaceTapDump)(nil)).Elem(),
	"TapConnect":            reflect.TypeOf((*TapConnect)(nil)).Elem(),
	"TapConnectReply":       reflect.TypeOf((*TapConnectReply)(nil)).Elem(),
	"TapDelete":             reflect.TypeOf((*TapDelete)(nil)).Elem(),
	"TapDeleteReply":        reflect.TypeOf((*TapDeleteReply)(nil)).Elem(),
	"TapModify":             reflect.TypeOf((*TapModify)(nil)).Elem(),
	"TapModifyReply":        reflect.TypeOf((*TapModifyReply)(nil)).Elem(),
}

var Functions = map[string]reflect.Value{
	"NewSwInterfaceTapDetails": reflect.ValueOf(NewSwInterfaceTapDetails),
	"NewSwInterfaceTapDump":    reflect.ValueOf(NewSwInterfaceTapDump),
	"NewTapConnect":            reflect.ValueOf(NewTapConnect),
	"NewTapConnectReply":       reflect.ValueOf(NewTapConnectReply),
	"NewTapDelete":             reflect.ValueOf(NewTapDelete),
	"NewTapDeleteReply":        reflect.ValueOf(NewTapDeleteReply),
	"NewTapModify":             reflect.ValueOf(NewTapModify),
	"NewTapModifyReply":        reflect.ValueOf(NewTapModifyReply),
}

var Variables = map[string]reflect.Value{}

var Consts = map[string]reflect.Value{
	"VlAPIVersion": reflect.ValueOf(VlAPIVersion),
}
