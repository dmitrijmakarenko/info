infoSys.controller('dbCntl', function ($scope, DataBase) {
    var loadData = function() {
        if ($scope.showTables) {
            DataBase.GetTables.go(function(data) {
                console.log(data);
                $scope.tables = data;
                $scope.views = [];
            });
        } else {
            DataBase.GetViews.go(function(data) {
                console.log(data);
                $scope.views = data;
                $scope.tables = [];
            });
        }
    };
    $scope.showTables = true;

    $scope.toggle = function(what) {
        if (what == "views") {
            $scope.showTables = false;
            loadData();
        } else {
            $scope.showTables = true;
            loadData();
        }
    };

    $scope.modify = function(table) {
        DataBase.Protect.go({table: table}, function(data) {
            loadData();
        });
    };

    loadData();
});