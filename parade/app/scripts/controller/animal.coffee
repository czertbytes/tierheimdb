angular.module('tdbApp')
.controller 'AnimalCtrl',
[
  '$scope'
  '$routeParams'
  'Animals'
  'Shelters'
  'Analytics'
  (
    $scope
    $routeParams
    Animals
    Shelters
    Analytics
  ) ->

    $scope.animal = {}
    $scope.shelter = {}

    $scope.init = ->
      shelterId = $routeParams.shelterId
      updateId = $routeParams.updateId
      animalId = $routeParams.animalId

      Analytics.trackPage "/{{shelterId}}/{{updateId}}/{{animalId}}"
      Analytics.trackTrans()

      Animals.get(shelterId, updateId, animalId).then setAnimal
      Shelters.get(shelterId).then setShelter

    setAnimal = (response) ->
      if response.status == 200
        $scope.animal = response.data
      else
        msg = "Getting animal detail failed! Status: #{response.status}"
        msg += "Error: #{response.data}"
        console.error msg

    setShelter = (response) ->
      if response.status == 200
        $scope.shelter = response.data
      else
        msg = "Getting shelter failed! Status: #{response.status}"
        msg += "Error: #{response.data}"
        console.error msg

    $scope.copyrightOwners = () ->
      $scope.shelter.fullName
]
