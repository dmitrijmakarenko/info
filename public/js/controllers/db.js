accessSettings.controller('dbCntl', function ($scope, DataBase) {
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
            if (data.error) {
                $scope.showErrorMsg(data.error)
            } else {
                loadData();
            }
        });
    };

    loadData();
});

accessSettings.controller('testCntl', function ($scope, Data, Test) {
    $scope.auth = function() {
        Test.Auth.go({user: $scope.userAuth, token: "8shnv6271jodoisxuyy274g1b"}, function(data) {
            console.log(data);
        });
    };

    $scope.reset = function() {
        Test.VCSreset.go(function(data) {
            console.log(data);
        });
    };


    $scope.goTest = function() {
        if ($scope.user && $scope.table) {

            var params = {};
            params.user = $scope.user;
            params.table = $scope.table;
            params.token = "8shnv6271jodoisxuyy274g1b";

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