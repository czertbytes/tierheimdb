angular.module('tdbApp')
.controller 'AnimalCtrl', ['$scope', '$routeParams', 'Animals'
(
  $scope
  $routeParams
  Animals
) ->

  $scope.animal = {}

  $scope.init = ->
    shelterId = $routeParams.shelterId
    updateId = $routeParams.updateId
    animalId = $routeParams.animalId
    Animals.get(shelterId, updateId, animalId)
      .then setAnimal

  setAnimal = (response) ->
    if response.status == 200
      $scope.animal = response.data
    else
      msg = "Getting animal detail failed! Status: #{response.status}"
      msg += "Error: #{response.data}"
      console.error msg

]
