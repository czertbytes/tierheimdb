package main

import (
	"net/http"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func APIv1GetSheltersHandler(w http.ResponseWriter, r *http.Request) {
	shelters, err := pb.GetEnabledShelters()
	if err != nil {
		internalServerError(w, err)
		return
	}

	response(w, shelters)
}

func APIv1PostSheltersHandler(w http.ResponseWriter, r *http.Request) {
	var shelters []*pb.Shelter
	err := unmarshalRequestBody(r, &shelters)
	if err != nil {
		badRequest(w, err)
		return
	}

	if _, err := pb.PutShelters(shelters); err != nil {
		internalServerError(w, err)
		return
	}

	responseCreated(w, shelters)
}

func APIv1DeleteSheltersHandler(w http.ResponseWriter, r *http.Request) {
	if err := pb.DeleteEnabledShelters(); err != nil {
		internalServerError(w, err)
		return
	}

	responseNoContent(w)
}

func APIv1GetShelterHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]

	shelter, err := pb.GetShelter(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, shelter)
}

func APIv1DeleteShelterHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]

	if err := pb.DeleteShelter(shelterId); err != nil {
		internalServerError(w, err)
		return
	}

	responseNoContent(w)
}

func APIv1GetShelterAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)
	shelterId := routeParams["shelterId"]

	update, err := pb.GetLastUpdate(shelterId)
	if err != nil {
		internalServerError(w, err)
		return
	}

	var animals []pb.Animal
	switch r.URL.Query().Get("type") {
	case "cats":
		animals, err = pb.GetCats(shelterId, update.Id)
		if err != nil {
			internalServerError(w, err)
			return
		}
	case "dogs":
		animals, err = pb.GetDogs(shelterId, update.Id)
		if err != nil {
			internalServerError(w, err)
			return
		}
	default:
		animals, err = pb.GetAnimals(shelterId, update.Id)
		if err != nil {
			internalServerError(w, err)
			return
		}
	}

	response(w, animals)
}

func APIv1PostShelterAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]

	var animals []*pb.Animal
	err := unmarshalRequestBody(r, &animals)
	if err != nil {
		badRequest(w, err)
		return
	}

	update := pb.NewUpdate(shelterId)
	if err := pb.PutUpdate(update); err != nil {
		internalServerError(w, err)
		return
	}

	for _, a := range animals {
		a.UpdateId = update.Id
		a.ShelterId = shelterId
	}

	if _, err := pb.PutAnimals(animals); err != nil {
		internalServerError(w, err)
		return
	}

	responseCreated(w, animals)
}

func APIv1DeleteShelterUpdateAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	/*
		routeParams := mux.Vars(r)
		shelterId := routeParams["shelterId"]
		updateId := routeParams["updateId"]

			if err := pb.DeleteAnimals(shelterId, updateId); err != nil {
				internalServerError(w, err)
				return
			}
	*/

	responseNoContent(w)
}

func APIv1GetShelterUpdateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)
	shelterId := routeParams["shelterId"]
	updateId := routeParams["updateId"]
	animalId := routeParams["animalId"]

	animal, err := pb.GetAnimal(shelterId, updateId, animalId)
	if err != nil {
		internalServerError(w, err)
		return
	}

	response(w, animal)
}

func APIv1DeleteShelterUpdateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	/*
		routeParams := mux.Vars(r)
		shelterId := routeParams["shelterId"]
		updateId := routeParams["updateId"]
		animalId := routeParams["animalId"]

			if err := pb.DeleteAnimals(pb.NewAnimalSearch(c).SetId(animalId).SetShelterId(shelterId)); err != nil {
				internalServerError(w, err)
				return
			}
	*/

	responseNoContent(w)
}

func APIv1SyncShelterSourcesHandler(w http.ResponseWriter, r *http.Request) {
	//shelterId := mux.Vars(r)["shelterId"]

	//c := appengine.NewContext(r)

	responseNoContent(w)
}

func APIv1GetShelterUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]

	updates, err := pb.GetUpdates(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, updates)
}

func APIv1DeleteShelterUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	/*
		    shelterId := mux.Vars(r)["shelterId"]

				if err := pb.DeleteUpdates(pb.NewUpdateSearch(c).SetShelterId(shelterId)); err != nil {
					badRequest(w, err)
					return
				}
	*/

	responseNoContent(w)
}

func APIv1GetShelterLastUpdateHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]

	update, err := pb.GetLastUpdate(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, update)
}

func APIv1DeleteShelterLastUpdateHandler(w http.ResponseWriter, r *http.Request) {
	/*
		    shelterId := mux.Vars(r)["shelterId"]

				if err := pb.DeleteUpdates(pb.NewUpdateSearch(c).SetShelterId(shelterId)); err != nil {
					badRequest(w, err)
					return
				}
	*/

	responseNoContent(w)
}

func APIv1GetShelterUpdateHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shelterId := params["shelterId"]
	updateId := params["updateId"]

	update, err := pb.GetUpdate(shelterId, updateId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, update)
}

func APIv1DeleteShelterUpdateHandler(w http.ResponseWriter, r *http.Request) {
	/*
		params := mux.Vars(r)
		shelterId := params["shelterId"]
		updateId := params["updateId"]

			if err := pb.DeleteUpdates(pb.NewUpdateSearch(c).SetId(updateId).SetShelterId(shelterId)); err != nil {
				badRequest(w, err)
				return
			}
	*/

	responseNoContent(w)
}

func APIv1GetShelterUpdatesDiffHandler(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	//shelterId := params["shelterId"]
	//updateId := params["updateId"]
	//otherUpdateId := params["otherUpdateId"]

	animals := []*pb.Animal{}

	response(w, animals)
}

func APIv1GetShelterUpdateAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shelterId := params["shelterId"]
	updateId := params["updateId"]

	animals, err := pb.GetAnimals(shelterId, updateId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, animals)
}
