angular.module('tdbApp')
.controller 'DocsCtrl', ['$scope', 'Analytics'
(
  $scope
  Analytics
) ->

  $scope.init = ->
    Analytics.trackPage "/docs"
    Analytics.trackTrans()
]
