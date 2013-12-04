angular.module('tdbApp')
.controller 'AboutCtrl', ['$scope', 'Analytics'
(
  $scope
  Analytics
) ->

  $scope.init = ->
    Analytics.trackPage "/about"
    Analytics.trackTrans()

]
