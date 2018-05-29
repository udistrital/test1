package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/financiera_mongo_crud/db"
	"github.com/udistrital/financiera_mongo_crud/models"
	_ "gopkg.in/mgo.v2"
)

// Operaciones Crud ArbolRubros
type ArbolRubrosController struct {
	beego.Controller
}

// @Title GetAll
// @Description get all objects
// @Success 200 ArbolRubros models.ArbolRubros
// @Failure 403 :objectId is empty
// @router / [get]
func (j *ArbolRubrosController) GetAll() {
	session, _ := db.GetSession()

	var query = make(map[string]interface{})

	if v := j.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				j.Data["json"] = errors.New("Consulta invalida")
				j.ServeJSON()
				return
			}

			if i, err := strconv.Atoi(kv[1]); err == nil {
				k, v := kv[0], i
				query[k] = v
			} else {
				k, v := kv[0], kv[1]
				query[k] = v
			}
		}
	}

	obs := models.GetAllArbolRubross(session, query)

	if len(obs) == 0 {
		j.Data["json"] = []string{}
	} else {
		j.Data["json"] = &obs
	}

	j.ServeJSON()
}

// @Title Get
// @Description get ArbolRubros by nombre
// @Param	nombre		path 	string	true		"El nombre de la ArbolRubros a consultar"
// @Success 200 {object} models.ArbolRubros
// @Failure 403 :uid is empty
// @router /:id [get]
func (j *ArbolRubrosController) Get() {
	id := j.GetString(":id")
	session, _ := db.GetSession()
	if id != "" {
		arbolrubros, err := models.GetArbolRubrosById(session, id)
		if err != nil {
			j.Data["json"] = err.Error()
		} else {
			j.Data["json"] = arbolrubros
		}
	}
	j.ServeJSON()
}

// @Title Borrar ArbolRubros
// @Description Borrar ArbolRubros
// @Param	objectId		path 	string	true		"El ObjectId del objeto que se quiere borrar"
// @Success 200 {string} ok
// @Failure 403 objectId is empty
// @router /:objectId [delete]
func (j *ArbolRubrosController) Delete() {
	session, _ := db.GetSession()
	objectId := j.Ctx.Input.Param(":objectId")
	result, _ := models.DeleteArbolRubrosById(session, objectId)
	j.Data["json"] = result
	j.ServeJSON()
}

// @Title Crear ArbolRubros
// @Description Crear ArbolRubros
// @Param	body		body 	models.ArbolRubros	true		"Body para la creacion de ArbolRubros"
// @Success 200 {int} ArbolRubros.Id
// @Failure 403 body is empty
// @router / [post]
func (j *ArbolRubrosController) Post() {
	var arbolrubros models.ArbolRubros
	json.Unmarshal(j.Ctx.Input.RequestBody, &arbolrubros)
	fmt.Println(arbolrubros)
	session, _ := db.GetSession()
	models.InsertArbolRubros(session, arbolrubros)
	j.Data["json"] = "insert success!"
	j.ServeJSON()
}

// @Title Update
// @Description update the ArbolRubros
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (j *ArbolRubrosController) Put() {
	objectId := j.Ctx.Input.Param(":objectId")

	var arbolrubros models.ArbolRubros
	json.Unmarshal(j.Ctx.Input.RequestBody, &arbolrubros)
	session, _ := db.GetSession()

	err := models.UpdateArbolRubros(session, arbolrubros, objectId)
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = "update success!"
	}
	j.ServeJSON()
}

// @Title Preflight options
// @Description Crear ArbolRubros
// @Param	body		body 	models.ArbolRubros	true		"Body para la creacion de ArbolRubros"
// @Success 200 {int} ArbolRubros.Id
// @Failure 403 body is empty
// @router / [options]
func (j *ArbolRubrosController) Options() {
	j.Data["json"] = "success!"
	j.ServeJSON()
}

// @Title Preflight options
// @Description Crear ArbolRubros
// @Param	body		body 	models.ArbolRubros true		"Body para la creacion de ArbolRubros"
// @Success 200 {int} ArbolRubros.Id
// @Failure 403 body is empty
// @router /:objectId [options]
func (j *ArbolRubrosController) ArbolRubrosDeleteOptions() {
	j.Data["json"] = "success!"
	j.ServeJSON()
}

// @Title Registra rubro
// @Description Convierte la estructura del api mid para registrarla en un documento mongo
// @Param	body		interface{}	true		"Body para la creacion de un nuevo rubro"
// @Success 201 {string} recibido
// @Failure 403 body is empty
// @router /registrarRubro [post]
func (j *ArbolRubrosController) RegistrarRubro() {
	try.This(func() {
		var (
			rubroData  interface{}
			rubroPadre string = ""
			err        error
		)
		session, _ := db.GetSession()
		json.Unmarshal(j.Ctx.Input.RequestBody, &rubroData)
		rubroDataHijo := rubroData.(map[string]interface{})["RubroHijo"].(map[string]interface{})

		nuevoRubro := models.ArbolRubros{
			Id:          rubroDataHijo["Codigo"].(string),
			Idpsql:      "-1",
			Nombre:      rubroDataHijo["Nombre"].(string),
			Descripcion: rubroDataHijo["Descripcion"].(string),
			Hijos:       nil}

		if rubroDataPadre := rubroData.(map[string]interface{})["RubroPadre"].(map[string]interface{}); rubroDataPadre["Codigo"] != nil {
			rubroPadre = rubroDataPadre["Codigo"].(string)
			nuevoRubro.Padre = rubroPadre
			updatedRubro, _ := models.GetArbolRubrosById(session, rubroPadre)
			updatedRubro.Hijos = append(updatedRubro.Hijos, rubroDataHijo["Codigo"].(string))
			session, _ = db.GetSession()
			err = models.RegistrarRubroTransacton(updatedRubro, nuevoRubro, session)
		} else {
			err = models.InsertArbolRubros(session, nuevoRubro)
		}

		if err != nil {
			panic(err.Error())
		} else {
			j.Data["json"] = map[string]interface{}{"Type": "sucess"}
		}
	}).Catch(func(e try.E) {
		beego.Error(e)
		j.Data["json"] = map[string]interface{}{"Type": "error"}
	})

	j.ServeJSON()
}

// @Title Eliminar rubro
// @Description recibe el idPsql del rubro desde api mid para eliminar el rubro
// @Param	body		interface{}	true		"Body para la eliminación de un rubro"
// @Success 201 {string} sucess
// @Failure 403 body is empty
// @router /eliminarRubro/:idPsql [get]
func (j *ArbolRubrosController) EliminarRubro() {
	try.This(func() {
		session, _ := db.GetSession()
		var (
			// rubroIdPsql interface{}
			err error
		)
		rubroIdPsql := j.GetString(":idPsql")
		rubroHijo, _ := models.GetArbolRubrosByIdPsql(session, rubroIdPsql)
		beego.Info("rubroHijo: ", rubroHijo)
		session, _ = db.GetSession()
		rubroPadre, _ := models.GetArbolRubrosById(session, rubroHijo.Padre)
		beego.Info("rubroPadre sin modificar: ", rubroPadre)
		rubroPadre.Hijos = remove(rubroPadre.Hijos, rubroHijo.Id)
		beego.Info("rubroPadre modificado: ", rubroPadre)
		session, _ = db.GetSession()
		err = models.EliminarRubroTransaccion(rubroPadre, rubroHijo, session)

		if err != nil {
			panic(err.Error())
		} else {
			j.Data["json"] = map[string]interface{}{"Type": "sucess"}
		}

	}).Catch(func(e try.E) {
		beego.Error(e)
		j.Data["json"] = map[string]interface{}{"Type": "error"}
	})
}

func remove(slice []string, object string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] == object {
			slice = append(slice[:i], slice[i+1:]...)
			return slice
		}
	}
	return slice
}

// @Title RaicesArbol
// @Description RaicesArbol
// @Param body body models.Rubro true "Body para la creacion de Rubro"
// @Success 200 {object} models.Object
// @Failure 404 body is empty
// @router /RaicesArbol [get]
func (j *ArbolRubrosController) RaicesArbol() {
	session, _ := db.GetSession()
	rubros, err := models.GetRaices(session)
	if err != nil {
		j.Data["json"] = err
	} else {
		j.Data["json"] = rubros
	}

	j.ServeJSON()
}

// @Title Preflight options
// @Description Construye el árbol a un nivel dependiendo de la raíz
// @Param body body stringtrue "Código de la raíz"
// @Success 200 {object} models.Object
// @Failure 404 body is empty
// @router /ArbolRubro/:raiz [get]
func (j *ArbolRubrosController) ArbolRubro() {
	nodoRaiz := j.GetString(":raiz")
	session, _ := db.GetSession()
	var arbolRubrosGrande []map[string]interface{}
	raiz, err := models.GetNodo(session, nodoRaiz)
	if err == nil {

		arbolRubros := make(map[string]interface{})
		arbolRubros["Id"], _ = strconv.Atoi(raiz.Idpsql)
		arbolRubros["Codigo"] = raiz.Id
		arbolRubros["Nombre"] = raiz.Nombre
		var hijos []interface{}
		for j := 0; j < len(raiz.Hijos); j++ {
			hijo := GetHijoRubro(raiz.Hijos[j])
			if len(hijo) > 0 {
				hijos = append(hijos, hijo)
			}
		}
		arbolRubros["Hijos"] = hijos
		arbolRubrosGrande = append(arbolRubrosGrande, arbolRubros)

		j.Data["json"] = arbolRubrosGrande
	} else {
		j.Data["json"] = err
	}

	j.ServeJSON()
}

func GetHijoRubro(id string) map[string]interface{} {
	session, _ := db.GetSession()
	rubroHijo, _ := models.GetNodo(session, id)
	hijo := make(map[string]interface{})

	if rubroHijo.Id != "" {
		hijo["Id"], _ = strconv.Atoi(rubroHijo.Idpsql)
		hijo["Codigo"] = rubroHijo.Id
		hijo["Nombre"] = rubroHijo.Nombre
		if len(rubroHijo.Hijos) == 0 {
			hijo["Hijos"] = nil
			return hijo
		}
	}
	return hijo
}
