angular.module('tdbApp')
.controller 'StateCtrl', ['$scope', '$routeParams', '$location', 'Shelters',
(
  $scope
  $routeParams
  $location
  Shelters
) ->

  $scope.states = []
  $scope.shelters = []
  $scope.shelterIds = []
  $scope.current =
    state: 'Loading ...'
    type: 'all'

  $scope.$watch('current.type', ->
    $scope.$broadcast('show-animals', true)
  , true)

  $scope.init = ->
    Shelters.getStatesMap().then (statesMap) ->
      $scope.states = $scope.getStates(statesMap)
      $scope.current.state = $scope.getCurrentStateByRouteParams($scope.states)
      $scope.shelters = $scope.getCurrentSheltersByRouteParams(statesMap)
      $scope.shelterIds = (shelter.id for shelter in $scope.shelters)

      $scope.$broadcast('show-animals', true)

  $scope.getStates = (statesMap) ->
    Object.keys(statesMap).sort()

  $scope.getCurrentStateByRouteParams = (states) ->
    for state in states
      if $routeParams.stateId == state.toLowerCase()
        return state
    states[0]

  $scope.getCurrentSheltersByRouteParams = (statesMap) ->
    statesMap[$scope.current.state]

  $scope.setCurrentState = (state) ->
    $scope.current.state = state
    $location.path state.toLowerCase()

]
