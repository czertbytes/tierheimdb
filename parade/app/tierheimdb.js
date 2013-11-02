(function(){
  "use strict";

  var apiBase = 'http://api.tierheimdb.de/v1';
  var animalsLimit = 20;

  var app = angular.module('TierheimDB', ['infinite-scroll']);

  app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
    $routeProvider
      .when('/tiere', {
        templateUrl: 'app/partials/tiere.html',
        controller: 'AnimalListCtrl'
      })
      .when('/docs', {
        templateUrl: 'app/partials/docs.html',
        controller: 'DocsCtrl'
      })
      .when('/about', {
        templateUrl: 'app/partials/about.html',
        controller: 'AboutCtrl'
      })
      .when('/:shelterId/:updateId/:animalId', {
        templateUrl: 'app/partials/animal-detail.html',
        controller: 'AnimalDetailCtrl'
      })
      .otherwise({
        redirectTo: '/tiere'
      });

    $locationProvider.html5Mode(true);
  }]);

  app.controller('AnimalListCtrl', ['$scope', '$http', '$q', function AnimalListCtrl($scope, $http, $q) {
    $scope.shelters = [];
    $scope.cities = [];
    $scope.animals = [];
    $scope.cats = 0;
    $scope.dogs = 0;

    $scope.page = 0;
    $scope.currentCity = 'Loading ...';
    $scope.currentShelters = [];
    $scope.typeFilter = 'all';
    $scope.busy = true;
    $scope.busyCounter = 0;

    $scope.$watch('typeFilter', function() {
      $scope.page = 0;
      $scope.fetchAnimals();
    });

    $scope.init = function() {
      $scope.fetchShelters();
    }

    $scope.refreshGallery = function() {
      $scope.busy = true;
      $scope.fetchLastUpdate().then($scope.fetchAnimals);
      $scope.busy = false;
    }

    var extractCities = function(shelters) {
      var cities = {};
      angular.forEach(shelters, function(shelter, index) {
        cities[shelter.city] = true;
      });

      return Object.keys(cities);
    }

    $scope.updateShelters = function(result) {
      if (result.status == 200) {
        $scope.shelters = result.data;
        $scope.cities = extractCities($scope.shelters);
        $scope.setCurrentCity($scope.cities[0]);
      }
    }

    $scope.fetchShelters = function() {
      return $http.get(apiBase + '/shelters').then($scope.updateShelters);
    }

    $scope.updateAnimalCounters = function(result) {
      if (result.status == 200) {
        $scope.cats += result.data[0].cats;
        $scope.dogs += result.data[0].dogs;
      }
    }

    $scope.fetchLastUpdate = function() {
      var promises = [];

      angular.forEach($scope.currentShelters, function(shelter, index) {
        var url = apiBase + '/shelter/' + shelter.id + '/updates?limit=1';
        promises.push($http.get(url).then($scope.updateAnimalCounters));
      });

      return $q.all(promises);
    }

    var prepareAnimal = function(animal) {
      if (animal.images && animal.images.length > 0) {
        return animal;
      } else {
          animal.images = [];
          animal.images.push({
            url: 'http://placehold.it/260x260&text=' + animal.name
          });
      }

      return animal;
    }

    $scope.updateAnimals = function(result) {
      if (result.status == 200) {
        angular.forEach(result.data, function(animal, index) {
          $scope.animals.push(prepareAnimal(animal));
        });
      }

      $scope.busyCounter -= 1;
      if ($scope.busyCounter == 0) {
        $scope.busy = false;
      }
    }

    $scope.fetchAnimals = function() {
      var promises = [];

      $scope.animals = [];
      angular.forEach($scope.currentShelters, function(shelter, index) {
        var url = apiBase + '/shelter/' + shelter.id + '/animals?limit=' + animalsLimit + '&type=' + $scope.getAnimalTypesParam();
        promises.push($http.get(url).then($scope.updateAnimals));
      });

      return $q.all(promises);
    }

    $scope.fetchNextAnimals = function() {
      var promises = [];

      if ($scope.busy == true) {
        return;
      }

      $scope.page += 1;
      var offset = $scope.page * animalsLimit;
      if (offset >= ($scope.cats + $scope.dogs)) {
        return;
      }

      $scope.busy = true;
      $scope.busyCounter = 0;

      angular.forEach($scope.currentShelters, function(shelter, index) {
        $scope.busyCounter += 1;
        var url = apiBase + '/shelter/' + shelter.id + '/animals?offset=' + offset + '&limit=' + animalsLimit + '&type=' + $scope.getAnimalTypesParam();
        promises.push($http.get(url).then($scope.updateAnimals));
      });

      return $q.all(promises);
    }

    $scope.getAnimalTypesParam = function() {
      if ($scope.typeFilter == 'cat' || $scope.typeFilter == 'dog') {
        return $scope.typeFilter;
      }

      return "";
    }

    $scope.setCurrentShelters = function() {
      $scope.currentShelters = [];
      angular.forEach($scope.shelters, function(shelter, index) {
        if (shelter.city == $scope.currentCity) {
          $scope.currentShelters.push(shelter);
        }
      });
    }

    $scope.setCurrentCity = function(city) {
      $scope.currentCity = city;
      $scope.setCurrentShelters();

      $scope.page = 0;
      $scope.cats = 0;
      $scope.dogs = 0;
      $scope.refreshGallery();
    }

    $scope.shelterFullNames = function() {
      var fullNames = [];
      angular.forEach($scope.currentShelters, function(shelter, index) {
        fullNames.push(shelter.fullName);
      });

      return fullNames.join(", ");
    }

  }]);

  app.controller('AnimalDetailCtrl', function AnimalDetailCtrl($scope, $http, $routeParams, $q) {
    $scope.animal = {};
    $scope.shelter = {};

    $scope.init = function() {
      $scope.fetchAnimal().then($scope.fetchShelter);
    }

    $scope.updateAnimal = function(result) {
      if (result.status = 200) {
        $scope.animal = result.data;
      }
    }

    $scope.fetchAnimal = function() {
      var url = apiBase + '/shelter/' + $routeParams.shelterId + '/update/' + $routeParams.updateId + '/animal/' + $routeParams.animalId;
      return $http.get(url).then($scope.updateAnimal);
    }

    $scope.updateShelter = function(result) {
      if (result.status == 200) {
        $scope.shelter = result.data;
      }
    }

    $scope.fetchShelter = function() {
      var url = apiBase + '/shelter/' + $routeParams.shelterId;
      return $http.get(url).then($scope.updateShelter);
    }

  });

  app.controller('DocsCtrl', function DocsCtrl($scope, $http, $routeParams, $q) {
  });

  app.controller('AboutCtrl', function AboutCtrl($scope, $http, $routeParams, $q) {
  });

})();
