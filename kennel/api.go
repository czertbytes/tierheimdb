package main

import (
	"net/http"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func APIv1GetAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	animals, err := pb.SearchAnimals(parseAnimalsHandlerQuery(r))
	if err != nil {
		internalServerError(w, err)
		return
	}

	response(w, animals)
}

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
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	shelter, err := pb.GetShelter(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, shelter)
}

func APIv1DeleteShelterHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	if err := pb.DeleteShelter(shelterId); err != nil {
		internalServerError(w, err)
		return
	}

	responseNoContent(w)
}

func APIv1GetShelterAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	update, err := pb.GetLastUpdate(shelterId)
	if err != nil {
		internalServerError(w, err)
		return
	}

	animals, err := pb.GetAnimals(shelterId, update.Id, r.URL.Query().Get("type"))
	if err != nil {
		internalServerError(w, err)
		return
	}

	response(w, animals)
}

func APIv1PostShelterAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	var animals []*pb.Animal
	if err := unmarshalRequestBody(r, &animals); err != nil {
		badRequest(w, err)
		return
	}

	u := pb.NewUpdate(shelterId)
	if err := pb.PutUpdate(u); err != nil {
		internalServerError(w, err)
		return
	}

	for _, a := range animals {
		a.UpdateId = u.Id
		a.ShelterId = shelterId
	}

	if _, err := pb.PutAnimals(animals); err != nil {
		internalServerError(w, err)
		return
	}

	responseCreated(w, animals)
}

func APIv1DeleteShelterUpdateAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, updateId, err := shelterIdUpdateIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	if err := pb.DeleteAnimals(shelterId, updateId); err != nil {
		internalServerError(w, err)
		return
	}

	responseNoContent(w)
}

func APIv1GetShelterUpdateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, updateId, animalId, err := shelterIdUpdateIdAnimalIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	a, err := pb.GetAnimal(shelterId, updateId, animalId)
	if err != nil {
		internalServerError(w, err)
		return
	}

	response(w, a)
}

func APIv1DeleteShelterUpdateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, updateId, animalId, err := shelterIdUpdateIdAnimalIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	if err := pb.DeleteAnimal(shelterId, updateId, animalId); err != nil {
		internalServerError(w, err)
		return
	}

	responseNoContent(w)
}

func APIv1SyncShelterSourcesHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	animals, err := pb.RunShelterSync(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	responseCreated(w, animals)
}

func APIv1GetShelterUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	updates, err := pb.GetUpdates(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, updates)
}

func APIv1DeleteShelterUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	if err := pb.DeleteUpdates(shelterId); err != nil {
		badRequest(w, err)
		return
	}

	responseNoContent(w)
}

func APIv1GetShelterLastUpdateHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	u, err := pb.GetLastUpdate(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, u)
}

func APIv1DeleteShelterLastUpdateHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	u, err := pb.GetLastUpdate(shelterId)
	if err != nil {
		badRequest(w, err)
		return
	}

	if err := pb.DeleteUpdate(shelterId, u.Id); err != nil {
		badRequest(w, err)
		return
	}

	responseNoContent(w)
}

func APIv1GetShelterUpdateHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, updateId, err := shelterIdUpdateIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	update, err := pb.GetUpdate(shelterId, updateId)
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, update)
}

func APIv1DeleteShelterUpdateHandler(w http.ResponseWriter, r *http.Request) {
	shelterId, updateId, err := shelterIdUpdateIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	if err := pb.DeleteUpdate(shelterId, updateId); err != nil {
		badRequest(w, err)
		return
	}

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
	shelterId, updateId, err := shelterIdUpdateIdFromRequest(r)
	if err != nil {
		badRequest(w, err)
		return
	}

	animals, err := pb.GetAnimals(shelterId, updateId, r.URL.Query().Get("type"))
	if err != nil {
		badRequest(w, err)
		return
	}

	response(w, animals)
}
