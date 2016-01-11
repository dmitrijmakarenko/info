infoSys.controller('tablesCntl', function ($scope, Entity, TestFunc) {
    Entity.List.go(function(data) {
        $scope.entities = data;
    });

    $scope.getEntity = function(entity) {
        window.location = "#/table/" + entity;
    };

    $scope.test = function() {
        console.log("test");
        TestFunc.Select.go(function(data) {
            console.log(data);
        });
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
