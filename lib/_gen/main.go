package main

import (
	"context"
	"fmt"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"os"
	"sort"
)

func main() {
	ctx := context.Background()
	rt := wazero.NewRuntime(ctx)
	wasm, err := os.ReadFile("lib/ts.wasm")
	check(err)
	mod, err := rt.CompileModule(ctx, wasm)
	check(err)

	fmt.Println("[IMPORT]")
	for _, efn := range mod.ImportedMemories() {
		moduleName, memoryName, isImport := efn.Import()
		if isImport {
			bounds := "[" + fmt.Sprint(efn.Min()) + ","
			if memoryMax, bounded := efn.Max(); bounded {
				bounds += fmt.Sprint(memoryMax)
			} else {
				bounds += "..."
			}
			bounds += "]"
			println(" mem ", moduleName+"."+memoryName, bounds)
		}
	}
	for _, ifn := range mod.ImportedFunctions() {
		moduleName, functionName, isImport := ifn.Import()
		if isImport {
			println("func ", moduleName+"."+functionName+valueTypesToString(ifn.ParamTypes())+" => "+valueTypesToString(ifn.ResultTypes()))
		}
	}
	fmt.Println("\n[EXPORT]")
	var exportedMemories []string
	for name, et := range mod.ExportedMemories() {
		bounds := "[" + fmt.Sprint(et.Min()) + ","
		if memoryMax, bounded := et.Max(); bounded {
			bounds += fmt.Sprint(memoryMax)
		} else {
			bounds += "..."
		}
		bounds += "]"
		exportedMemories = append(exportedMemories, " mem  "+name+" "+bounds)
	}
	sort.Strings(exportedMemories)
	for _, em := range exportedMemories {
		println(em)
	}

	var exportedFunctions []string
	for name, efn := range mod.ExportedFunctions() {
		exportedFunctions = append(exportedFunctions, "func  "+name+valueTypesToString(efn.ParamTypes())+" => "+valueTypesToString(efn.ResultTypes()))
	}
	sort.Strings(exportedFunctions)
	for _, ef := range exportedFunctions {
		println(ef)
	}

}

func valueTypesToString(vs []api.ValueType) string {
	s := "("
	for i, v := range vs {
		if i > 0 {
			s += ", "
		}
		s += api.ValueTypeName(v)
	}
	s += ")"
	return s
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
