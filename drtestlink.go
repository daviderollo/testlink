package drtestlink

import (
	"strconv"

	"sort"

	log "drollo.it/drlog"
)

const nomemodulo = "drtestlink"

var err error

//ReportTestLink esegui report di testlink
func ReportTestLink(tiporeport string, svr string, nomefile string, chiave string, nomeprogetto string, tipoOut string, nomeplan string, nomebuild string) error {
	err = collegaTestLink(svr)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	inizializzaFile(nomefile, tipoOut)
	switch tiporeport {
	case "progetto":
		err = casoProgetto(svr, nomefile, chiave, nomeprogetto, tipoOut)
		if err != nil {
			log.Logga(nomemodulo).Error(err)
			return err
		}
	case "plan":
		err = casoPlan(svr, nomefile, chiave, nomeprogetto, tipoOut, nomeplan, nomebuild)
		if err != nil {
			log.Logga(nomemodulo).Error(err)
			return err
		}
	}
	err = chiudiFile(nomefile, tipoOut)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	return nil
}

//CasoProgetto gestisce tutti i TC in caso di progetto
func casoProgetto(svr string, nomefile string, chiave string, nomeprogetto string, tipoOut string) error {
	var idprogetto string
	var elencoTC []string
	var mytc testcase

	//recupero l'id del progetto
	log.Logga(nomemodulo).Info("Predispongo recupero ID progetto:" + nomeprogetto)
	idprogetto, err = recuperaIDProjectNew(chiave, nomeprogetto)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	log.Logga(nomemodulo).Debug("Recuperato il projectID: " + idprogetto)

	//Leggo tutti i  TS del progetto e li metto in elencoSuites
	log.Logga(nomemodulo).Info("Recupero le TS del progetto")
	err = recuperaTSfromProjectNew(chiave, idprogetto)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	// Ciclo i testsuites
	log.Logga(nomemodulo).Debug("Recuperate le test suite: n° " + strconv.Itoa(len(elencoSuites)))
	for i, myTS := range elencoSuites {
		log.Logga(nomemodulo).Info("Ciclo tutti i TS trovati. Numero: " + strconv.Itoa(i))

		// recupero tutti i TC dalla TS
		elencoTC, err = recuperaTCfromTS(chiave, myTS.id)
		if err != nil {
			log.Logga(nomemodulo).Error(err)
			return err
		}
		//Ciclo tutti i TC della TS
		log.Logga(nomemodulo).Debug("Recuperati i TC dalla test suite")
		for _, myTC := range elencoTC {
			log.Logga(nomemodulo).Debug("Richiedo TC id:" + myTC)
			mytc, err = recuperaTC(chiave, myTC)
			if err != nil {
				log.Logga(nomemodulo).Error(err)
				return err
			}
			elencoTCs = append(elencoTCs, mytc)
		}
	}
	//Ordino i testcases
	log.Logga(nomemodulo).Info("Ordino i TC")
	sort.Sort(elencoTCs)
	log.Logga(nomemodulo).Info("Elenco TC: " + strconv.Itoa(len(elencoTCs)))

	//Scrivo su file
	log.Logga(nomemodulo).Info("Scrivo i risultati")
	err = scriviStructToFile(tipoOut)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	log.Logga(nomemodulo).Info("Fine caso Progetto")
	return nil
}

//casoPlan gestisce tutti i dati in caso di progetto
func casoPlan(svr string, nomefile string, chiave string, nomeprogetto string, tipoOut string, nomeplan string, nomebuild string) error {
	var idPlan string

	//recupero l'id del progetto
	log.Logga(nomemodulo).Info("Predispongo recupero ID su progetto:" + nomeprogetto)
	idPlan, err = recuperaIDPlanbyName(chiave, nomeprogetto, nomeplan)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	log.Logga(nomemodulo).Info("Recuperato il planID: " + idPlan)

	//recupero TC da plan
	log.Logga(nomemodulo).Info("Recupero prima TS del progetto")
	err = recuperaTCdaPlan(chiave, idPlan)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}

	log.Logga(nomemodulo).Info("Ciclo tutti i TC trovati")
	//Ciclo tutti i TC trovati
	for _, myTS := range elencoidTC {
		var mycasoditest testcase

		//recupero i TC dalla TS
		mycasoditest, err = recuperaTC(chiave, myTS)
		if err != nil {
			log.Logga(nomemodulo).Error(err)
			return err
		}
		elencoTCs = append(elencoTCs, mycasoditest)
	}
	//Ordino i testcases
	sort.Sort(elencoTCs)

	//Scrivo su file
	log.Logga(nomemodulo).Info("Scrivo i risultati")
	err = scriviStructToFile(tipoOut)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	log.Logga(nomemodulo).Info("Fine casoPlan")
	return nil
}
