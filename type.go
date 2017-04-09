package drtestlink

type testcase struct {
	idTC          string
	idTS          string
	nomeTS        string
	nome          string
	sommario      string
	precondizioni string
	elencopassi   []passotestcase
}
type testcases []testcase

func (a testcases) Len() int           { return len(a) }
func (a testcases) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a testcases) Less(i, j int) bool { return a[i].idTC < a[j].idTC }

type passotestcase struct {
	numerostep      string
	azione          string
	risultatoatteso string
}

type suite struct {
	id   string
	nome string
}
