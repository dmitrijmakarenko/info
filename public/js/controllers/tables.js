accessSettings.controller('tablesCntl', function ($scope, Entity) {
    Entity.List.go(function(data) {
        $scope.entities = data;
    });

    $scope.getEntity = function(entity) {
        window.location = "#/table/" + entity;
    };

});

accessSettings.controller('tableCntl', function ($scope, $routeParams, Entity) {
    //$scope.loading = true;
    var table = $routeParams.tableId;

    Entity.Get.go({id: table}, function(data) {
        //console.log("table data", data);
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.entity = data;
            if (!$scope.entity.rows) $scope.entity.rows = [];
        }
    });
});
