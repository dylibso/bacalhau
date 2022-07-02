// Code generated by "stringer -type=EngineType --trimprefix=Engine"; DO NOT EDIT.

package executor

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[engineUnknown-0]
	_ = x[EngineNoop-1]
	_ = x[EngineDocker-2]
	_ = x[EngineWasm-3]
	_ = x[EngineLanguage-4]
	_ = x[EnginePythonWasm-5]
	_ = x[engineDone-6]
}

const _EngineType_name = "engineUnknownNoopDockerWasmLanguagePythonWasmengineDone"

var _EngineType_index = [...]uint8{0, 13, 17, 23, 27, 35, 45, 55}

func (i EngineType) String() string {
	if i < 0 || i >= EngineType(len(_EngineType_index)-1) {
		return "EngineType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EngineType_name[_EngineType_index[i]:_EngineType_index[i+1]]
}
