angular.module('tdbApp')
.controller 'ShelterCtrl',
[
  '$scope'
  '$routeParams'
  'Shelters'
  'Analytics'
  (
    $scope
    $routeParams
    Shelters
    Analytics
  ) ->


    $scope.shelterIds = []
    $scope.current =
      shelter: {}
      type: 'all'

    $scope.$watch('current.type', ->
      $scope.$broadcast('show-animals', true)
    , true)

    $scope.init = ->
      Analytics.trackPage "/#{$routeParams.stateId}/#{$routeParams.shelterId}"
      Analytics.trackTrans()

      Shelters.getAll().then (shelters) ->
        $scope.current.shelter = $scope.getShelterByRouteParams(shelters.data)
        $scope.shelterIds = [ $scope.current.shelter.id ]

        $scope.$broadcast('show-animals', true)

    $scope.getShelterByRouteParams = (shelters) ->
      for shelter in shelters
        if $routeParams.shelterId == shelter.id
          return shelter
      shelters[0]

    $scope.copyrightOwners = () ->
      $scope.current.shelter.fullName
]
