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

    if (rule != "!new") {
        Rules.Get.go({id: rule}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.desc = data.desc;
            }
        });
    }

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
