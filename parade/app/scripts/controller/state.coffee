angular.module('tdbApp')
.controller 'StateCtrl',
[
  '$scope'
  '$routeParams'
  '$location'
  'Shelters'
  'Analytics'
  (
    $scope
    $routeParams
    $location
    Shelters
    Analytics
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
      Analytics.trackPage "/#{$routeParams.stateId}"
      Analytics.trackTrans()

      Shelters.getStatesMap().then (statesMap) ->
        $scope.states = $scope.getStates(statesMap)
        $scope.current.state = $scope.getCurrentState($scope.states)
        $scope.shelters = $scope.getCurrentShelters(statesMap)
        $scope.shelterIds = (shelter.id for shelter in $scope.shelters)

        $scope.$broadcast('show-animals', true)

    $scope.getStates = (statesMap) ->
      Object.keys(statesMap).sort()

    $scope.getCurrentState = (states) ->
      for state in states
        if $routeParams.stateId == state.toLowerCase()
          return state
      states[0]

    $scope.getCurrentShelters = (statesMap) ->
      statesMap[$scope.current.state]

    $scope.setCurrentState = (state) ->
      $scope.current.state = state
      $location.path state.toLowerCase()

    $scope.copyrightOwners = () ->
      (shelter.fullName for shelter in $scope.shelters).join(', ')
]
