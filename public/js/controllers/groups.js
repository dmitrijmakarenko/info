infoSys.controller('groupsCntl', function ($scope, Groups) {
    $scope.groups = [{id: 'group1'}, {id: 'group2'}, {id: 'group3'}];

    $scope.go = function(group) {
        window.location = "#/groups/" + group;
    };

    Groups.List.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.groups = data.groups||[];
        }
    });
});

infoSys.controller('groupCntl', function ($scope, $routeParams, Groups) {
    var group = $routeParams.group;

    $scope.createMode = (group != "!new");

    if (group != "!new") {
        Groups.Get.go({id: group}, function(data) {
            console.log(data);
            $scope.name = data.name;
        });
    }

    $scope.groupUsers = [{name: "u1"}, {name: "u2"}, {name: "u3"}];
    $scope.otherUsers = [{name: "u111"}, {name: "u222"}];

    $scope.saveSettings = function() {
        Groups.Update.go({id: group, name: $scope.name}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Отдел сохранен");
                window.location = "#/groups/";
            }
        });
    };

    $scope.deleteGroup = function() {
        Groups.Delete.go({id: group}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Отдел удален");
                window.location = "#/groups/";
            }
        });
    };
});