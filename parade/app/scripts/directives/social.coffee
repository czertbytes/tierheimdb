angular.module('tdbApp')
.directive 'socialButtons', ->
  restrict: 'A'
  scope:
    url: '@'
  templateUrl: 'views/social-buttons.html'
  link: ($scope, elm, attrs) ->
    gapi.plusone.go('social-buttons')
    FB.init({
      logging: false
      xfbml: true
    })

