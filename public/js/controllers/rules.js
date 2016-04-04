accessSettings.controller('rulesCntl', function ($scope, Rules) {

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

accessSettings.controller('ruleCntl', function ($scope, $routeParams, Rules) {
    var rule = $routeParams.rule,
        userNameById = {},
        groupNameById = {};

    $scope.actions = [];
    $scope.createMode = (rule != "!new");
    $scope.selectUsers = true;

    $scope.operations = [
        {operation: "r", text: "Просмотр"},
        {operation: "w", text: "Изменение"}
    ];

    Rules.Data.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.users = data.users||[];
            $scope.groups = data.groups||[];
            for (var i = 0; i < $scope.users.length; i++) {
                if ($scope.users[i].id) userNameById[$scope.users[i].id] = $scope.users[i].name||$scope.users[i].id;
            }
            for (var i = 0; i < $scope.groups.length; i++) {
                if ($scope.groups[i].id) groupNameById[$scope.groups[i].id] = $scope.groups[i].name||$scope.groups[i].id;
            }
        }
    });

    if (rule != "!new") {
        Rules.Get.go({id: rule}, function(data) {
            console.log(data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.desc = data.desc;
                $scope.actions = data.actions||[];
            }
        });
    }

    $scope.getName = function(isUser, id) {
        if (isUser) {
            return userNameById[id]||"";
        } else {
            return groupNameById[id]||"";
        }
    };

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
