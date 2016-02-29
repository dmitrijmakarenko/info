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
    $scope.actions = [];

    $scope.createMode = (rule != "!new");
    $scope.selectUsers = true;

    $scope.operations = [
        {operation: "select", text: "Просмотр"},
        {operation: "insert", text: "Добавление"},
        {operation: "update", text: "Изменение"}
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
        if ($scope.operationSelected) {
            item.operation = $scope.operationSelected.operation;
        }
        console.log(item);
        if (item.object && item.operation) {
            $scope.actions.push(item);
        }
    };

    $scope.removeOption = function(idx) {
        $scope.actions.splice(idx, 1);
    };


    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.desc = $scope.desc;
        compileSettings.actions = $scope.actions;
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
