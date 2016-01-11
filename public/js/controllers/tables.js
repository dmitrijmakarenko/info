infoSys.controller('tablesCntl', function ($scope, Entity) {
    Entity.List.go(function(data) {
        $scope.entities = data;
    });

    $scope.getEntity = function(entity) {
        window.location = "#/table/" + entity;
    }
});

infoSys.controller('tableCntl', function ($scope, $routeParams, Entity) {
    //$scope.loading = true;
    var table = $routeParams.tableId;

    Entity.Get.go({id: table}, function(data) {
        if (!data.error) {
            $scope.entity = data;
            if (!$scope.entity.rows) $scope.entity.rows = [];
        }
        console.log(data);
    });
});
