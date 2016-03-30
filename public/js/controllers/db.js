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

accessSettings.controller('vcsCntl', function ($scope, VCS) {
    var loadData = function() {
        VCS.Tables.go(function(data) {
            console.log("vcs", data);
            if (data.error) {
                $scope.showErrorMsg(data.error)
            } else {
                $scope.vcsTables = data.tablesVcs||[];
                $scope.otherTables = [];
                var allTables = data.tablesAll||[];
                for (var  i = 0; i < allTables.length; i++) {
                    var addedVcs = false;
                    for (var  j = 0; j < $scope.vcsTables.length; j++) {
                        if (allTables[i] == $scope.vcsTables[j]) addedVcs = true;
                    }
                    if (!addedVcs) $scope.otherTables.push(allTables[i]);
                }
            }
        });
    };

    $scope.addToVcs = function(table) {
        VCS.AddToVcs.go({table: table}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error)
            } else {
                loadData();
            }
        });
    };

    $scope.removeFromVcs = function(table) {
        VCS.RemoveFromVcs.go({table: table}, function(data) {
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
    var token;
    $scope.auth = function() {
        Test.Auth.go({user: $scope.userName, pass: $scope.userPass}, function(data) {
            //console.log("auth", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg(data.token);
                token = data.token;
            }
        });
    };

    $scope.reset = function() {
        Test.Reset.go(function(data) {
            //console.log("reset", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Успешно");
            }
        });
    };

    $scope.init = function() {
        Test.Init.go(function(data) {
            //console.log("init", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Успешно");
            }
        });
    };

    $scope.work = function() {
        Test.Work.go(function(data) {
            //console.log("work", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Успешно");
            }
        });
    };

    $scope.compile = function() {
        Test.Compile.go(function(data) {
            console.log("compile", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Успешно");
            }
        });
    };

    $scope.copyToFile = function() {
        Test.CopyToFile.go(function(data) {
            //console.log("copyToFile", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Успешно");
            }
        });
    };

    $scope.copyFromFile = function() {
        Test.CopyFromFile.go(function(data) {
            //console.log("copyFromFile", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Успешно");
            }
        });
    };

    $scope.install = function() {
        Test.Install.go(function(data) {
            console.log("install", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Успешно");
            }
        });
    };

    $scope.goTest = function() {
        if ($scope.table) {

            var params = {};
            params.table = $scope.table;
            params.token = token||$scope.token||"";

            Data.Get.go({params: JSON.stringify(params)}, function(data) {
                console.log("get data", data);
                if (data.error) {
                    $scope.entity = {};
                    $scope.showErrorMsg(data.error);
                } else {
                    $scope.entity = data;
                    if (!$scope.entity.rows) $scope.entity.rows = [];
                }
            });
        }
    };
});