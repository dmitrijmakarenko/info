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
        console.log(data);
        if (data.error) {
            $scope.showErrorMsg(data.error)
        } else {
            $scope.users = data.users||[];
            $scope.groups = data.groups||[];
        }
    });

    $scope.options = [];
    $scope.options.push({t: "user1"});
    $scope.options.push({t: "grp1"});
    $scope.options.push({t: "user2"});

    $scope.selectObject = function(select) {
        $scope.selectUsers = select;
    };

    $scope.addOption = function() {
        $scope.options.push({t: "user2"});
    };

    $scope.removeOption = function(idx) {
        $scope.options.splice(idx, 1);
    };


    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.desc = $scope.desc;
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
