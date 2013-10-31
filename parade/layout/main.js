(function(){
  "use strict";

  var apiBase = 'http://api.tierheimdb.de/v1';
  var animalsLimit = 20;

  var app = angular.module('tdb', ['infinite-scroll']);

  app.controller('AnimalListCtrl', function AnimalListCtrl($scope, $http, $q) {
    $scope.shelters = [];
    $scope.cities = [];
    $scope.animals = [];
    $scope.cats = 0;
    $scope.dogs = 0;

    $scope.page = 0;
    $scope.currentCity = 'Loading ...';
    $scope.typeFilter = {
      cats: true,
      dogs: true
    };
    $scope.busy = false;
    $scope.busyCounter = 0;

    $scope.$watch('typeFilter', function() {
      $scope.fetchAnimals();
    }, true);

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
      $.each(shelters, function(index, value) {
        cities[value.city] = true;
      });

      return Object.keys(cities);
    }

    $scope.updateShelters = function(result) {
      if (result.status == 200) {
        $scope.shelters = result.data;
        $scope.cities = extractCities($scope.shelters);
        $scope.currentCity = $scope.cities[0];

        $scope.refreshGallery();
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

      $.each($scope.shelters, function(index, value) {
        if (value.city == $scope.currentCity) {
          var url = apiBase + '/shelter/' + value.id + '/updates?limit=1';
          promises.push($http.get(url).then($scope.updateAnimalCounters));
        }
      });

      return $q.all(promises);
    }

    var prepareAnimal = function(animal) {
      if (animal.images == null || animal.images.length == 0) {
          animal.images = [];
          animal.images.push({
            url: 'http://placehold.it/260x260&text=' + animal.name
          });
      }

      return animal;
    }

    $scope.updateAnimals = function(result) {
      if (result.status == 200) {
        $.each(result.data, function(index, animal) {
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
      $.each($scope.shelters, function(index, value) {
        if (value.city == $scope.currentCity) {
          var url = apiBase + '/shelter/' + value.id + '/animals?limit=' + animalsLimit + '&type=' + $scope.getAnimalTypesParam();
          promises.push($http.get(url).then($scope.updateAnimals));
        }
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
      // tohle potrebuje fixnout!
      if (offset >= ($scope.cats + $scope.dogs)) {
        return;
      }

      $scope.busy = true;
      $scope.busyCounter = 0;

      $.each($scope.shelters, function(index, value) {
        if (value.city == $scope.currentCity) {
          $scope.busyCounter += 1;
          var url = apiBase + '/shelter/' + value.id + '/animals?offset=' + offset + '&limit=' + animalsLimit + '&type=' + $scope.getAnimalTypesParam();
          promises.push($http.get(url).then($scope.updateAnimals));
        }
      });

      return $q.all(promises);
    }

    $scope.getAnimalTypesParam = function() {
      var animalTypes = [];
      if ($scope.typeFilter.cats == true) {
          animalTypes.push("cat");
      }
      if ($scope.typeFilter.dogs == true) {
          animalTypes.push("dog");
      }

      return animalTypes.join(",");
    }

    $scope.setCurrentCity = function(city) {
      $scope.currentCity = city;
      $scope.page = 0;
      $scope.cats = 0;
      $scope.dogs = 0;
      $scope.refreshGallery();
    }
  });

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

})();
