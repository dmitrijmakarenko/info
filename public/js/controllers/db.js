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

infoSys.controller('testCntl', function ($scope, Data) {
    $scope.goTest = function() {
        if ($scope.user && $scope.table) {

            var params = {};
            params.user = $scope.user;
            params.table = $scope.table;

            Data.Get.go({params: JSON.stringify(params)}, function(data) {
                console.log(data);
                if (!data.error) {
                    $scope.entity = data;
                    if (!$scope.entity.rows) $scope.entity.rows = [];
                } else {
                    $scope.entity = {};
                    $scope.showErrorMsg(data.error);
                }
            });
        }
    };
});