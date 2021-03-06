angular.module('tdbApp',
[
  'ngRoute'
  'infinite-scroll'
  'angular-google-analytics'
])
.config ['$routeProvider', '$locationProvider', 'AnalyticsProvider'
, ($routeProvider, $locationProvider, AnalyticsProvider) ->
  $locationProvider.html5Mode on

  route = $routeProvider.when

  route '/docs',
    templateUrl: '/views/docs.html'
    controller: 'DocsCtrl'

  route '/about',
    templateUrl: '/views/about.html'
    controller: 'AboutCtrl'

  route '/:stateId',
    templateUrl: '/views/state.html'
    controller: 'StateCtrl'

  route '/:stateId/:shelterId',
    templateUrl: '/views/shelter.html'
    controller: 'ShelterCtrl'

  route '/:shelterId/:updateId/:animalId',
    templateUrl: '/views/animal.html'
    controller: 'AnimalCtrl'

  $routeProvider.otherwise
    redirectTo: '/berlin'

  AnalyticsProvider.setAccount 'UA-44691876-1'
]
