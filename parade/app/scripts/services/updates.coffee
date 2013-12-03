angular.module('tdbApp')
.factory 'Updates', ['$http', '$q',
(
  $http
  $q
) ->

  fetchLast = (shelterId) ->
    url = "http://api.tierheimdb.de/v1/shelter/#{shelterId}/updates?limit=1"
    $http.get(url)

  getLast = (shelterIds) ->
    $q.all((fetchLast(shelterId) for shelterId in shelterIds))

  {
    getLast
  }
]
