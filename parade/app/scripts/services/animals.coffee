angular.module('tdbApp')
.factory 'Animals', ['$http', '$q',
(
  $http
  $q
) ->

  get = (shelterId, updateId, animalId) ->
    url = "http://api.tierheimdb.de/v1/shelter/#{shelterId}"
    url += "/update/#{updateId}/animal/#{animalId}"
    $http.get(url)

  getAll= (shelterIds, type, offset=0) ->
    $q.all(fetchLast(shelterId, type, offset) for shelterId in shelterIds)

  fetchLast = (shelterId, type, offset) ->
    url = "http://api.tierheimdb.de/v1/shelter/#{shelterId}/animals"
    url += "?offset=#{offset}&limit=20&type=#{type}"
    $http.get(url)

  {
    get
    getAll
  }
]
