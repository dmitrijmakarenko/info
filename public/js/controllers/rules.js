infoSys.controller('rulesCntl', function ($scope, Rules) {

    $scope.go = function(rule) {
        window.location = "#/rules/" + rule;
    };

    Rules.List.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.rules = data.rules||[];
        }
    });

});

infoSys.controller('ruleCntl', function ($scope, $routeParams, Rules) {
    var rule = $routeParams.rule;
    $scope.options = [];

    $scope.createMode = (rule != "!new");
    $scope.selectUsers = true;

    $scope.operations = [
        {action: "select", text: "Просмотр"},
        {action: "insert", text: "Добавление"},
        {action: "update", text: "Изменение"}
    ];

    if (rule != "!new") {
        Rules.Get.go({id: rule}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.desc = data.desc;
            }
        });
    }

    Rules.Data.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.users = data.users||[];
            $scope.groups = data.groups||[];
        }
    });

    $scope.selectObject = function(select) {
        $scope.selectUsers = select;
    };

    $scope.addOption = function() {
        var item = {};
        if ($scope.selectUsers) {
            item.object = $scope.userSelected.id||null;
            item.isUser = true;
        } else {
            item.object = $scope.groupSelected.id||null;
            item.isUser = false;
        }
        if ($scope.actionSelected) {
            item.action = $scope.actionSelected.action;
        }
        if (item.object && item.action) {
            $scope.options.push(item);
        }
    };

    $scope.removeOption = function(idx) {
        $scope.options.splice(idx, 1);
    };


    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.desc = $scope.desc;
        compileSettings.options = $scope.options;
        Rules.Update.go({id: rule, settings: JSON.stringify(compileSettings)}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Правило сохранено");
                window.location = "#/rules/";
            }
        });
    };

    $scope.deleteRule = function() {
        Rules.Delete.go({id: rule}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Правило удалено");
                window.location = "#/rules/";
            }
        });
    };

});
