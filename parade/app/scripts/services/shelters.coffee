angular.module('tdbApp')
.factory 'Shelters', ['$http',
(
  $http
) ->

  get = (shelterId) ->
    $http.get "http://api.tierheimdb.de/v1/shelter/#{shelterId}"

  getAll = ->
    $http.get("http://api.tierheimdb.de/v1/shelters")

  getStatesMap = ->
    getAll().then (response) ->
      map = {}
      if response.status == 200
        for shelter in response.data
          if map.hasOwnProperty(shelter.state)
            map[shelter.state].push(shelter)
          else
            map[shelter.state] = [shelter]

      return map

  {
    get
    getAll
    getStatesMap
  }
]
