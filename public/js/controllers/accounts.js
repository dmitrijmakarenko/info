accessSettings.controller('accountsCntl', function ($scope, Accounts) {

    $scope.go = function(account) {
        window.location = "#/accounts/" + account;
    };

    $scope.showSpinner();
    Accounts.List.go(function(data) {
        $scope.hideSpinner();
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

    $scope.tempRules = [];
    $scope.tempSettings = [];

    $scope.timeSettings = [
        {time: -1, text: "Неограниченно"},
        {time: 60, text: "Один час"},
        {time: 120, text: "Два часа"},
        {time: 180, text: "Три часа"},
        {time: 360, text: "Шесть часов"},
        {time: 720, text: "Двеннадцать часов"},
        {time: 1440, text: "Один день"}
    ];

    $scope.showSpinner();
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
        if (account == "!new") $scope.hideSpinner();
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.tables = data||[];
        }
    });

    if (account != "!new") {
        Accounts.Get.go({id: account}, function(data) {
            $scope.hideSpinner();
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

    $scope.showTempAccessSettingsDlg = function() {
        $scope.tempRules = [];
        ngDialog.open({
            template: 'accountAcsTmpDlgCntl',
            controller: 'accountAcsTempDlgCntl',
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
        if($scope.tepmAll && $scope.timeSelected && $scope.timeSelected.time) {
            compileSettings.tempRule = $scope.tableRuleSelected.id;
        } else if (!$scope.tepmAll && $scope.tempRules && $scope.tempRules.length > 0) {
            compileSettings.tempRules = $scope.tempRules;
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
        if ($scope.tableSelected && $scope.dlgRuleSelected) {
            var item = {};
            item.table = $scope.tableSelected.name;
            item.rule = $scope.dlgRuleSelected.id;
            item.ruleDesc = $scope.dlgRuleSelected.desc;
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

accessSettings.controller('accountAcsTempDlgCntl', function ($scope, ngDialog) {
    $scope.addTable = function() {
        if ($scope.tableSelected && $scope.dlgTimeSelected) {
            var item = {};
            item.table = $scope.tableSelected.name;
            item.time = $scope.dlgTimeSelected.time;
            item.timeText = $scope.dlgTimeSelected.text;
            $scope.tempSettings.push(item);
        }
    };

    $scope.removeTable = function(idx) {
        $scope.tempSettings.splice(idx, 1);
    };

    $scope.acceptSettings = function() {
        $scope.tempSettings.forEach(function(item) {
            $scope.tempRules.push({table: item.table, time: item.time});
        });
        ngDialog.close();
    };
});