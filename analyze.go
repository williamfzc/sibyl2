package sibyl2

import (
	"github.com/dominikbraun/graph"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

// These functions are designed on the top of query.go
// for some higher levels usages
// Starts with `Analyze`

const refLimit = 1024

func isFuncNameInvalid(funcName string) bool {
	// ignore length < 4 functions
	// current calculation can not get the correct results for them
	// which still takes lots of calc time
	tooShort := len(funcName) < 4
	return tooShort
}

func AnalyzeFuncGraph(funcFiles []*extractor.FunctionFileResult, symbolFiles []*extractor.SymbolFileResult) (*FuncGraph, error) {
	reverseCallGraph := graph.New((*FunctionWithPath).GetDesc, graph.Directed())
	callGraph := graph.New((*FunctionWithPath).GetDesc, graph.Directed())

	// speed up cache
	funcFileMap := make(map[string]*extractor.FunctionFileResult, len(funcFiles))
	for _, each := range funcFiles {
		funcFileMap[each.Path] = each
	}

	symbolMap := make(map[string]map[string][]*extractor.Symbol, len(symbolFiles))
	for _, each := range symbolFiles {
		functions, ok := funcFileMap[each.Path]
		// this file only contains symbols
		if !ok {
			continue
		}
		// out of function scope
		validSymbols := make([]*extractor.Symbol, 0)
		for _, eachS := range each.Units {
			for _, eachF := range functions.Units {
				if eachF.BodySpan.HasInteraction(eachS.GetSpan()) {
					validSymbols = append(validSymbols, eachS)
					break
				}
			}
		}
		symbolMap[each.Path] = make(map[string][]*extractor.Symbol)
		for _, eachSymbol := range validSymbols {
			if symbolList, ok := symbolMap[each.Path][eachSymbol.GetIndexName()]; ok {
				symbolList = append(symbolList, eachSymbol)
			} else {
				symbolList := make([]*extractor.Symbol, 0)
				symbolList = append(symbolList, eachSymbol)
				symbolMap[each.Path][eachSymbol.GetIndexName()] = symbolList
			}
		}
	}
	core.Log.Infof("symbol clean up")

	funcRefMap := make(map[string][]*SymbolWithPath, 0)
	for symbolPath, symbolNameMap := range symbolMap {
		for _, eachFuncFile := range funcFiles {
			for _, eachFunc := range eachFuncFile.Units {
				index := eachFunc.GetIndexName()
				if isFuncNameInvalid(index) {
					continue
				}

				funcName := eachFunc.GetIndexName()
				refs, ok := symbolNameMap[funcName]
				if !ok {
					continue
				}
				refWithPaths := make([]*SymbolWithPath, 0, len(refs))
				for _, eachRef := range refs {
					swp := &SymbolWithPath{
						Symbol: eachRef,
						Path:   symbolPath,
					}
					refWithPaths = append(refWithPaths, swp)
				}
				funcRefMap[funcName] = append(funcRefMap[funcName], refWithPaths...)
			}
		}
	}
	core.Log.Infof("symbol refs finished")

	// fill graph with vertex
	for _, eachFuncFile := range funcFiles {
		for _, eachFunc := range eachFuncFile.Units {
			// multi graphs shared
			fwp := &FunctionWithPath{
				eachFunc,
				eachFuncFile.Path,
				eachFuncFile.Language,
			}
			err := reverseCallGraph.AddVertex(fwp)
			if err != nil {
				return nil, err
			}
			err = callGraph.AddVertex(fwp)
			if err != nil {
				return nil, err
			}
		}
	}
	core.Log.Infof("vertex filled")

	// build relationship
	for _, eachFuncFile := range funcFiles {
		core.Log.Infof("file %s , methods: %d", eachFuncFile.Path, len(eachFuncFile.Units))
		for _, eachFunc := range eachFuncFile.Units {
			refs, ok := funcRefMap[eachFunc.GetIndexName()]
			if !ok {
				continue
			}

			// in some languages (like java) which has `override`
			// will create thousands of refs for some special methods (toString, etc.)
			// which makes the final graph very, very large
			// and at the most time these methods will not be analyzed
			if len(refs) > refLimit {
				// happen very easily in big repo, suppress it
				//core.Log.Warnf("func %s exceed the ref limit %d, now %d", eachFunc.GetIndexName(), refLimit, len(refs))
				continue
			}

			for _, each := range refs {
				targetFuncFile, ok := funcFileMap[each.Path]
				if !ok {
					continue
				}
				for _, eachMatchFunc := range targetFuncFile.Units {
					if eachMatchFunc.BodySpan.HasInteraction(each.GetSpan()) {
						// match
						// exclude itself
						if eachMatchFunc.GetDesc() == eachFunc.GetDesc() {
							continue
						}

						// eachFunc referenced by eachMatchFunc
						reverseCallGraph.AddEdge(eachFunc.GetDesc(), eachMatchFunc.GetDesc())
						callGraph.AddEdge(eachMatchFunc.GetDesc(), eachFunc.GetDesc())
						break
					}
				}
			}
		}
	}
	fg := &FuncGraph{
		ReverseCallGraph: WrapFuncGraph(reverseCallGraph),
		CallGraph:        WrapFuncGraph(callGraph),
	}
	return fg, nil
}
