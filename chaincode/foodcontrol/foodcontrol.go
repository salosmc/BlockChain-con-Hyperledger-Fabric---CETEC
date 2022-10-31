package main

import (
	//modulo para convertir json a bytes
	"encoding/json"
	//modulo para mostrar por pantalla
	"fmt"
	//modulo de hyperledger
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for control the food
type SmartContract struct {
	contractapi.Contract
}

//IMPORTANTE!!!! EL ACTIVO QUE QUEREMOS PERSISITIR

type Alumno struct {
	Nombre string `json:"nombre"`
	Apellido string `json:"apellido"`
	Padron int `json:"padron"`
	Materia string `json:"materia"`
	Asistencia bool `json:"asistencia"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, alumnoId string, nombre string, apellido string, padron int ,materia string, asistencia bool) error {

	//Validaciones de sintaxis

	//validaciones de negocio

	//validamos primero si el alumno ya existe en la blockchain
	resAlumno, err :=s.Query(ctx, alumnoId)
	if resAlumno != nil {
		fmt.Printf("Alumno already exist error: %s", err.Error())
		return err
	}

	alumno := Alumno{
		Nombre:  nombre,
		Apellido: apellido,
		Padron: padron,
		Materia: materia,
		Asistencia: asistencia,
	}

	//transformo alumno a bytes
	//lo resivo como json, opero en go como structura y lo persisito en bytes.
	alumnoAsBytes, err := json.Marshal(alumno)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	//PutState es el que nos permite guardar en el libro distribuido
	return ctx.GetStub().PutState(alumnoId, alumnoAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, alumnoId string) (*Alumno, error) {

	//busco el alumno por su id en la blockchain  
	alumnoAsBytes, err := ctx.GetStub().GetState(alumnoId)

	//validamos si tenemos un error
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	//validamos si no encontro el alumno
	if alumnoAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", alumnoId)
	}

	//declaramos una estructura Alumno
	alumno := new(Alumno)

	//convertimos de bytes a la estructura alumno
	err = json.Unmarshal(alumnoAsBytes, alumno)
	//validamos si existio algun error
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error. %s", err.Error())
	}

	//retornamos el alumno o null
	return alumno, nil
}


//Food describes basic details of what makes up a food
/*
type Food struct {
	Farmer  string `json:"farmer"`
	Variety string `json:"variety"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, foodId string, farmer string, variety string) error {

	//Validaciones de sintaxis

	//validaciones de negocio

	food := Food{
		Farmer:  farmer,
		Variety: variety,
	}

	foodAsBytes, err := json.Marshal(food)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	return ctx.GetStub().PutState(foodId, foodAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, foodId string) (*Food, error) {

	foodAsBytes, err := ctx.GetStub().GetState(foodId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if foodAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", foodId)
	}

	food := new(Food)

	err = json.Unmarshal(foodAsBytes, food)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error. %s", err.Error())
	}

	return food, nil
}
*/

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create foodcontrol chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting foodcontrol chaincode: %s", err.Error())
	}
}
