package drtestlink

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/Luxurioust/excelize"
	"github.com/kolo/xmlrpc"

	log "drollo.it/drlog"
)

var client *xmlrpc.Client
var elencoTC testcase
var elencoSuites []suite
var elencoidTC []string
var elencoTCs testcases
var primasuite bool

//Collegamento a TestLink
func collegaTestLink(svr string) error {
	log.Logga(nomemodulo).Info("Mi collego a testlink")
	client, err = xmlrpc.NewClient(svr, nil)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}
	log.Logga(nomemodulo).Info("Connessione OK")

	xlsx = excelize.CreateFile()

	return nil
}

// recupera l'id dal nome del progetto
func recuperaIDProjectNew(keysec string, nomeprogetto string) (string, error) {
	var result interface{}
	var risultatoIDProject map[string]interface{}
	var idPrj string
	var ok bool

	log.Logga(nomemodulo).Debug("Predispongo chiamata recuperaIDProject su progetto:" + nomeprogetto)
	parTL := make(map[string]string)
	parTL["devKey"] = keysec
	parTL["testprojectname"] = nomeprogetto
	log.Logga(nomemodulo).Debug("Chiamo metodo tl.getTestProjectByName con key:" + keysec + " progetto: " + nomeprogetto)
	err = client.Call("tl.getTestProjectByName", parTL, &result)
	if err != nil {
		log.Logga(nomemodulo).Error("Errore in getTestProjectByName")
		return "", err
	}
	log.Logga(nomemodulo).Debug("Chiamata getTestProjectByName riuscita")
	log.Logga(nomemodulo).Debug(result)
	if risultatoIDProject, ok = result.(map[string]interface{}); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(result).Type())
		return "", errors.New("Errore in mapping risultati getTestProjectByName")
	}
	if idPrj, ok = risultatoIDProject["id"].(string); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(risultatoIDProject["id"]).Type())
		return "0", errors.New("Errore in mapping id")
	}
	log.Logga(nomemodulo).Debug(idPrj)

	return idPrj, nil
}

// recupera i TC dalla testsuite
func recuperaTCfromTS(keysec string, idtestsuite string) ([]string, error) {
	var result interface{}
	var risultatoTC []interface{}
	var ok bool
	var risultatoTCdettaglio map[string]interface{}
	var elencotestcases []string
	var idtc string

	log.Logga(nomemodulo).Debug("Predispongo chiamata getTestCasesForTestSuite  idtestsuite:" + idtestsuite)

	parTL := make(map[string]string)
	parTL["devKey"] = keysec
	parTL["testsuiteid"] = idtestsuite
	parTL["deep"] = "true"
	parTL["details"] = "full"

	log.Logga(nomemodulo).Debug("Chiamo metodo tl.getTestCasesForTestSuite con idtestsuite:" + idtestsuite)
	err = client.Call("tl.getTestCasesForTestSuite", parTL, &result)
	if err != nil {
		log.Logga(nomemodulo).Error("Errore in getTestCasesForTestSuite")
		return nil, err
	}
	log.Logga(nomemodulo).Debug("Chiamata getTestCasesForTestSuite riuscita")

	//Verifico risultato
	if risultatoTC, ok = result.([]interface{}); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(result).Type())
		return nil, errors.New("Errore in mapping risultati recuperaTC su TC")
	}
	//Leggo tutti i TC
	log.Logga(nomemodulo).Debug("Ciclo tutti i TC per il progetto")
	for i, value := range risultatoTC {
		log.Logga(nomemodulo).Debug("Ciclo TC n°: " + strconv.Itoa(i))

		if risultatoTCdettaglio, ok = value.(map[string]interface{}); !ok {
			log.Logga(nomemodulo).Debug(reflect.ValueOf(value).Type())
			return nil, errors.New("Errore in mapping risultati risultatoTSs")
		}
		if idtc, ok = risultatoTCdettaglio["id"].(string); !ok {
			log.Logga(nomemodulo).Error("Errore in lettura id")
			return nil, errors.New("Errore in mapping risultato singola suite id")
		}
		log.Logga(nomemodulo).Debug(idtc)
		elencotestcases = append(elencotestcases, idtc)
		log.Logga(nomemodulo).Debug("Numero TC:" + strconv.Itoa(len(elencoSuites)))
	}
	//recupero il nome della TS
	log.Logga(nomemodulo).Debug("-------fine TCdaTS-------")

	return elencotestcases, nil
}

// recupera primo TS dal progetto specifico
func recuperaTSfromProjectNew(keysec string, idprogetto string) error {
	var result interface{}
	var risultatoTSs []interface{}
	var risultatoTS map[string]interface{}
	var ok bool
	var miasuite suite

	log.Logga(nomemodulo).Debug("Predispongo chiamata recuperaTSfromProjectNew su idprogetto:" + idprogetto)
	parTL := make(map[string]string)
	parTL["devKey"] = keysec
	parTL["testprojectid"] = idprogetto
	log.Logga(nomemodulo).Debug("Chiamo metodo tl.getFirstLevelTestSuitesForTestProject con key:" + keysec + " idprogetto: " + idprogetto)
	err = client.Call("tl.getFirstLevelTestSuitesForTestProject", parTL, &result)
	if err != nil {
		log.Logga(nomemodulo).Error(err)
		return err
	}

	log.Logga(nomemodulo).Debug("Chiamata getFirstLevelTestSuitesForTestProject riuscita")
	log.Logga(nomemodulo).Debug(result)
	if risultatoTSs, ok = result.([]interface{}); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(result).Type())
		return errors.New("Errore in mapping risultati  risultatoIDsuite")
	}
	//Leggo tutti i TS
	log.Logga(nomemodulo).Debug("Ciclo tutti i TS per il progetto")
	for i, value := range risultatoTSs {
		log.Logga(nomemodulo).Debug("Ciclo TS n°: " + strconv.Itoa(i))

		if risultatoTS, ok = value.(map[string]interface{}); !ok {
			log.Logga(nomemodulo).Debug(reflect.ValueOf(value).Type())
			return errors.New("Errore in mapping risultati risultatoTSs")
		}
		if _, ok = risultatoTS["id"].(string); !ok {
			log.Logga(nomemodulo).Error("Errore in lettura id")
			return errors.New("Errore in mapping risultato singola suite id")
		}
		log.Logga(nomemodulo).Debug(risultatoTS["id"].(string))
		miasuite.id = risultatoTS["id"].(string)

		if _, ok = risultatoTS["name"].(string); !ok {
			log.Logga(nomemodulo).Error("Errore in lettura name")
			return errors.New("Errore in mapping risultato singola suite id")
		}
		log.Logga(nomemodulo).Debug(risultatoTS["name"].(string))
		miasuite.nome = risultatoTS["name"].(string)

		elencoSuites = append(elencoSuites, miasuite)
		log.Logga(nomemodulo).Debug("Numero Suites:" + strconv.Itoa(len(elencoSuites)))
	}
	return nil
}

// recupera id plan da nome plan
func recuperaIDPlanbyName(keysec string, nomeprogetto string, nomeplan string) (string, error) {
	var result interface{}
	var risultatoTSs []interface{}
	var risultatoTS map[string]interface{}
	var idTS string
	var ok bool

	log.Logga(nomemodulo).Info("Predispongo chiamata getTestPlanByName  su progetto-plan:" + nomeprogetto + "-" + nomeplan)

	parTL := make(map[string]string)
	parTL["devKey"] = keysec
	parTL["testprojectname"] = nomeprogetto
	parTL["testplanname"] = nomeplan

	log.Logga(nomemodulo).Debug("Chiamo metodo tl.getTestPlanByName")
	err = client.Call("tl.getTestPlanByName", parTL, &result)
	if err != nil {
		log.Logga(nomemodulo).Error("Errore in getTestPlanByName")
		return "0", err
	}
	log.Logga(nomemodulo).Debug("Chiamata getTestPlanByName riuscita")
	log.Logga(nomemodulo).Debug(result)
	if risultatoTSs, ok = result.([]interface{}); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(result).Type())
		return "0", errors.New("Errore in mapping risultati getTestPlanByName")
	}
	for _, valorizziamo := range risultatoTSs {
		log.Logga(nomemodulo).Debug("Ciclo i TS")
		if risultatoTS, ok = valorizziamo.(map[string]interface{}); !ok {
			log.Logga(nomemodulo).Debug(reflect.ValueOf(valorizziamo).Type())
			return "", errors.New("Errore in mapping risultati risultatoTSs su singolo item")
		}
		if idTS, ok = risultatoTS["id"].(string); !ok {
			log.Logga(nomemodulo).Debug(reflect.ValueOf(risultatoTS["id"]).Type())
			return "", errors.New("Errore in mapping idts")
		}
		log.Logga(nomemodulo).Debug(risultatoTS["id"])
	}
	log.Logga(nomemodulo).Debug(risultatoTS)
	return idTS, nil
}

// recupera i TC dal plan e li valorizza in elencoTC
func recuperaTCdaPlan(keysec string, idplan string) error {
	var result interface{}
	var ok bool
	var risultatoTCs map[string]interface{}
	var risultatoTC []interface{}
	var risultatoTCdettaglio map[string]interface{}

	log.Logga(nomemodulo).Info("Predispongo chiamata getTestCasesForTestPlan  su idplan:" + idplan)

	parTL := make(map[string]string)
	parTL["devKey"] = keysec
	parTL["testplanid"] = idplan

	log.Logga(nomemodulo).Debug("Chiamo metodo tl.getTestCasesForTestPlan")
	err = client.Call("tl.getTestCasesForTestPlan", parTL, &result)
	if err != nil {
		log.Logga(nomemodulo).Error("Errore in getTestCasesForTestPlan")
		return err
	}
	log.Logga(nomemodulo).Debug("Chiamata getTestCasesForTestPlan riuscita")

	if risultatoTCs, ok = result.(map[string]interface{}); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(result).Type())
		return errors.New("Errore in mapping risultati recuperaTCdaPlan su TCS")
	}
	//Leggo tutti i TC
	for i, value := range risultatoTCs {
		log.Logga(nomemodulo).Debug(i)
		log.Logga(nomemodulo).Debug(value)
		log.Logga(nomemodulo).Debug(reflect.ValueOf(value).Type())
		if risultatoTC, ok = value.([]interface{}); !ok {
			log.Logga(nomemodulo).Debug(reflect.ValueOf(value).Type())
			return errors.New("Errore in mapping risultati recuperaTCdaPlan su TC")
		}
		//Leggo le platform + i dettagli
		for _, valorizziamo := range risultatoTC {
			if risultatoTCdettaglio, ok = valorizziamo.(map[string]interface{}); !ok {
				log.Logga(nomemodulo).Debug(reflect.ValueOf(valorizziamo).Type())
				return errors.New("Errore in mapping risultati recuperaTCdaPlan su TCdettaglio")
			}
			if _, ok = risultatoTCdettaglio["tc_id"].(string); !ok {
				return errors.New("Errore in lettura tc_id di recuperaTCdaPlan su TCdettaglio")
			}
			log.Logga(nomemodulo).Debug(risultatoTCdettaglio["tcase_name"].(string))
			elencoidTC = append(elencoidTC, risultatoTCdettaglio["tc_id"].(string))
		}
	}
	log.Logga(nomemodulo).Debug(elencoTC)
	return nil
}

// recupera nome suite da idsuite
func recuperaTSnameFromID(keysec string, idTS string) (string, error) {
	var result interface{}
	var ok bool
	var risultatoTSs map[string]interface{}
	var ritornaTSname string

	log.Logga(nomemodulo).Info("Predispongo chiamata recuperaTSnameFromID  su idsuite:" + idTS)

	parTL := make(map[string]string)
	parTL["devKey"] = keysec
	parTL["testsuiteid"] = idTS

	log.Logga(nomemodulo).Debug("Chiamo metodo tl.getTestSuiteByID")
	err = client.Call("tl.getTestSuiteByID", parTL, &result)
	if err != nil {
		log.Logga(nomemodulo).Error("Errore in getTestSuiteByID")
		return "0", err
	}
	log.Logga(nomemodulo).Debug("Chiamata getTestSuiteByID riuscita")

	log.Logga(nomemodulo).Debug(result)
	if risultatoTSs, ok = result.(map[string]interface{}); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(result).Type())
		return "0", errors.New("Errore in mapping risultati getTestSuiteByID")
	}
	ritornaTSname = risultatoTSs["name"].(string)
	log.Logga(nomemodulo).Debug("Nome suite:" + ritornaTSname)
	return ritornaTSname, nil
}

//recupera singolo testcase da id del test case
func recuperaTC(keysec string, idTC string) (testcase, error) {
	var result interface{}
	var risultatoTC []interface{}
	var ok bool
	var risultatoTCdettaglio map[string]interface{}
	var risultatoTCsteps []interface{}
	var risultatoTCstepsDettaglio map[string]interface{}
	var elenchiamoTC testcase
	var nomeSuite string
	var passi []passotestcase

	log.Logga(nomemodulo).Info("Predispongo chiamata recuperaTC  su idcase:" + idTC)
	parTL := make(map[string]string)
	parTL["devKey"] = keysec
	parTL["testcaseid"] = idTC
	log.Logga(nomemodulo).Debug("Chiamo metodo tl.getTestCase")
	err = client.Call("tl.getTestCase", parTL, &result)
	if err != nil {
		log.Logga(nomemodulo).Error("Errore in getTestCase")
		return elenchiamoTC, err
	}
	log.Logga(nomemodulo).Debug("Chiamata getTestCase riuscita")

	//Verifico risultato
	if risultatoTC, ok = result.([]interface{}); !ok {
		log.Logga(nomemodulo).Debug(reflect.ValueOf(result).Type())
		return elenchiamoTC, errors.New("Errore in mapping risultati recuperaTC su TC")
	}
	//Leggo i TC
	log.Logga(nomemodulo).Debug("-------inizio TC-------")
	log.Logga(nomemodulo).Debug(result)
	for _, valorizziamo := range risultatoTC {
		log.Logga(nomemodulo).Debug("Ciclo i TC n° " + strconv.Itoa(len(risultatoTC)))

		if risultatoTCdettaglio, ok = valorizziamo.(map[string]interface{}); !ok {
			log.Logga(nomemodulo).Debug(reflect.ValueOf(valorizziamo).Type())
			return elenchiamoTC, errors.New("Errore in mapping risultati recuperaTC su TCdettaglio")
		}
		if myvalore, ok := risultatoTCdettaglio["id"].(string); ok {
			elenchiamoTC.idTC = myvalore
		}
		if myvalore, ok := risultatoTCdettaglio["testsuite_id"].(string); ok {
			elenchiamoTC.idTS = myvalore
		}
		if myvalore, ok := risultatoTCdettaglio["summary"].(string); ok {
			elenchiamoTC.sommario = myvalore
		}
		if myvalore, ok := risultatoTCdettaglio["preconditions"].(string); ok {
			elenchiamoTC.precondizioni = myvalore
		}
		if myvalore, ok := risultatoTCdettaglio["name"].(string); ok {
			if myvalore2, ok := risultatoTCdettaglio["full_tc_external_id"].(string); ok {
				elenchiamoTC.nome = myvalore2 + ":" + myvalore
			}
		}
		//cerco se ci sono gli steps
		if risultatoTCsteps, ok = risultatoTCdettaglio["steps"].([]interface{}); !ok {
			log.Logga(nomemodulo).Error("Errore in mapping risultati recuperaTC su TCdettaglioSteps")
			passi = append(passi, passotestcase{"", "", ""})
			elenchiamoTC.elencopassi = passi
			//log.Logga(nomemodulo).Debug(reflect.ValueOf(risultatoTCdettaglio["steps"]).Type())
		} else {
			// Trovati steps
			log.Logga(nomemodulo).Debug("Ciclo gli steps. N°: " + strconv.Itoa(len(risultatoTCsteps)))

			for _, valorestep := range risultatoTCsteps {
				var miostep passotestcase
				// Trovato singolo step
				log.Logga(nomemodulo).Debug("Trovato step")
				if risultatoTCstepsDettaglio, ok = valorestep.(map[string]interface{}); !ok {
					log.Logga(nomemodulo).Debug(reflect.ValueOf(valorestep).Type())
					log.Logga(nomemodulo).Error("Errore in mapping risultati recuperaTC su singolo step")
				} else {

					if myvalore, ok := risultatoTCstepsDettaglio["expected_results"].(string); ok {
						miostep.risultatoatteso = myvalore
					}
					if myvalore, ok := risultatoTCstepsDettaglio["step_number"].(string); ok {
						miostep.numerostep = myvalore
					}
					if myvalore, ok := risultatoTCstepsDettaglio["actions"].(string); ok {
						miostep.azione = myvalore
					}
					passi = append(passi, miostep)
				}
			}
			elenchiamoTC.elencopassi = passi
		}
		//recupero il nome della TS
		nomeSuite, err = recuperaTSnameFromID(keysec, elenchiamoTC.idTS)
		if err != nil {
			log.Logga(nomemodulo).Error(err)
			return elenchiamoTC, err
		}
		elenchiamoTC.nomeTS = nomeSuite
		log.Logga(nomemodulo).Debug("-------fine TC-------")
	}
	return elenchiamoTC, nil
}
