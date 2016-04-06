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
    var account = $routeParams.account,
        descRuleById = {};

    $scope.createMode = (account != "!new");

    $scope.tableAll = true;
    $scope.ruleTable = "all";

    $scope.tepmAll = true;
    $scope.ruleTemp = "all";

    $scope.tableRules = [];
    $scope.tableSettings = [];

    Rules.List.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.rules = data.rules||[];
            $scope.rules.forEach(function(item) {
                descRuleById[item.id] = item.desc;
            });
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
                if (data.tableRule) {
                    $scope.tableAll = true;
                    $scope.ruleTable = "all";
                    $scope.rules.forEach(function(item) {
                        if (item.id == data.tableRule) {
                            $scope.tableRuleSelected = {};
                            $scope.tableRuleSelected = item;
                        }
                    });
                } else if (data.tableRules && data.tableRules.length > 0) {
                    $scope.tableSettings = [];
                    $scope.tableAll = false;
                    $scope.ruleTable = "some";
                    data.tableRules.forEach(function(item) {
                        var itemSet = item;
                        itemSet.ruleDesc = $scope.getDescRule(item.rule);
                        $scope.tableSettings.push(itemSet);
                    });
                }
            }
        });
    }

    $scope.getDescRule = function(id) {
        return descRuleById[id]||"";
    };

    $scope.tableOnSelect = function(v) {
        $scope.tableAll = (v == "all");
    };
    $scope.tempOnSelect = function(v) {
        $scope.tepmAll = (v == "all");
    };

    $scope.showAccessSettingsDlg = function() {
        $scope.tableRules = [];
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
        if ($scope.tableAll && $scope.tableRuleSelected && $scope.tableRuleSelected.id) {
            compileSettings.tableRule = $scope.tableRuleSelected.id;
        } else if (!$scope.tableAll && $scope.tableRules && $scope.tableRules.length > 0) {
            compileSettings.tableRules = $scope.tableRules;
        }
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
        $scope.tableSettings.forEach(function(item) {
            $scope.tableRules.push({table: item.table, rule: item.rule});
        });
        ngDialog.close();
    };
});