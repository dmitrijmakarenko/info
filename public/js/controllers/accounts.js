accessSettings.controller('accountsCntl', function ($scope, Accounts) {

    $scope.go = function(account) {
        window.location = "#/accounts/" + account;
    };

    Accounts.List.go(function(data) {
        $scope.accounts = [];
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            var accounts = data.accounts||[];
            for (var i = 0; i < accounts.length; i ++) {
                var item = {};
                item.id = accounts[i].id;
                item.name = accounts[i].name||accounts[i].id;
                $scope.accounts.push(item);
            }
        }
    });
});

accessSettings.controller('accountCntl', function ($scope, $routeParams, Accounts, Rules, DataBase, ngDialog) {
    var account = $routeParams.account;

    $scope.createMode = (account != "!new");

    $scope.tableAll = true;
    $scope.ruleTable = "all";

    $scope.tepmAll = true;
    $scope.ruleTemp = "all";

    Rules.List.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.rules = data.rules||[];
        }
    });
    DataBase.GetTables.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.tables = data||[];
        }
    });

    if (account != "!new") {
        Accounts.Get.go({id: account}, function(data) {
            //console.log("get user", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.id = account;
                $scope.name = data.name;
                $scope.position = data.position;
            }
        });
    }

    $scope.tableOnSelect = function(v) {
        $scope.tableAll = (v == "all");
    };
    $scope.tempOnSelect = function(v) {
        $scope.tepmAll = (v == "all");
    };

    $scope.showAccessSettingsDlg = function() {
        ngDialog.open({
            template: 'accountAcsDlgCntl',
            controller: 'accountAcsDlgCntl',
            disableAnimation: true,
            showClose: false,
            scope: $scope
        });
    };

    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.id = $scope.id;
        compileSettings.name = $scope.name;
        compileSettings.password = $scope.password;
        compileSettings.position = $scope.position;
        Accounts.Update.go({id: account, settings: JSON.stringify(compileSettings)}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Пользователь сохранен");
                window.location = "#/accounts/";
            }
        });
    };

    $scope.deleteAccount = function() {
        Accounts.Delete.go({id: account}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Пользователь удален");
                window.location = "#/accounts/";
            }
        });
    };
});

accessSettings.controller('accountAcsDlgCntl', function ($scope, ngDialog) {
    $scope.tableSettings = [];

    $scope.addTable = function() {
        if ($scope.tableSelected && $scope.ruleSelected) {
            var item = {};
            item.table = $scope.tableSelected.name;
            item.rule = $scope.ruleSelected.id;
            item.ruleDesc = $scope.ruleSelected.desc;
            $scope.tableSettings.push(item);
        }
    };

    $scope.removeTable = function(idx) {
        $scope.tableSettings.splice(idx, 1);
    };

    $scope.acceptSettings = function() {
        ngDialog.close();
    };
});