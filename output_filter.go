package main

import (
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"strings"
)

const FilterTypeCaller = "caller"
const FilterTypeCallee = "callee"

func isFilterFn(caller, callee *callgraph.Node, inFns, ignoreFns []string, dpMap map[string]bool, filterType string) bool {
	if filterType == FilterTypeCaller {
		return checkTargetCaller(caller, callee, inFns, ignoreFns, dpMap)
	} else if filterType == FilterTypeCallee {
		return checkTargetCallee(caller, callee, inFns, ignoreFns, dpMap)
	}
	return true
}
func checkTargetCallee(caller, callee *callgraph.Node, ignoreFns, exFns []string, dpMap map[string]bool) bool {
	for _, fn := range exFns {
		if strings.HasPrefix(caller.Func.Name(), fn) {
			return false
		}
	}
	if len(ignoreFns) > 0 {
		if dpMap[callee.Func.Name()] {
			return true
		}
		for _, fn := range ignoreFns {
			if strings.HasPrefix(callee.Func.Name(), fn) {
				return true
			}
		}
		return false
	}
	return true
}
func checkTargetCaller(caller, callee *callgraph.Node, inFns, ignoreFns []string, dpMap map[string]bool) bool {
	return checkTargetCallee(callee, caller, inFns, ignoreFns, dpMap)
}

func isFocused(edge *callgraph.Edge, focusPkg []*types.Package) bool {
	fromFocused := false
	toFocused := false
	if len(focusPkg) == 0 || focusPkg[0] == nil {
		return true
	}
	for _, pkg := range focusPkg {
		caller := edge.Caller
		callee := edge.Callee
		if pkg != nil && (caller.Func.Pkg.Pkg.Path() == pkg.Path() || callee.Func.Pkg.Pkg.Path() == pkg.Path()) {
			return true
		}
		for _, e := range caller.In {
			if !isSynthetic(e) && pkg != nil &&
				e.Caller.Func.Pkg.Pkg.Path() == pkg.Path() {
				fromFocused = true
				break
			}
		}
		for _, e := range callee.Out {
			if !isSynthetic(e) && pkg != nil &&
				e.Callee.Func.Pkg.Pkg.Path() == pkg.Path() {
				toFocused = true
				break
			}
		}
		if fromFocused && toFocused {
			logf("edge semi-focus: %s", edge)
			return true
		}
	}
	return false
}

func inIncludes(node *callgraph.Node, inPkgs []string) bool {
	pkgPath := node.Func.Pkg.Pkg.Path()
	for _, p := range inPkgs {
		if strings.HasPrefix(pkgPath, p) {
			return true
		}
	}
	return false
}

func inLimits(node *callgraph.Node, exPkgs []string) bool {
	pkgPath := node.Func.Pkg.Pkg.Path()
	for _, p := range exPkgs {
		if strings.HasPrefix(pkgPath, p) {
			return true
		}
	}
	return false
}

func inIgnores(node *callgraph.Node, ignorePkgs []string) bool {
	pkgPath := node.Func.Pkg.Pkg.Path()
	for _, p := range ignorePkgs {
		if strings.HasPrefix(pkgPath, p) {
			return true
		}
	}
	return false
}
func isInter(edge *callgraph.Edge) bool {
	//caller := edge.Caller
	callee := edge.Callee
	if callee.Func.Object() != nil && !callee.Func.Object().Exported() {
		return true
	}
	return false
}
