package drtestlink

import (
	"strconv"

	"github.com/Luxurioust/excelize"
	"github.com/jaytaylor/html2text"

	log "drollo.it/drlog"
)

var xlsx *excelize.File
var nomefoglioxsl string

//inizializza file XLS
func inizializzaXls() error {
	log.Logga(nomemodulo).Info("Imposto XLS")
	xlsx.SetActiveSheet(1)
	nomefoglioxsl = xlsx.GetSheetName(1)
	xlsx.SetCellValue(nomefoglioxsl, "A1", "ID Test Suite")
	xlsx.SetCellValue(nomefoglioxsl, "B1", "Test Suite")
	xlsx.SetCellValue(nomefoglioxsl, "C1", "ID Test Case")
	xlsx.SetCellValue(nomefoglioxsl, "D1", "Test Case")
	xlsx.SetCellValue(nomefoglioxsl, "E1", "Sommario")
	xlsx.SetCellValue(nomefoglioxsl, "F1", "Precondizioni")
	xlsx.SetCellValue(nomefoglioxsl, "G1", "Step")
	xlsx.SetCellValue(nomefoglioxsl, "H1", "Azioni")
	xlsx.SetCellValue(nomefoglioxsl, "I1", "Risultati")
	err = xlsx.SetCellStyle(nomefoglioxsl, "A1", "I1", `{"fill":{"type":"pattern","color":["#2b567c"],"pattern":1}}`)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	err = xlsx.SetCellStyle(nomefoglioxsl, "A1", "I1", `{"alignment": {
        "horizontal": "left",
        "shrink_to_fit": true,
        "vertical": "top",
        "wrap_text": true}}`)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	xlsx.SetColWidth(nomefoglioxsl, "A", "A", 10)
	xlsx.SetColWidth(nomefoglioxsl, "B", "B", 40)
	xlsx.SetColWidth(nomefoglioxsl, "C", "C", 10)
	xlsx.SetColWidth(nomefoglioxsl, "D", "D", 40)
	xlsx.SetColWidth(nomefoglioxsl, "E", "F", 80)
	xlsx.SetColWidth(nomefoglioxsl, "G", "G", 4)
	xlsx.SetColWidth(nomefoglioxsl, "H", "I", 80)
	return nil
}

//scrivi file XSL
func scriviFileXls(nomefile string) error {
	err = xlsx.WriteTo(nomefile)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	return nil
}

//Inizializza file
func inizializzaFile(nomefile string, tipo string) error {

	switch tipo {
	case "xsl":
		inizializzaXls()
		return nil
	case "json":
		return nil
	case "testo":
		return nil
	default:
		log.Logga(nomemodulo).Error("Tipo file sconosciuto su inizializzaFile")
		return nil
	}
}

//scrivi strutturta elenco testcase su file
func scriviStructToFile(tipoOut string) error {
	var primasuite bool
	//Ciclo tutti i testcase
	i := 2
	for elencomyTC, tc := range elencoTCs {
		log.Logga(nomemodulo).Info("Leggo riga:" + strconv.Itoa(elencomyTC) + " TC: " + tc.nome)
		primasuite = true
		for _, passo := range tc.elencopassi {
			log.Logga(nomemodulo).Info("Leggo riga:" + strconv.Itoa(elencomyTC) + " TC step: " + passo.numerostep)
			err = scriviRigaDaStruct(primasuite, i, tc, passo, tipoOut)
			if err != nil {
				log.Logga(nomemodulo).Error(err)
				return err
			}
			primasuite = false
			i++
		}
	}
	return nil
}

//scrive una riga su file
func scriviRigaDaStruct(primasuite bool, riga int, tc testcase, passo passotestcase, tipoOut string) error {
	switch tipoOut {
	case "xsl":
		if primasuite == true {

			// Set value of a cell.
			log.Logga(nomemodulo).Debug("Scrivo su foglio:" + nomefoglioxsl)
			log.Logga(nomemodulo).Debug("Scrivo " + "A" + strconv.Itoa(riga) + " su xls: " + strconv.Itoa(riga) + " " + tc.idTS)
			xlsx.SetCellValue(nomefoglioxsl, "A"+strconv.Itoa(riga), tc.idTS)

			xlsx.SetCellValue(nomefoglioxsl, "B"+strconv.Itoa(riga), tc.nomeTS)
			log.Logga(nomemodulo).Debug("Scrivo B su xls: " + strconv.Itoa(riga) + " " + tc.nomeTS)
			xlsx.SetCellValue(nomefoglioxsl, "C"+strconv.Itoa(riga), tc.idTC)
			log.Logga(nomemodulo).Debug("Scrivo C su xls: " + strconv.Itoa(riga) + " " + tc.idTC)
			xlsx.SetCellValue(nomefoglioxsl, "D"+strconv.Itoa(riga), tc.nome)
			log.Logga(nomemodulo).Debug("Scrivo D su xls: " + strconv.Itoa(riga) + " " + tc.nome)

			xlsx.SetCellValue(nomefoglioxsl, "E"+strconv.Itoa(riga), pulisciHTML(tc.sommario))
			xlsx.SetCellValue(nomefoglioxsl, "F"+strconv.Itoa(riga), pulisciHTML(tc.precondizioni))
			log.Logga(nomemodulo).Debug("Scrivo numero step su xls: " + strconv.Itoa(riga) + " " + passo.numerostep)
		}
		xlsx.SetCellValue(nomefoglioxsl, "G"+strconv.Itoa(riga), pulisciHTML(passo.numerostep))
		xlsx.SetCellValue(nomefoglioxsl, "H"+strconv.Itoa(riga), pulisciHTML(passo.azione))
		xlsx.SetCellValue(nomefoglioxsl, "I"+strconv.Itoa(riga), pulisciHTML(passo.risultatoatteso))
		/*		err = xlsx.SetCellStyle(nomefoglioxsl, "A"+strconv.Itoa(riga), "I"+strconv.Itoa(riga), `{"alignment": {
				"horizontal": "left",
				"shrink_to_fit": true,
				"vertical": "top",
				"wrap_text": true}}`)*/
		if err != nil {
			log.Logga(nomemodulo).Error(err)
			return err
		}
		return nil
	case "json":
		return nil
	case "testo":
		return nil
	default:
		log.Logga(nomemodulo).Error("Tipo file sconosciuto su scriviRiga")
		return nil
	}
}

//Chiudi file
func chiudiFile(nomefile string, tipo string) error {
	switch tipo {
	case "xsl":
		scriviFileXls(nomefile)
		return nil
	case "json":
		return nil
	case "testo":
		return nil
	default:
		log.Logga(nomemodulo).Error("Tipo file sconosciuto su chiudiFile")
		return nil
	}
}

func pulisciHTML(testo interface{}) string {
	var stringa string

	if testo != nil {
		stringa, err = html2text.FromString(testo.(string))
		if err != nil {
			log.Logga(nomemodulo).Error(err)
			return "**-n.d.-**"
		}
	} else {
		return "Vuoto"
	}
	return stringa
}
