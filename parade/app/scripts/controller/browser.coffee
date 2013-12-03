angular.module('tdbApp')
.controller 'BrowserCtrl', ['$scope', 'Updates', 'Animals',
(
  $scope
  Updates
  Animals
) ->

  $scope.updates = []
  $scope.animals = []
  $scope.busy = true
  $scope.page = 0

  $scope.init = ->

  $scope.$on('show-animals', (value) ->
    $scope.updates = []
    Updates.getLast($scope.shelterIds).then setUpdates

    $scope.page = 0
    $scope.animals = []
    type = $scope.current.type
    Animals.getAll($scope.shelterIds, type).then setAnimals
  )

  $scope.fetchNext = ->
    if $scope.busy == true
      return

    $scope.page += 1
    offset = $scope.page * 20
    if offset >= maxOffset()
      return

    console.log "page: #{$scope.page} offset: #{offset}"
    $scope.busy = true

    type = $scope.current.type
    Animals.getAll($scope.shelterIds, type, offset).then setAnimals

  setUpdates = (responses) ->
    for response in responses
      if response.status == 200
        $scope.updates.push response.data[0]

  setAnimals = (responses) ->
    for response in responses
      if response.status == 200
        fixedAnimals = (prepareAnimal(animal) for animal in response.data)
        $scope.animals = $scope.animals.concat fixedAnimals
    $scope.busy = false

  prepareAnimal = (animal) ->
    if animal.images && animal.images.length > 0
      animal
    else
      animal.images = [
        { url: "http://placehold.it/260x260&text=#{animal.name}" }
      ]
      animal

  maxOffset = () ->
    cats = 0
    dogs = 0
    for update in $scope.updates
      cats += update.cats
      dogs += update.dogs

    cats + dogs

]
